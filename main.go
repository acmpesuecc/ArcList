package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	ID   int
	Task string
}

var db *sql.DB
var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	// Open SQLite database
	var err error
	db, err = sql.Open("sqlite3", "./sqlite.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize the database
	createTable()

	// Route handlers
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/search", searchHandler) // New search handler

	log.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func createTable() {
	query := `
    CREATE TABLE IF NOT EXISTS todos (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        task TEXT
    );
    `
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, task FROM todos")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		rows.Scan(&todo.ID, &todo.Task)
		todos = append(todos, todo)
	}

	tpl.ExecuteTemplate(w, "index.html", todos)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		task := r.FormValue("task")
		if task != "" {
			_, err := db.Exec("INSERT INTO todos (task) VALUES (?)", task)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		renderTaskList(w)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := r.FormValue("id")
		if id != "" {
			_, err := db.Exec("DELETE FROM todos WHERE id = ?", id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		renderTaskList(w)
	}
}

// New search handler
func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	query = strings.ToLower(query)

	rows, err := db.Query("SELECT id, task FROM todos")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		rows.Scan(&todo.ID, &todo.Task)
		if strings.Contains(strings.ToLower(todo.Task), query) {
			todos = append(todos, todo)
		}
	}

	tpl.ExecuteTemplate(w, "tasklist", todos)
}

func renderTaskList(w http.ResponseWriter) {
	rows, err := db.Query("SELECT id, task FROM todos")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		rows.Scan(&todo.ID, &todo.Task)
		todos = append(todos, todo)
	}

	tpl.ExecuteTemplate(w, "tasklist", todos)
}
