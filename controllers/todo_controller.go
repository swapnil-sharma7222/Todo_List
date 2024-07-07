package controllers

import (
	"time"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"todo-app/models"
	"todo-app/repositories"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)


type TodoController struct {
	repo repositories.TodoRepository
}

// Creates a new TodoController with the given repository
func NewTodoController(repo repositories.TodoRepository) *TodoController {
	return &TodoController{repo: repo}
}

// Creates a new todo item
func (c *TodoController) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if todo.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	if todo.UserID == "" {
		http.Error(w, "UserID is required", http.StatusBadRequest)
		return
	}

	todo.ID = gocql.TimeUUID()
	currentTime := time.Now().Unix()
	todo.Created = currentTime
	todo.Updated = currentTime

	err = c.repo.Create(&todo)
	if err != nil {
		log.Println("Error inserting todo item:", err)
		http.Error(w, "Failed to create TODO item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

// Fetching a todo item by its ID
func (c *TodoController) GetTodoByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoID := vars["id"]

	id, err := gocql.ParseUUID(todoID)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	todo, err := c.repo.GetByID(id)
	if err != nil {
		if err == gocql.ErrNotFound {
			http.Error(w, "Todo item not found", http.StatusNotFound)
		} else {
			log.Println("Error retrieving todo item:", err)
			http.Error(w, "Failed to retrieve TODO item", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// Updating an existing todo item by its ID
func (c *TodoController) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoID := vars["id"]

	id, err := gocql.ParseUUID(todoID)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	var updatedTodo models.Todo
	err = json.NewDecoder(r.Body).Decode(&updatedTodo)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if updatedTodo.ID != id {
		http.Error(w, "Todo ID mismatch", http.StatusBadRequest)
		return
	}

	err = c.repo.Update(&updatedTodo)
	if err != nil {
		log.Println("Error updating todo item:", err)
		http.Error(w, "Failed to update TODO item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Deleting a todo item by its ID
func (c *TodoController) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoID := vars["id"]

	id, err := gocql.ParseUUID(todoID)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	err = c.repo.Delete(id)
	if err != nil {
		log.Println("Error deleting todo item:", err)
		http.Error(w, "Failed to delete TODO item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Fetching a list of todo items
func (c *TodoController) ListTodos(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	pageStr := queryParams.Get("page")
	sizeStr := queryParams.Get("size")
	status := queryParams.Get("status")
	lastPageToken := queryParams.Get("lastPageToken")

	page, size := 1, 10
	var err error
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			http.Error(w, "Invalid page number", http.StatusBadRequest)
			return
		}
	}
	if sizeStr != "" {
		size, err = strconv.Atoi(sizeStr)
		if err != nil {
			http.Error(w, "Invalid page size", http.StatusBadRequest)
			return
		}
	}

	todos, err := c.repo.List(status, page, size, lastPageToken)
	if err != nil {
		http.Error(w, "Failed to fetch TODO items", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}
