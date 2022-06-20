package main

import (
	"database/sql"
	_ "encoding/json"
	"fmt"
	"net/http"
	"time"
	"log"

	handler "Assigment2Golang/Handler"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres" 
	password = "Abcdzfgh123" //ganti sesuai nama password postgres
	dbname   = "db-go-sql"
)

var (
	db *sql.DB

	err error
)

const PORT = ":8080"

func main() {
	// mmemastikan db connect atau tidak
	db, err = sql.Open("postgres", ConnectDbPsql(host, user, password, dbname, port))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Succesfully connected to database")

	// handler crud
	r := mux.NewRouter()
	itemHandler := handler.NewItemHandler(db)
	r.HandleFunc("/orders", itemHandler.ItemsHandler)
	r.HandleFunc("/orders/{id}", itemHandler.ItemsHandler)

	fmt.Println("Now listening on port 0.0.0.0" + PORT)
	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0" + PORT,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func ConnectDbPsql(host, user, password, name string, port int) string {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbname)
	return psqlInfo
}