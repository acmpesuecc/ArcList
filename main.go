package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type Todo struct {
	Title   string    `json:"title"`
	DueDate time.Time `json:"dueDate"`
}

var (
	todos = []Todo{}
	mu    sync.Mutex
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/todos", getTodos).Methods("GET")
	r.HandleFunc("/todos", addTodo).Methods("POST")

	http.Handle("/", http.FileServer(http.Dir("./"))) // Serve static files
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	// Sort todos by due date
	sort.Slice(todos, func(i, j int) bool {
		return todos[i].DueDate.Before(todos[j].DueDate)
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func addTodo(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse the due date
	dueDate, err := time.Parse("2006-01-02", r.FormValue("dueDate"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	todo.DueDate = dueDate

	todos = append(todos, todo)
	w.WriteHeader(http.StatusCreated)
}
