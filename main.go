package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

func generatePassword() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func connectDB(username, password, dbhost, dbname string) (*sql.DB, error) {
	return connectDBWithRetry(username, password, dbhost, dbname, 2)
}

func connectDBWithRetry(username, password, dbhost, dbname string, retryMinutes int) (*sql.DB, error) {
	dsnString := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", username, password, dbhost, dbname)
	var db *sql.DB
	var err error

	for i := 0; i < retryMinutes*60; i++ {
		db, err = sql.Open("postgres", dsnString)
		if err != nil {
			log.Printf("sql.Open failed with %s \n", dsnString)
			return nil, err
		}
		err = db.Ping()
		if err == nil {
			log.Printf("Connected to %s/%s \n", dbhost, dbname)
			break
		}
		log.Printf("Failed to connect to %s/%s on attempt %d.  Retrying...  \n", dbhost, dbname, i+1)
		time.Sleep(time.Second)
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}

func execSQL(db *sql.DB, query string) error {
	_, err := db.Exec(query)
	return err
}

func main() {

	config, err := readConfig("ochami.yaml")
	if err != nil {
		log.Fatal(err)
	}

	username := os.Getenv("DB_USER")
	if username == "" {
		log.Fatal("DB_USER is required")
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		log.Fatal("DB_PASSWORD is required")
	}

	dbhost := os.Getenv("DB_HOST")
	if dbhost == "" {
		log.Fatal("DB_HOST is required")
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		log.Fatal("DB_NAME is required")
	}

	db, err := connectDB(username, password, dbhost, dbname)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for _, database := range config.Databases {
		err = execSQL(db, fmt.Sprintf("CREATE DATABASE %s;", database.Name))
		if err != nil {
			log.Fatal(err)
		}

		for _, user := range database.Users {
			if user.Password == "" {
				user.Password, err = generatePassword()
				if err != nil {
					log.Fatal(err)
				}
			}

			err = execSQL(db, fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s';", user.Name, user.Password))
			if err != nil {
				log.Fatal(err)
			}

			err = execSQL(db, fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s;", database.Name, user.Name))
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	err = writeConfig("ochami.yaml", config)
	if err != nil {
		log.Fatal(err)
	}
}
