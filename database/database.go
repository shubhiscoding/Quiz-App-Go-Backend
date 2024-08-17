package database

import (
    "database/sql"
    "log"

    _ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB initializes the database connection and creates the necessary tables.
func InitDB() {
    var err error
    DB, err = sql.Open("sqlite3", "./quiz.db")
    if err != nil {
        log.Fatalf("Failed to connect to the database: %v", err)
    }

    // Verify the connection is valid
    if err := DB.Ping(); err != nil {
        log.Fatalf("Failed to establish a database connection: %v", err)
    }

    createTables()
}

// createTables creates the necessary tables if they do not already exist.
func createTables() {
    createUsersTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        email TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL,
        points INTEGER DEFAULT 0
    );`

    createGamesTable := `
    CREATE TABLE IF NOT EXISTS games (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        level INTEGER NOT NULL
    );`

    createUserGamesTable := `
    CREATE TABLE IF NOT EXISTS user_games (
        user_id INTEGER,
        game_id INTEGER,
        point INTEGER DEFAULT 0,
        gamepoint INTEGER DEFAULT 0,
        PRIMARY KEY (user_id, game_id),
        FOREIGN KEY (user_id) REFERENCES users (id),
        FOREIGN KEY (game_id) REFERENCES games (id)
    );`

    tables := []string{createUsersTable, createGamesTable, createUserGamesTable}

    for _, table := range tables {
        if _, err := DB.Exec(table); err != nil {
            log.Fatalf("Failed to create table: %v", err)
        }
    }
}
