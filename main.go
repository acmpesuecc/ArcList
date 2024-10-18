package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	ID       int
	Task     string
	Position int
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
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/updateOrder", updateOrderHandler)
	http.HandleFunc("/search", searchHandler)

	log.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func createTable() {
	query := `
    CREATE TABLE IF NOT EXISTS todos (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        task TEXT,
        position INTEGER
    );
    `
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	// Ensure the position column is populated
	migrateAddPositionColumn()
}

// Migrate to add 'position' column if it doesn't exist
func migrateAddPositionColumn() {
	_, err := db.Query("SELECT position FROM todos LIMIT 1")
	if err != nil {
		_, err = db.Exec("ALTER TABLE todos ADD COLUMN position INTEGER")
		if err != nil {
			log.Fatal(err)
		}
		
		_, err = db.Exec("UPDATE todos SET position = id WHERE position IS NULL")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, task, position FROM todos ORDER BY position")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		rows.Scan(&todo.ID, &todo.Task, &todo.Position)
		todos = append(todos, todo)
	}

	tpl.ExecuteTemplate(w, "index.html", todos)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		task := r.FormValue("task")
		if task != "" {
			_, err := db.Exec("INSERT INTO todos (task, position) VALUES (?, (SELECT COALESCE(MAX(position), 0) + 1 FROM todos))", task)
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

func editHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/edit/"):]

	var todo Todo
	err := db.QueryRow("SELECT id, task, position FROM todos WHERE id = ?", id).Scan(&todo.ID, &todo.Task, &todo.Position)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tpl.ExecuteTemplate(w, "edit.html", todo)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := r.FormValue("id")
		task := r.FormValue("task")

		if id != "" && task != "" {
			_, err := db.Exec("UPDATE todos SET task = ? WHERE id = ?", task, id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			renderTaskList(w)
		}
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")

	rows, err := db.Query("SELECT id, task, position FROM todos WHERE task LIKE ? ORDER BY position", "%"+query+"%")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		rows.Scan(&todo.ID, &todo.Task, &todo.Position)
		todos = append(todos, todo)
	}

	tpl.ExecuteTemplate(w, "tasklist", todos)
}

// Update task order based on drag-and-drop
func updateOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var request struct {
			Order []int `json:"order"`
		}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tx, err := db.Begin()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for i, id := range request.Order {
			_, err := tx.Exec("UPDATE todos SET position = ? WHERE id = ?", i+1, id)
			if err != nil {
				tx.Rollback()
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		err = tx.Commit()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}

func renderTaskList(w http.ResponseWriter) {
	rows, err := db.Query("SELECT id, task, position FROM todos ORDER BY position")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		rows.Scan(&todo.ID, &todo.Task, &todo.Position)
		todos = append(todos, todo)
	}

	tpl.ExecuteTemplate(w, "tasklist", todos)
}
