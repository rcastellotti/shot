// this is just a placeholder with some grok-generated junk

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

const dbFile = "shot.db"

type User struct {
	ID       int
	Email    string
	Verified bool
}

func main() {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTables(db)

	http.HandleFunc("/register", registerHandler(db))
	http.HandleFunc("/submit", authMiddleware(submitHandler(db)))
	http.HandleFunc("/verify", verifyHandler(db))

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createTables(db *sql.DB) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			verified BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS scripts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			url TEXT UNIQUE NOT NULL,
			proof TEXT,
			script_hash TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
	}

	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func registerHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if req.Email == "" || req.Password == "" {
			http.Error(w, "Email and password required", http.StatusBadRequest)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("INSERT INTO users (email, password_hash) VALUES (?, ?)", req.Email, hash)
		if err != nil {
			if err.Error() == "UNIQUE constraint failed: users.email" { // sqlite specific
				http.Error(w, "Email already registered", http.StatusConflict)
			} else {
				http.Error(w, "Internal error", http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Registered, pending admin verification"})
	}
}

func authMiddleware(next func(http.ResponseWriter, *http.Request, *User)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var u User
		var storedHash string
		err := db.QueryRow("SELECT id, email, password_hash, verified FROM users WHERE email = ?", email).
			Scan(&u.ID, &u.Email, &storedHash, &u.Verified)
		if err == sql.ErrNoRows {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)); err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !u.Verified {
			http.Error(w, "Account pending verification", http.StatusForbidden)
			return
		}

		next(w, r, &u)
	}
}

func submitHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, *User) {
	return func(w http.ResponseWriter, r *http.Request, user *User) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			URL   string `json:"url"`
			Proof string `json:"proof"`
			Hash  string `json:"hash"` // sha256 hex
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if req.URL == "" || req.Hash == "" {
			http.Error(w, "URL and hash required", http.StatusBadRequest)
			return
		}

		_, err := db.Exec("INSERT INTO scripts (user_id, url, proof, script_hash) VALUES (?, ?, ?, ?)",
			user.ID, req.URL, req.Proof, req.Hash)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed: scripts.url") {
				http.Error(w, "URL already submitted", http.StatusConflict)
			} else {
				http.Error(w, "Internal error", http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Script submitted"})
	}
}

func verifyHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "GET only", http.StatusMethodNotAllowed)
			return
		}

		queryURL := r.URL.Query().Get("url")
		if queryURL == "" {
			http.Error(w, "url query param required", http.StatusBadRequest)
			return
		}

		var scriptHash string
		err := db.QueryRow("SELECT script_hash FROM scripts WHERE url = ?", queryURL).Scan(&scriptHash)
		if err == sql.ErrNoRows {
			http.Error(w, "Script not registered", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"expected_hash": scriptHash})
	}
}
