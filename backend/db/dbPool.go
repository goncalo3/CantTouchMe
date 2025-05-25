package db

import (
	"backend/config"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db   *sql.DB
	once sync.Once
)

// InitDB initializes the database connection pool.
func InitDB(cfg *config.DbConfig) {
	once.Do(func() {
		var err error
		// Use cfg.DbHost and cfg.DbName for the database connection
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", cfg.DbUser, cfg.DbPwd, cfg.DbHost, cfg.Port, cfg.DbName)
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("Error opening database: %v", err)
		}

		db.SetConnMaxLifetime(time.Minute * 3)
		db.SetMaxOpenConns(50)
		db.SetMaxIdleConns(50)

		err = db.Ping()
		if err != nil {
			log.Fatalf("Error connecting to database: %v", err)
		}
		fmt.Println("Database connection pool initialized successfully")
	})
}

// GetDB returns the database connection pool.
func GetDB() *sql.DB {
	if db == nil {
		log.Fatal("Database connection pool not initialized. Call InitDB first.")
	}
	return db
}

// CloseDB closes the database connection pool.
func CloseDB() {
	if db != nil {
		db.Close()
		fmt.Println("Database connection pool closed.")
	}
}
