package database

import (
	"database/sql"
	"log"
	"time"

	"net/url"
	"strings"

	_ "github.com/lib/pq"
)

func InitDB(connectionString string) (*sql.DB, error) {
	// Check if connection string is URL encoded
	if strings.HasPrefix(connectionString, "postgres%") || strings.HasPrefix(connectionString, "postgresql%") {
		if decoded, err := url.PathUnescape(connectionString); err == nil {
			connectionString = decoded
		}
	}

	// Parse as URL to ensure sslmode is set (required for Supabase)
	if u, err := url.Parse(connectionString); err == nil && (u.Scheme == "postgres" || u.Scheme == "postgresql") {
		q := u.Query()
		if q.Get("sslmode") == "" {
			q.Set("sslmode", "require")
			u.RawQuery = q.Encode()
			connectionString = u.String()
		}
		log.Printf("Connecting to DB: %s", u.Redacted())
		log.Printf("Parsed Host: %s", u.Host)
		log.Printf("Parsed User: %s", u.User.Username())
	} else {
		log.Printf("Failed to parse URL or not a URL connection string: %v", err)
	}

	// Open database
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings (optional tapi recommended)
	db.SetMaxOpenConns(25)
	db.SetConnMaxLifetime(100 * time.Minute)
	db.SetMaxIdleConns(5)

	log.Println("Database connected successfully")
	return db, nil
}
