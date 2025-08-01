package main

import (
	"database/sql"
    _ "github.com/mattn/go-sqlite3"
	"log"
    "golang.org/x/crypto/bcrypt"
)

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

// maybe return a User pointer
func GetUser(db *sql.DB, username string) User {
	user := User{}
	const sentence = `SELECT name, password FROM users WHERE name = ?`
	row := db.QueryRow(sentence, username)
	err := row.Scan(&user.Name, &user.Password)
	if err != nil {
		log.Printf("cannot found user '%s': %v\n", username, err)
	}

	return user 
}

func GetUsers(db *sql.DB) []User {
	users := []User{}
	const sentence = `SELECT name, password FROM users`

	rows, err := db.Query(sentence)
	if err != nil {
		log.Println("error querying users:", err)
		return users
	}
	defer rows.Close()

	for rows.Next() {
		user := User{}

		err := rows.Scan(&user.Name, &user.Password)
		if err != nil {
			log.Println("error getting row:", err)	
			return users
		}

		users = append(users, user)
	}

	return users
}
