package main

import (
	"fmt"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

func main() {
	fmt.Println("Hello, World!")

	db, err := sql.Open("pgx", "host=localhost port=5432 user=postgres password=password dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	createTablesIfNotExists(db)
	// insertUser(db, "John Doe", "john.doe@example.com", 25, "male")
	insertUser(db, "Jane Doe", "jane.doe@example.com", 25, "female")
	defer db.Close()
}

func insertUser(db *sql.DB, name string, email string, age int, sex string) int {
	var userId int
	tx, err := db.Begin()

	if err != nil {
		log.Fatal(err)
	}

	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", name, email)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.QueryRow("SELECT id FROM users WHERE name = $1 AND email = $2", name, email).Scan(&userId)
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx.Exec("INSERT INTO profile (user_id, age, sex) VALUES ($1, $2, $3)", userId, age, sex)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	return userId
}

func createTablesIfNotExists(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE
		)
	`)
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS profile (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			age INT NOT NULL,
			sex VARCHAR(255) NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}