package repositories

import (
	"log"
	"todo-app/models"

	"github.com/gocql/gocql"
)

type TodoRepository interface {
	Create(todo *models.Todo) error
	GetByID(id gocql.UUID) (*models.Todo, error)
	Update(todo *models.Todo) error
	Delete(id gocql.UUID) error
	List(status string, page, size int, lastPageToken string) ([]models.Todo, error)
}

type ScyllaTodoRepository struct {
	// ScyllaDB session for database operations
	session *gocql.Session
}

// New instance of ScyllaTodoRepository
func NewTodoRepository(session *gocql.Session) TodoRepository {
	return &ScyllaTodoRepository{session: session}
}

// Add new TODO item to the database
func (r *ScyllaTodoRepository) Create(todo *models.Todo) error {
	return r.session.Query(
		"INSERT INTO todos (id, user_id, title, description, status, created, updated) VALUES (?, ?, ?, ?, ?, ?, ?)",
		todo.ID, todo.UserID, todo.Title, todo.Description, todo.Status, todo.Created, todo.Updated,
	).Exec()
}

// Gets a todo item by its ID
func (r *ScyllaTodoRepository) GetByID(id gocql.UUID) (*models.Todo, error) {
	var todo models.Todo
	err := r.session.Query(
		"SELECT id, user_id, title, description, status, created, updated FROM todos WHERE id = ?",
		id).Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.Status, &todo.Created, &todo.Updated)
	return &todo, err
}

// Updates an existing todo item in the database
func (r *ScyllaTodoRepository) Update(todo *models.Todo) error {
	return r.session.Query(
		"UPDATE todos SET title = ?, description = ?, status = ?, updated = ? WHERE id = ?",
		todo.Title, todo.Description, todo.Status, todo.Updated, todo.ID).Exec()
}

// Deletes a todo item by its ID
func (r *ScyllaTodoRepository) Delete(id gocql.UUID) error {
	return r.session.Query("DELETE FROM todos WHERE id = ?", id).Exec()
}

// Gets todo items from the database
func (r *ScyllaTodoRepository) List(status string, page, size int, lastPageToken string) ([]models.Todo, error) {
	var todos []models.Todo
	var query string
	var iter *gocql.Iter

	// Construct the final query
	if status != "" {
		if lastPageToken != "" {
			query = "SELECT id, user_id, title, description, status, created, updated FROM todos WHERE status = ? AND token(id) > token(?) LIMIT ? ALLOW FILTERING"
			iter = r.session.Query(query, status, lastPageToken, size).Iter()
		} else {
			query = "SELECT id, user_id, title, description, status, created, updated FROM todos WHERE status = ? LIMIT ? ALLOW FILTERING"
			iter = r.session.Query(query, status, size).Iter()
		}
	} else {
		if lastPageToken != "" {
			query = "SELECT id, user_id, title, description, status, created, updated FROM todos WHERE token(id) > token(?) LIMIT ? ALLOW FILTERING"
			iter = r.session.Query(query, lastPageToken, size).Iter()
		} else {
			query = "SELECT id, user_id, title, description, status, created, updated FROM todos LIMIT ? ALLOW FILTERING"
			iter = r.session.Query(query, size).Iter()
		}
	}

	// Final Result
	for {
		todo := models.Todo{}
		if !iter.Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.Status, &todo.Created, &todo.Updated) {
			break
		}
		todos = append(todos, todo)
	}

	// Error Handlers
	if err := iter.Close(); err != nil {
		log.Println("Error querying todos:", err)
		return nil, err
	}

	return todos, nil
}


