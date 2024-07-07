// package utils

// import (
// 	"log"

// 	"github.com/gocql/gocql"
// )

// func InitScyllaDB() (*gocql.Session, error) {
// 	cluster := gocql.NewCluster("127.0.0.1")
// 	cluster.Port = 9042
// 	cluster.Keyspace = "todo_app"
// 	cluster.Consistency = gocql.Quorum
// 	session, err := cluster.CreateSession()
// 	if err != nil {
// 		log.Fatal("Error connecting to ScyllaDB: ", err)
// 		return nil, err
// 	}
// 	return session, nil
// }


package utils

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocql/gocql"
)

func SetupScyllaDB() (*gocql.Session, error) {
	cluster := gocql.NewCluster("127.0.0.1") // ScyllaDB host
	cluster.Port = 9042                      // ScyllaDB port
	cluster.Keyspace = "todo_app"            // ScyllaDB keyspace

	// Attempt to create a session
	session, err := cluster.CreateSession()
	if err == nil {
		return session, nil
	}

	// Check if the error is due to the keyspace not existing
	if !strings.Contains(err.Error(), "Keyspace 'todo_app' does not exist") {
		return nil, fmt.Errorf("failed to create session for todo_app keyspace: %v", err)
	}

	// Connect to the system keyspace to create the desired keyspace
	cluster.Keyspace = "system"
	session, err = cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session for system keyspace: %v", err)
	}

	// Create the todo_app keyspace
	err = session.Query(`CREATE KEYSPACE todo_app WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}`).Exec()
	if err != nil {
		return nil, fmt.Errorf("failed to create keyspace todo_app: %v", err)
	}

	// Connect to the newly created keyspace
	cluster.Keyspace = "todo_app"
	session, err = cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session for todo_app keyspace: %v", err)
	}

	// Create the todos table if it does not exist
	err = session.Query(`
		CREATE TABLE IF NOT EXISTS todos (
			id UUID PRIMARY KEY,
			user_id TEXT,
			title TEXT,
			description TEXT,
			status TEXT,
			created TIMESTAMP,
			updated TIMESTAMP
		)
	`).Exec()
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	log.Println("Keyspace 'todo_app' and table 'todos' created successfully")
	return session, nil
}
