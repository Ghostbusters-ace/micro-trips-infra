# Micro-Trips — Core Application Services & Software Engineering

Ce dépôt contient l'intégralité de l'intelligence métier et du code source de l'application **Micro-Trips**. Développé entièrement en **Go (Golang)** pour des raisons d'efficience, de concourance native (Goroutines) et de légèreté des binaires, ce projet adopte une architecture découplée, orientée événements (*Event-Driven*).

---

## Architecture Logicielle : Clean Architecture Décortiquée

Chaque microservice applicatif rejette le modèle monolithique ou la structure en "script unique" au profit d'un découpage en **Clean Architecture**. L'organisation interne isole strictement le code technique des règles métier à travers 4 couches étanches :

```
          +---------------------------------------------------+
          | 🌐 couche transport : Handlers (HTTP / REST)      |
          |      +-----------------------------------------+  |
          |      | 🧠 couche métier : Services (Logique)   |  |
          |      |      +-------------------------------+  |  |
          |      |      | 💾 couche accès : Repositories | |  |
          |      |      |      +---------------------+   | |  |
          |      |      |      | 💎 domaine : Models |   | |  |
          |      |      |      +---------------------+   | |  |
          |      |      +-------------------------------+  |  |
          |      +-----------------------------------------+  |
          +---------------------------------------------------+

```

### Description des Couches Applicatives

1. **`internal/models/` (Domaine) :** Définition pure des entités structurelles (`Trip`, `Booking`). Aucun tag de framework ou de bibliothèque externe n'y pénètre. C'est le cœur invariant de l'application.
2. **`internal/repository/` (Infrastructure de données) :** Encapsule l'accès aux bases de données PostgreSQL. Il implémente des interfaces strictes. La logique métier ne sait pas *comment* les données sont enregistrées, elle délègue cette tâche au repository.
3. **`internal/service/` (Logique Métier Pure) :** C'est ici que résident les règles décisionnelles (ex: valider si un e-mail est correct, calculer un statut, orchestrer la création). Cette couche interagit avec la couche d'accès aux données uniquement via des abstractions (interfaces Go), permettant de substituer la vraie base de données par des simulacres (*mocks*) lors des tests unitaires.
4. **`internal/handler/` (Interface de Transport) :** Points d'entrée réseau (API REST / JSON). Cette couche décode les requêtes HTTP, intercepte les erreurs pour injecter les codes HTTP appropriés (`201 Created`, `400 Bad Request`, `500 Server Error`), et gère le middleware d'observabilité Prometheus.

---

## plumbing: Focus Ingénierie : Le Pattern Event-Driven (Publisher/Consumer)

Pour garantir une tolérance aux pannes maximale, la création d'une réservation et l'envoi de la notification par e-mail sont totalement asynchrones :

### Le Publisher : `apps/booking`

Lorsqu'un client émet une requête `POST /bookings`, le cycle interne s'active :

1. Le `BookingHandler` réceptionne et valide le payload JSON.
2. Le `BookingService` persiste la réservation en base de données avec le statut initial `PENDING` via le `BookingRepository`.
3. Le module **`internal/messaging/publisher.go`** prend le relais. Il ouvre un canal AMQP sur le broker RabbitMQ, sérialise un payload structuré appelé `EventPayload` au format JSON, et publie le message de manière persistante (DeliveryMode: 2) dans l'exchange ciblant la file d'attente `bookings_queue`.
4. L'API HTTP répond immédiatement au client en moins de 10ms, sans attendre la finalisation du traitement de l'e-mail.

### Le Consumer : `apps/worker`

Le microservice `Notification Worker` s'exécute en tâche de fond comme un démon d'infrastructure indépendant :

1. Au démarrage, son point d'entrée `main.go` appelle `messaging.InitRabbitMQ()` pour établir une connexion permanente.
2. Il déclare la queue `bookings_queue` et invoque la méthode de consommation non-bloquante `Channel.Consume`.
3. Via une boucle de scrutation infinie (`for d := range msgs`) s'exécutant dans une **Goroutine dédiée**, il intercepte chaque message dès sa publication sur le broker.
4. Le payload JSON brut est désérialisé dans la structure `EventPayload`. Si le décodage réussit, il invoque le package `internal/mailer`, qui compose un e-mail au standard MIME (gérant le charset UTF-8 et les headers de sujet) et l'expédie via le protocole SMTP sur le port `1025` de l'outil **MailHog**.

---

## Architecture des Dockerfiles (Optimisation Multi-Stage)

L'écriture de nos fichiers `Dockerfile` suit scrupuleusement les exigences de production industrielle (sécurité, isolation et réduction drastique de la surface d'attaque). Chaque composant compile ses binaires via un processus **Multi-Stage** :

```dockerfile
# ========================================================
# STAGE 1 : Environnement de Compilation (Builder)
# ========================================================
FROM golang:1.22-alpine AS builder

# Définition du répertoire de travail
WORKDIR /app

# Cache des dépendances Go : évite de retélécharger les modules si go.mod n'a pas bougé
COPY go.mod go.sum ./
RUN go mod download

# Copie de l'intégralité des sources du service
COPY . .

# Compilation statique du binaire Go
# CGO_ENABLED=0 élimine les dépendances aux librairies C dynamiques de l'hôte
# GOOS=linux cible le noyau Linux de production
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/api/main.go

# ========================================================
# STAGE 2 : Image d'Exécution Sécurisée (Runner)
# ========================================================
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Récupération exclusive du binaire compilé au stage 1 (aucun code source, aucun outil de build)
COPY --from=builder /app/main .

# Création d'un utilisateur non-privilégié pour interdire l'exécution en tant que root
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Exposition du port applicatif
EXPOSE 8080

# Commande de démarrage
CMD ["./main"]

```

### Avantages de cette Approche FinOps / SecOps :

* **Poids Plume :** L'image finale pèse moins de **25 Mo** (au lieu de 800 Mo si l'environnement complet de Go était conservé).
* **Sécurité Accrue :** L'absence de compilateurs, d'outils système (comme `gcc`) et le passage à un utilisateur `non-root` empêchent l'exploitation de failles de sécurité par injection ou escalade de privilèges.

---

## Cycle de Vie des Données : Auto-Migration & Seeding

Afin de respecter les paradigmes GitOps, nos applications n'exigent aucune intervention humaine ou exécution de scripts SQL manuels à la création du cluster :

* **Auto-Migration :** Au démarrage, la méthode `InitDB()` de chaque service se connecte à PostgreSQL. Une fonction `autoMigrate(db *sql.DB)` analyse l'état du moteur de données et exécute automatiquement une requête `CREATE TABLE IF NOT EXISTS`. Si les tables manquent, elles se créent dynamiquement sans interruption de service.
* **Database Seeding :** Pour l'API `Catalog`, le code Go intègre une couche d'amorçage. Après avoir validé l'existence de la table `trips`, le service compte le nombre d'enregistrements. Si la table est vide, il injecte de manière programmatique notre catalogue initial (Paris, Tokyo, Bali), garantissant un environnement fonctionnel dès la première seconde.

---

## Pipeline d'Intégration Continue (GitHub Actions)

Le fichier `.github/workflows/ci.yml` orchestre automatiquement la validation de la qualité logicielles à chaque `git push` ou `pull_request` :

1. **Linting & Formating :** Analyse statique du code (`go fmt` et `go vet`) pour traquer les failles syntaxiques ou les variables inutilisées.
2. **Unit Testing :** Exécution de la suite de tests via la commande `go test -v ./...` pour garantir la non-régression de la logique métier.
3. **Build & Push :** Si les tests sont au vert, la pipeline s'authentifie sur notre registre de conteneurs privé. Elle génère une version de l'image Docker tagguée avec le SHA unique du commit Git.
4. **Authentification Sécurisée (Docker Secrets) :** Pour permettre aux noeuds Kubernetes du cluster local de télécharger l'image depuis notre registre privé sans laisser de mot de passe traîner, la pipeline s'appuie sur un secret Kubernetes de type `kubernetes.io/dockerconfigjson`. Ce secret, chiffré par l'administrateur via l'outil `kubeseal`, est versionné de manière totalement sûre sous la forme d'un fichier `SealedSecret` dans le dépôt GitOps.

---

## Organisation Détaillée du Code Source

```text
.
├── .github/
│   └── workflows/          # Config de la pipeline CI (GitHub Actions)
└── apps/
    ├── booking/            # Microservice de gestion des réservations
    │   ├── internal/
    │   │   ├── database/   # Initialisation du pool SQL + Auto-migration
    │   │   ├── handler/    # Contrôleur HTTP / Points d'entrée REST
    │   │   ├── messaging/  # Logique du Publisher RabbitMQ
    │   │   ├── models/     # Modèles du domaine (Booking struct)
    │   │   └── repository/ # Requêtes SQL et accès à PostgreSQL
    │   ├── main.go         # Bootstrap et Injection des dépendances du service
    │   └── Dockerfile      # Processus de build multi-stage
    ├── catalog/            # Microservice de consultation du catalogue
    │   ├── internal/
    │   │   ├── database/   # Initialisation PostgreSQL + Database Seeding
    │   │   ├── handler/    # Traitement des requêtes GET /trips
    │   │   ├── models/     # Modèles du domaine (Trip struct)
    │   │   └── repository/ # Lecture SQL du catalogue
    │   └── main.go
    └── worker/             # Démon de notification asynchrone
        ├── internal/
        │   ├── mailer/     # Service d'envoi SMTP (génération MIME/UTF-8)
        │   └── messaging/  # Logique du Consumer RabbitMQ (Goroutine infinie)
        └── main.go

```