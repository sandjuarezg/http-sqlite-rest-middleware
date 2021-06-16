package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sandjuarezg/http-sqlite-rest-middleware/server/database/model"
	"github.com/sandjuarezg/http-sqlite-rest-middleware/server/functionality"
)

type message struct {
	Body string `json:"body"`
}

func main() {
	functionality.SqlMigration()

	db, err := sql.Open("sqlite3", "./database/user.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.Handle("/add", postAdd(db))
	http.Handle("/show", getShow(db))
	http.Handle("/search", getSearch(db))

	fmt.Println("Listening on localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func postAdd(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("Accept") == "application/json" {

			w.Header().Set("Content-Type", "application/json")

			defer fmt.Printf("Response from %s\n", r.URL.RequestURI())

			var addU functionality.User
			var err = json.NewDecoder(r.Body).Decode(&addU)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = model.AddUser(db, addU)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = json.NewEncoder(w).Encode(message{Body: "Insert data successfully"})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		} else {
			w.WriteHeader(http.StatusNotAcceptable)
		}
	})
}

func getShow(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("Accept") == "application/json" {

			defer fmt.Printf("Response from %s\n", r.URL.RequestURI())

			w.Header().Set("Content-Type", "application/json")

			users, err := model.ShowUser(db)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = json.NewEncoder(w).Encode(users)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		} else {
			w.WriteHeader(http.StatusNotAcceptable)
		}
	})
}

func getSearch(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("Accept") == "application/json" {

			defer fmt.Printf("Response from %s\n", r.URL.RequestURI())

			w.Header().Set("Content-Type", "application/json")

			var searchU functionality.User
			var err = json.NewDecoder(r.Body).Decode(&searchU)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			user, err := model.SearchUser(db, searchU.Id)
			if err != nil {
				if err == sql.ErrNoRows {
					err = json.NewEncoder(w).Encode(message{Body: "User not found"})
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					return
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			err = json.NewEncoder(w).Encode(user)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		} else {
			w.WriteHeader(http.StatusNotAcceptable)
		}
	})
}
