package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
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
			log.Println("error parsing form:",err);
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

        err = db.QueryRow("SELECT name FROM users WHERE name = ?", username).Scan()
		if err == nil {
            log.Println("username:",username,"already exists!")
			return
        } else if err != sql.ErrNoRows {
			fmt.Println("db error:", err)
            return
        }

		user := User{Name: username, Password: password}
		CreateUser(db, user)
	}).Methods("POST")

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Content-Type", "application/json")
//		json.NewEncoder(w).Encode()
		err := r.ParseForm()
		if err != nil {
			log.Println("error parsing form:",err);
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		var hashedPassword string
        err = db.QueryRow("SELECT password FROM users WHERE name = ?", username).Scan(&hashedPassword)
		if err == nil {
			log.Printf("username: %s not found\n", username)
			return
        } 

		passwordMatch := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if passwordMatch != nil {
			log.Println("password doesnt match:", err)
			return
		}

	}).Methods("POST")

	router.HandleFunc("/get-users", func(w http.ResponseWriter, r *http.Request) {
  		w.Header().Set("Content-Type", "application/json")
  		json.NewEncoder(w).Encode(GetUsers(db))
	}).Methods("GET")

	CreateUsersTable(db)
	CreateUser(db, User{Name: "skibidi", Password: "asd123"})
	users := GetUsers(db)
	fmt.Println(users)

	const port = "8081"
	fmt.Printf("server listening in localhost:%s\n",port)
	http.ListenAndServe("localhost:" + port, router)
}
