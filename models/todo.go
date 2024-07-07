package models

import "github.com/gocql/gocql"

// Todo represents a TODO item
type Todo struct {
	ID          gocql.UUID `json:"id"`
	UserID      string     `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Created     int64      `json:"created"`
	Updated     int64      `json:"updated"`
}
