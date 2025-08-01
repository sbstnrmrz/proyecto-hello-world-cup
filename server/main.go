package main

import (
	"log"
	"database/sql"
    _ "github.com/mattn/go-sqlite3"

	"net/http"
	"github.com/gorilla/mux"
)

func main()  {
	db, err := sql.Open("sqlite3", "database.db");
	if err != nil {
		log.Fatalln("error opening db:", err);
	}
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/user/{}",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			out, err :=  c.GetDepartamentos(dbpool)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
				return
			}

			fmt.Println(r.URL)
			json.NewEncoder(w).Encode(out)
		}).Methods("GET")
	CreateUsersTable(db)
	CreateUser(db, User{Name: "skibidi", Password: "asd123"})


}
