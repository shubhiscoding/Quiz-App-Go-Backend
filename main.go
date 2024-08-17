package main

import (
    "encoding/json"
    "log"
    "net/http"
    "quiz-backend/database"
    "quiz-backend/models"
    "quiz-backend/controllers"
    "crypto/sha256"
    "fmt"

    "github.com/gorilla/mux"
)

func main() {
    database.InitDB()
    defer database.DB.Close()

    r := mux.NewRouter()

    r.HandleFunc("/users", getUsers).Methods("GET")
    r.HandleFunc("/users", controllers.Register).Methods("POST")
    r.HandleFunc("/login", controllers.Login).Methods("POST")
    r.HandleFunc("/games", controllers.CreateGame).Methods("GET")
    r.HandleFunc("/user-games", createUserGame).Methods("POST")

    log.Fatal(http.ListenAndServe(":8080", r))
}

func getUsers(w http.ResponseWriter, r *http.Request) {
    rows, err := database.DB.Query("SELECT id, name, email, points FROM users")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var users []models.User
    for rows.Next() {
        var u models.User
        if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Points); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        users = append(users, u)
    }

    json.NewEncoder(w).Encode(users)
}

func createUserGame(w http.ResponseWriter, r *http.Request) {
    var userGame models.UserGame
    json.NewDecoder(r.Body).Decode(&userGame)

    _, err := database.DB.Exec("INSERT INTO user_games (user_id, game_id) VALUES (?, ?)", userGame.UserID, userGame.GameID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Update user points
    var gamePoint int
    err = database.DB.QueryRow("SELECT gamepoint FROM games WHERE id = ?", userGame.GameID).Scan(&gamePoint)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    _, err = database.DB.Exec("UPDATE users SET points = points + ? WHERE id = ?", gamePoint, userGame.UserID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(userGame)
}

func hashPassword(password string) string {
    h := sha256.New()
    h.Write([]byte(password))
    return fmt.Sprintf("%x", h.Sum(nil))
}
