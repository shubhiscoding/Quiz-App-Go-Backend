package controllers

import (
    "encoding/json"
    "net/http"
    "quiz-backend/database"
    "quiz-backend/models"
    "crypto/sha256"
    "fmt"
)

func Register(w http.ResponseWriter, r *http.Request) {
    var user models.User
    json.NewDecoder(r.Body).Decode(&user)

    hashedPassword := hashPassword(user.Password)
    if len(user.Password) < 6 {
        http.Error(w, "Password must be at least 6 characters long", http.StatusBadRequest)
        return
    }
    user.Password = hashedPassword

    result, err := database.DB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", user.Name, user.Email, user.Password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    id, _ := result.LastInsertId()
    user.ID = int(id)
    user.Password = "" // Don't return the password in the response

    json.NewEncoder(w).Encode(user)
}

func Login(w http.ResponseWriter, r *http.Request) {
    var loginData struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    json.NewDecoder(r.Body).Decode(&loginData)

    var user models.User
    err := database.DB.QueryRow("SELECT id, name, email, password, points FROM users WHERE email = ?", loginData.Email).Scan(
        &user.ID, &user.Name, &user.Email, &user.Password, &user.Points,
    )
    if err != nil {
        http.Error(w, "Invalid email or password", http.StatusUnauthorized)
        return
    }

    hashedPassword := hashPassword(loginData.Password)
    if user.Password != hashedPassword {
        http.Error(w, "Invalid email or password", http.StatusUnauthorized)
        return
    }

    user.Password = "" // Don't return the password in the response

    json.NewEncoder(w).Encode(user)
}

func hashPassword(password string) string {
    h := sha256.New()
    h.Write([]byte(password))
    return fmt.Sprintf("%x", h.Sum(nil))
}
