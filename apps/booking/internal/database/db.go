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
	dbname := "booking_db"

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var db *sql.DB
	var err error

	for i := 1; i <= 5; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				log.Println("✅ Connexion à PostgreSQL (Booking) réussie !")
				
				autoMigrate(db)
				
				db.SetMaxOpenConns(25)
				db.SetMaxIdleConns(25)
				db.SetConnMaxLifetime(5 * time.Minute)
				
				return db
			}
		}
		log.Printf("Base injoignable (tentative %d/5)...", i)
		time.Sleep(3 * time.Second)
	}
	log.Fatalf("❌ Échec critique de connexion : %v", err)
	return nil
}

func autoMigrate(db *sql.DB) {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS bookings (
		id SERIAL PRIMARY KEY,
		trip_id INT NOT NULL,
		user_email VARCHAR(255) NOT NULL,
		status VARCHAR(50) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("❌ Erreur lors de la création de la table bookings : %v", err)
	}

	log.Println("✅ Table 'bookings' vérifiée/initialisée avec succès.")
}