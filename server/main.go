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
		err := r.ParseForm()
		if err != nil {
			log.Println("error parsing form:",err);
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		var name string
        err = db.QueryRow("SELECT name FROM users WHERE name = ?", username).Scan(&name)
		if err == nil {
            log.Println("username:",username,"already exists!")
			http.Error(w, fmt.Sprintf("Username: %s already exists", username), http.StatusBadRequest)
			return
        } else if err != sql.ErrNoRows {
			http.Error(w, "Database error", http.StatusBadRequest)
			fmt.Println("db error:", err)
            return
        }

		user := User{Name: username, Password: password}
		CreateUser(db, user)
		http.Error(w, fmt.Sprintf("User: %s created successfully", username), http.StatusOK)
	}).Methods("POST")

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Println("error parsing form:",err);
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		var hashedPassword string
        err = db.QueryRow("SELECT password FROM users WHERE name = ?", username).Scan(&hashedPassword)
		if err != nil {
			log.Printf("username: %s not found\n", username)
			http.Error(w, fmt.Sprintf("Username or password incorrect"), http.StatusBadRequest)
			return
        } 

		passwordMatch := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if passwordMatch != nil {
			log.Println("password doesnt match:", err)
			http.Error(w, fmt.Sprintf("Username or password incorrect"), http.StatusBadRequest)
			return
		}

		http.Error(w, fmt.Sprintf("Login successfull"), http.StatusOK)
	}).Methods("POST")

	router.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
	})

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
