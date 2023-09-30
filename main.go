package main

import (
	"database/sql"
	"log"
	"net"
	"net/http"
)

func main() {
	var db *sql.DB = connectToDatabase(DB_USERNAME, DB_PASSWORD, DB_ADDR)
	go cacheTrucks(db)
	
	http.HandleFunc("/trucks", trucksHandler(db))

	listener, error := net.Listen("tcp", ":8080")
	if error != nil {
		log.Fatal(error)
	}
	http.Serve(listener, nil)

}
