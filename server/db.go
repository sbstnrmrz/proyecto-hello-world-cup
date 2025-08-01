package main

import (
	"database/sql"
    _ "github.com/mattn/go-sqlite3"
	"log"
    "golang.org/x/crypto/bcrypt"
)

type User struct {
	Name string
	Password string
}


func CreateUsersTable(db *sql.DB) {
	const sentence = `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`

	_, err := db.Exec(sentence)
	if err != nil {
		log.Println("users table already created")
		return
	}

	log.Println("users table created")
}

func CreateUser(db *sql.DB, user User) {
	const sentence = `INSERT INTO users (name, password) VALUES (?, ?)`

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("error hashing password")
		return
	}

	_, err = db.Exec(sentence, user.Name, hashedPassword)
	if err != nil {
		log.Printf("error creating user '%s': %v\n", user, err)
		return
	}

	log.Printf("user: '%s' created\n", user.Name)
}

func GetUser(db *sql.DB, name string) User {
	const query = `SELECT name, password FROM user WHERE name = ? LIMIT 1`

	result, err := db.Exec(query, name)

	if err != nil {
		log.Printf("error getting user '%s': %v\n", name, err)
		return User{};
	}

	log.Printf("user: '%s' selected\n", name)
	return 
}
