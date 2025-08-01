package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func GenerateSessionToken() string {
    b := make([]byte, 32)
    rand.Read(b)
    return base64.URLEncoding.EncodeToString(b)
}

func main()  {
	db, err := sql.Open("sqlite3", "database.db");
	if err != nil {
		log.Fatalln("error opening db:", err);
	}
	defer db.Close()

	var sessions = make(map[string]string)

	router := mux.NewRouter()

	router.HandleFunc("/create-user", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Println("error parsing form:",err);
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		regexp.Compile("[a-zA-Z1-9.]+@unet.edu.ve$")

		var savedEmail string
        err = db.QueryRow("SELECT email FROM users WHERE email = ?", email).Scan(&savedEmail)
		if err == nil {
			log.Println("account with email:",email,"already exists!")
			http.Error(w, fmt.Sprintf("account with email: %s", email), http.StatusBadRequest)
			return
        } else if err != sql.ErrNoRows {
			http.Error(w, "Database error", http.StatusBadRequest)
			fmt.Println("db error:", err)
            return
        }

		nickFromEmail := strings.Split(email, "@")[0]
		log.Printf("nick for %s: %s\n", email, nickFromEmail)

		user := User{
			Nick: nickFromEmail,
			Email: email,
			Password: password,
			TipoUsuario: 1,
		}
		CreateUser(db, user)
		http.Error(w, fmt.Sprintf("User for email: %s created successfully", email), http.StatusOK)
	}).Methods("POST")

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Println("error parsing form:",err);
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		var hashedPassword string
        err = db.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&hashedPassword)
		if err != nil {
			log.Printf("account with email: %s not found\n", email)
			http.Error(w, fmt.Sprintf("account with email: %s doesnt exists", email), http.StatusBadRequest)
			return
        } 

		passwordMatch := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if passwordMatch != nil {
			log.Println("password doesnt match:", err)
			http.Error(w, fmt.Sprintf("email or password incorrect"), http.StatusBadRequest)
			return
		}

		sessionToken := GenerateSessionToken()
		sessions[sessionToken] = email 

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			MaxAge:   5600,
			HttpOnly: true,
			Secure: true,
			Path: "/",
			SameSite: http.SameSiteLaxMode,
		})

		http.Error(w, fmt.Sprintf("Login successfull"), http.StatusOK)
	}).Methods("POST")

	router.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
	})

	router.HandleFunc("/get-users", func(w http.ResponseWriter, r *http.Request) {
  		w.Header().Set("Content-Type", "application/json")
  		json.NewEncoder(w).Encode(GetUsers(db))
	}).Methods("GET")

	DropUsersTable(db)
	CreateUsersTable(db)

	const port = "8081"
	fmt.Printf("server listening in localhost:%s\n",port)
	http.ListenAndServe("localhost:" + port, router)
	user := GetUser(db, "pedro.sanchez")
	fmt.Println(user)
}
