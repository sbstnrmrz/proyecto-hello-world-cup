package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main()  {
	db, err := sql.Open("sqlite3", "database.db");
	if err != nil {
		log.Fatalln("error opening db:", err);
	}
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/create-user", func(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Content-Type", "application/json")
//		json.NewEncoder(w).Encode()
		err := r.ParseForm()
		if err != nil {
			log.Println("error creating user:",err);
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

        err = db.QueryRow("SELECT username FROM users WHERE username = ?", username).Scan()
		if err == nil {
            log.Println("Username:",username,"already exists!")
			return
        } else if err != sql.ErrNoRows {
			fmt.Println("db error:", err)
            return
        }

		user := User{Name: username, Password: password}
		CreateUser(db, user)
	}).Methods("GET")

	router.HandleFunc("/get-users", func(w http.ResponseWriter, r *http.Request) {
  		w.Header().Set("Content-Type", "application/json")
  		json.NewEncoder(w).Encode(GetUsers(db))
	}).Methods("GET")

	CreateUsersTable(db)
	CreateUser(db, User{Name: "skibidi", Password: "asd123"})
	users := GetUsers(db)
	fmt.Println(users)
}
