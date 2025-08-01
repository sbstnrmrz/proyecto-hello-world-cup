package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)
func enableCORS(w *http.ResponseWriter) {
	// Puedes usar '*' para permitir cualquier origen (útil para desarrollo, pero no recomendado en producción sin precauciones)
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:4321") // ¡Cambia esto por el origen de tu frontend!
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true") // Si necesitas manejar cookies o credenciales
}

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
		enableCORS(&w)
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
		enableCORS(&w)
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

		sessionToken := GenerateSessionToken()
		sessions[sessionToken] = username
		log.Println(sessionToken)

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
		})

		http.Error(w, fmt.Sprintf("Login successfull"), http.StatusOK)
  		json.NewEncoder(w).Encode(sessionToken)
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
