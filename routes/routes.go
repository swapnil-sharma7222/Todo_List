package routes

import (
	"todo-app/controllers"
	"todo-app/repositories"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

func NewRouter(session *gocql.Session) *mux.Router {
	// Importing the necessary functions
	repo := repositories.NewTodoRepository(session)
	controller := controllers.NewTodoController(repo)

	router := mux.NewRouter()
  // Endpoint to create todo
	router.HandleFunc("/todos", controller.CreateTodo).Methods("POST")

  // Endpoint to get a particular todo by id
	router.HandleFunc("/todos/{id}", controller.GetTodoByID).Methods("GET")

  // Endpoint to update a particular todo by id
	router.HandleFunc("/todos/{id}", controller.UpdateTodo).Methods("PUT")

  // Endpoint to delete a particular todo by id
	router.HandleFunc("/todos/{id}", controller.DeleteTodo).Methods("DELETE")
  
  // Endpoint to get all the todo
	router.HandleFunc("/todos", controller.ListTodos).Methods("GET")

	return router
}
