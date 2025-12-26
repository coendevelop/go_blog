package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// Post represents a blog entry
type Post struct {
	ID      int
	Title   string
	Content string
}

var db *sql.DB

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Check for session cookie
	_, err := r.Cookie("session")
	isLoggedIn := err == nil // true if cookie exists, false if not

	rows, _ := db.Query("SELECT id, title, content FROM posts")
	var posts []Post

	for rows.Next() {
		var p Post
		rows.Scan(&p.ID, &p.Title, &p.Content)
		posts = append(posts, p)
	}

	// Pass data to template
	data := struct {
		Posts      []Post
		IsLoggedIn bool
	}{
		Posts:      posts,
		IsLoggedIn: isLoggedIn,
	}

	tmpl := template.Must(template.ParseFiles("templates/index.tmpl"))
	tmpl.Execute(w, data)
}

func newPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		title := r.FormValue("title")
		content := r.FormValue("content")
		db.Exec("INSERT INTO posts (title, content) VALUES (?, ?)", title, content)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	tmpl := template.Must(template.ParseFiles("templates/new.tmpl"))
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

	errr := initDB(db)
	if errr != nil {
		log.Fatal(errr)
	}

	// Insert a test user (password: "password123")
	// In a real app, you'd have a signup page.
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), 14)
	db.Exec("INSERT OR IGNORE INTO users (username, password) VALUES (?, ?)", "admin", hashedPassword)

	/* Routes
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/new", AuthRequired(newPostHandler)) // Only logged-in users can reach this
	http.HandleFunc("/login", loginHandler)*/

	// Recover Middleware
	http.HandleFunc("/", recoverMiddleware(indexHandler))
	http.HandleFunc("/new", recoverMiddleware(AuthRequired(newPostHandler)))
	http.HandleFunc("/login", recoverMiddleware(loginHandler))

	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
