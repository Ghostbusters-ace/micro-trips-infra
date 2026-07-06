package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := "catalog_db"

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var db *sql.DB
	var err error

	for i := 1; i <= 5; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				log.Println("✅ Connexion à PostgreSQL (Catalog) réussie !")
				
				autoMigrate(db)
				
				return db
			}
		}
		log.Printf("⏳ Base injoignable (tentative %d/5)...", i)
		time.Sleep(3 * time.Second)
	}
	log.Fatalf("❌ Échec critique de connexion : %v", err)
	return nil
}

func autoMigrate(db *sql.DB) {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS trips (
		id SERIAL PRIMARY KEY,
		destination VARCHAR(100) NOT NULL,
		price DECIMAL(10,2) NOT NULL,
		description TEXT
	);`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("❌ Erreur lors de la création de la table : %v", err)
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM trips").Scan(&count)
	
	if count == 0 {
		log.Println("Base de données vide : Insertion des voyages par défaut (Seeding)...")
		insertData := `
		INSERT INTO trips (destination, price, description) VALUES 
		('Paris, France', 150.00, 'Un week-end romantique'),
		('Tokyo, Japon', 1200.00, 'Immersion totale dans la culture japonaise'),
		('Bali, Indonésie', 800.00, 'Détente et plages paradisiaques');`
		
		_, err = db.Exec(insertData)
		if err != nil {
			log.Fatalf("❌ Erreur lors de l'insertion des données : %v", err)
		}
	} else {
		log.Println("La table trips contient déjà des données, on passe l'initialisation.")
	}
}