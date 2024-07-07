<h1 align="center">Backend Assignment Todo List</h1>

## Objective:
    Develop a TODO API using Golang and ScyllaDB that supports basic CRUD operations and includes pagination functionality for the list endpoint.


## Requirements:
    ● Set up a Golang project and integrate ScyllaDB as the database for storing TODO items.
    Ensure that items in the database are stored user-wise.
    ● Implement endpoints for creating, reading, updating, and deleting TODO items for a single user at a   time. Each TODO item should have at least the following properties: id, user_id, title, description status, created, updated.
    ● Implement a paginated list endpoint to retrieve TODO items.
    ● Provide support for filtering based on TODO item status (e.g., pending, completed).


## 🖥️ Run Locally

-Clone this repository
```git clone https://github.com/swapnil-sharma7222/Todo-List.git```
-Install Docker Desktop.
-Run `docker-compose up -d` to initialize a ScyllaDB instance running on port 9042.
-Go to the cmd/api directory and run `go run main.go` to start the service.


Note: If you are facing issues connecting the Go backend services with ScyllaDB, follow these steps:

Open the terminal and enter the command to execute the CQL shell in the ScyllaDB instance:

```docker exec -it scylla-db cqlsh```
Note the IP (172.23.0.2) (it might change according to your Docker configuration) and change it in the utils/db.go file.

In the CQL shell terminal, create the keyspace manually with the following command:

```CREATE KEYSPACE todo_app WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};```

Use the `todo_app keyspace`:


```USE todo_app;```

Create the todos table with the following command:


```CREATE TABLE IF NOT EXISTS todos (id UUID PRIMARY KEY, user_id TEXT, title TEXT, description TEXT, status TEXT, created TIMESTAMP, updated TIMESTAMP);```


## ⚙️ API Documentation
  ### `POST` /todos

    Create a new todo item.

    URL: `http://localhost:8080/todos`
    Method: POST
    Body: JSON
    ```
    json example
    {
      "user_id": "004599c1-69d7-47b8-8306-a144d7265538",
      "title": "Complete project tasks",
      "description": "Finish coding the backend and write documentation",
      "status": "Pending"
    }
    ```


  ### `PUT` /todos/{id}

    Update a specific todo item.

    URL: `http://localhost:8080/todos/{id}`
    Method: PUT
    URL Parameters:
    id: ID of the todo item to update


  ### `DELETE` /todos/{id}

    Delete a specific todo item.

    URL: `http://localhost:8080/todos/{id}`
    Method: DELETE
    URL Parameters:
    id: ID of the todo item to delete


  ### `GET` /todos/{id}

    Retrieve details of a specific todo item.

    URL: `http://localhost:8080/todos/{id}`
    Method: GET
    URL Parameters:
    id: ID of the todo item to retrieve


  ### `GET` /todos

    Retrieve all todo items based on filters.

    URL: `http://localhost:8080/v1/todos`
    Method: GET
    Query Parameters:
    status: Filter by status (pending or completed)
    size: Number of items to retrieve (default 10)
    lastPageToken: Offset for pagination (default 0)


## Features Implemented:

-Implemented CRUD routes for interaction between server and ScyllaDB.
-The APIs are paginated for easy data retrieval.
-The application's DB part is Dockerized and is stateful through volumes.
-Support for filtering based on TODO item status (e.g., pending, completed).


## Current Architecture:

-Containerized approach to solve the problem statement.
-Interfaces for the server and DB are interacting with each other for the backend application.

## Future Scope:

-This is a basic implementation of the problem statement.
-Depending on the scale, the architecture can be scaled horizontally using Nginx load balancing.
-The app can use a queueing mechanism like RabbitMQ or BullMQ to introduce pub-sub architecture for better performance.
-The Go server can be containerized for improved deployment.
-Introduction of goroutines would increase the overall throughput of the service.


## 👨🏻‍💻 Developer's Talk
Developed by <a href="https://github.com/swapnil-sharma7222">Swapnil Sharma</a>

<a href="https://github.com/swapnil-sharma7222/Todo-List">This</a> is a small effort from my side to build a small-scale project using Golang and ScyllaDB technologies. The experience taught me many things, as well as the challenges involved in overcoming problems encountered during the development phase. I consider this project very relevant to me as a full-stack developer.

<br/>