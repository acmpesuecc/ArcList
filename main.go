package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	ID   int
	Task string
	Due  string // Date without time
}

var db *sql.DB
var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
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
	http.HandleFunc("/search", searchHandler)

	log.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// Create table with `due_date`
func createTable() {
	query := `
    CREATE TABLE IF NOT EXISTS todos (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        task TEXT,
        due_date DATE
    );
    `
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, task, DATE(due_date) FROM todos ORDER BY due_date ASC") // Use DATE to ensure no time
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		rows.Scan(&todo.ID, &todo.Task, &todo.Due)
		todos = append(todos, todo)
	}

	tpl.ExecuteTemplate(w, "index.html", todos)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		task := r.FormValue("task")
		due := r.FormValue("due_date") // This is already in YYYY-MM-DD format from HTML input

		if task != "" {
			_, err := db.Exec("INSERT INTO todos (task, due_date) VALUES (?, ?)", task, due)
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

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")

	rows, err := db.Query("SELECT id, task, DATE(due_date) FROM todos WHERE task LIKE ?", "%"+query+"%") // Use DATE to strip time
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		rows.Scan(&todo.ID, &todo.Task, &todo.Due)
		todos = append(todos, todo)
	}

	tpl.ExecuteTemplate(w, "tasklist", todos)
}

func renderTaskList(w http.ResponseWriter) {
	rows, err := db.Query("SELECT id, task, DATE(due_date) FROM todos ORDER BY due_date ASC") // Ensure date only
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		rows.Scan(&todo.ID, &todo.Task, &todo.Due)
		todos = append(todos, todo)
	}

	tpl.ExecuteTemplate(w, "tasklist", todos)
}
