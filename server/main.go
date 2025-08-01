package main

import (
	"log"
	"database/sql"
    _ "github.com/mattn/go-sqlite3"
)

type User struct {
	Name string
	Password string
}

func main()  {
	db, err := sql.Open("sqlite3", "database.db");
	if err != nil {
		log.Fatalln("error opening db:", err);
	}
	defer db.Close()

	CreateUsersTable(db)
	CreateUser(db, User{Name: "skibidi", Password: "asd123"})
}
