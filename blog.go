package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// Post represents a blog entry
type Post struct {
	ID      int
	Title   string
	Content string
}

var db *sql.DB

func indexHandler(w http.ResponseWriter, r *http.Request) {
	rows, _ := db.Query("SELECT id, title, content FROM posts")
	var posts []Post

	for rows.Next() {
		var p Post
		rows.Scan(&p.ID, &p.Title, &p.Content)
		posts = append(posts, p)
	}

	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, posts)
}

func newPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		title := r.FormValue("title")
		content := r.FormValue("content")
		db.Exec("INSERT INTO posts (title, content) VALUES (?, ?)", title, content)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	tmpl := template.Must(template.ParseFiles("new.html"))
	tmpl.Execute(w, nil)
}

func main() {
	var err error
	// Initialize SQLite connection
	db, err = sql.Open("sqlite3", "./blog.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create table if it doesn't exist
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS posts (id INTEGER PRIMARY KEY, title TEXT, content TEXT)")
	statement.Exec()

	// Routes
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/new", newPostHandler)

	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
