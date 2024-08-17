package database

import (
    "database/sql"
    "log"

    _ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
    var err error
    DB, err = sql.Open("sqlite3", "./quiz.db")
    if err != nil {
        log.Fatal(err)
    }

    createTables()
}

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
        score INTEGER,
        gamepoint INTEGER
    );`

    createUserGamesTable := `
    CREATE TABLE IF NOT EXISTS user_games (
        user_id INTEGER,
        game_id INTEGER,
        PRIMARY KEY (user_id, game_id),
        FOREIGN KEY (user_id) REFERENCES users (id),
        FOREIGN KEY (game_id) REFERENCES games (id)
    );`

    _, err := DB.Exec(createUsersTable)
    if err != nil {
        log.Fatal(err)
    }

    _, err = DB.Exec(createGamesTable)
    if err != nil {
        log.Fatal(err)
    }

    _, err = DB.Exec(createUserGamesTable)
    if err != nil {
        log.Fatal(err)
    }
}