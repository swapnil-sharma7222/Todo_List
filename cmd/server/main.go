package main

import (
	"log"
	"net/http"
	"todo-app/routes"
	"todo-app/utils"
)

func main() {
	// Initializing database
	session, err := utils.SetupScyllaDB()
	if err != nil {
		log.Fatal("Error initializing ScyllaDB: ", err)
	}
	defer session.Close()
	// Initializing server
	router := routes.NewRouter(session)

	// Listening on Port 8080
	log.Fatal(http.ListenAndServe(":8080", router))
	log.Println("Server running at http://localhost:8080/")
}
