package main

import (
	"database/sql"
    _ "github.com/mattn/go-sqlite3"
	"log"
    "golang.org/x/crypto/bcrypt"
)

type User struct {
	Nick string
	Email string
	Password string
	TipoUsuario int 
	DateCreated string
	LastSession string
}

func CreateUsersTable(db *sql.DB) {
	const sentence = `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nick TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		tipo_usuario INTEGER NOT NULL,
		date_created TEXT DEFAULT CURRENT_TIMESTAMP,
		last_session TEXT DEFAULT CURRENT_TIMESTAMP 
	);`

	_, err := db.Exec(sentence)
	if err != nil {
		log.Println("users table already created")
		return
	}

	log.Println("users table created")
}

func CreateUser(db *sql.DB, user User) {
	const sentence = `INSERT INTO users (nick, email, password, tipo_usuario) VALUES (?, ?, ?, ?)`

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("error hashing password")
		return
	}

	_, err = db.Exec(sentence, user.Nick, user.Email, hashedPassword, user.TipoUsuario)
	if err != nil {
		log.Printf("error creating user '%s': %v\n", user, err)
		return
	}

	log.Printf("user with email: '%s' created\n", user.Email)
}

// maybe return a User pointer
func GetUser(db *sql.DB, nick string) User {
	user := User{}
	const sentence = `SELECT nick, email, tipo_usuario, date_created, last_session FROM users WHERE nick = ?`
	row := db.QueryRow(sentence, nick)
	err := row.Scan(&user.Nick, &user.Email, &user.DateCreated, &user.LastSession)
	if err != nil {
		log.Printf("cannot found user '%s': %v\n", nick, err)
	}

	return user 
}

func GetUsers(db *sql.DB) []User {
	users := []User{}
	const sentence = `SELECT nick, email, date_created, last_session FROM users`

	rows, err := db.Query(sentence)
	if err != nil {
		log.Println("error querying users:", err)
		return users
	}
	defer rows.Close()

	for rows.Next() {
		user := User{}

		err := rows.Scan(&user.Nick, &user.Email, &user.DateCreated, &user.LastSession)
		if err != nil {
			log.Println("error getting row:", err)	
			return users
		}

		users = append(users, user)
	}

	return users
}

func DropUsersTable(db *sql.DB) {
	const sentence = `DROP TABLE IF EXISTS users;`
	_, err := db.Exec(sentence)
	if err != nil {
		log.Println("error dropping users table:", err)
		return
	}

	log.Println("users table dropped")
}
