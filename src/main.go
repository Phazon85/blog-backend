package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type user struct {
	Email *string
}

func main() {

	//define configuration
	//getting user credentials
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "changeme"
		dbname   = "test"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	//connection to database
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	//defining router
	r := mux.NewRouter()
	r.HandleFunc("/", getHomeHandler(db)).Methods("GET")
	r.HandleFunc("/", postHomeHandler(db)).Methods("POST")

	//running the server
	log.Fatal(http.ListenAndServe(":8080", r))
}

func getHomeHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//get user from databse
		var selectedUsers []*user
		rows, err := db.Query("SELECT email FROM user_email")
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			selectedUser := &user{}
			err := rows.Scan(&selectedUser.Email)
			if err != nil {
				panic(err)
			}
			selectedUsers = append(selectedUsers, selectedUser)
		}

		//respond to the user with the email
		resp, err := json.Marshal(selectedUsers)
		if err != nil {
			panic(err)
		}
		w.Write(resp)
	}
}

func postHomeHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//convert the request from json to something Go can read using the User struct
		var newUser user
		err := json.NewDecoder(r.Body).Decode(&newUser)
		if err != nil {
			panic(err)
		}
		//take value and use a SQL statement to insert it into our DB
		sqlStatement := fmt.Sprintf("INSERT INTO user_email (email) VALUES ('%s');", *newUser.Email)
		_, err = db.Exec(sqlStatement)
		if err != nil {
			panic(err)
		}
		//give user confirmation that it worked
		w.Write([]byte("SUCCESS - Dexter's Laboratory Voice"))
	}
}
