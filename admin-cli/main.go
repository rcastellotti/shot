// this is just a placeholder with some grok-generated junk
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

func adminCLI() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: admin list-users | admin verify <email>")
		os.Exit(1)
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	switch os.Args[2] {
	case "list-users":
		rows, err := db.Query("SELECT id, email, verified, created_at FROM users ORDER BY created_at DESC")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		fmt.Println("ID\tEmail\tVerified\tCreated")
		for rows.Next() {
			var id int
			var email string
			var verified bool
			var created time.Time
			rows.Scan(&id, &email, &verified, &created)
			status := "pending"
			if verified {
				status = "verified"
			}
			fmt.Printf("%d\t%s\t%s\t%s\n", id, email, status, created.Format(time.RFC3339))
		}
	case "verify":
		if len(os.Args) < 4 {
			fmt.Println("Usage: admin verify <email>")
			os.Exit(1)
		}
		email := os.Args[3]
		result, err := db.Exec("UPDATE users SET verified = TRUE WHERE email = ?", email)
		if err != nil {
			log.Fatal(err)
		}
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			fmt.Println("No user found with that email")
		} else {
			fmt.Println("User verified")
		}
	default:
		fmt.Println("Unknown command")
	}
}
