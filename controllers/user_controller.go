package controllers

import (
    "encoding/json"
    "net/http"
    "quiz-backend/models"
    "strconv"
    "quiz-backend/services"
    "fmt"
)

var userService = services.NewUserService()

func Register(w http.ResponseWriter, r *http.Request) {
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, fmt.Sprintf("Register: invalid request: %v", err), http.StatusBadRequest)
        return
    }

    if len(user.Password) < 6 {
        http.Error(w, "Register: password must be at least 6 characters long", http.StatusBadRequest)
        return
    }

    err := userService.RegisterUser(user.Name, user.Email, user.Password)
    if err != nil {
        if err.Error() == "user already exists" {
            http.Error(w, "Register: user already exists", http.StatusConflict)
            return
        }
        http.Error(w, fmt.Sprintf("Register: failed to register user: %v", err), http.StatusInternalServerError)
        return
    }

    user.Points = 0
    user.Password = ""
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func Login(w http.ResponseWriter, r *http.Request) {
    var loginData struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
        http.Error(w, fmt.Sprintf("Login: invalid request: %v", err), http.StatusBadRequest)
        return
    }

    user, err := userService.AuthenticateUser(loginData.Email, loginData.Password)
    if err != nil {
        if err.Error() == "invalid email or password" {
            http.Error(w, "Login: invalid email or password", http.StatusUnauthorized)
            return
        }
        http.Error(w, fmt.Sprintf("Login: failed to authenticate user: %v", err), http.StatusInternalServerError)
        return
    }

    user.Password = ""
    json.NewEncoder(w).Encode(user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
    users, err := userService.GetUsers()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(users)
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
    // Extract the user ID from the URL parameters
    idStr := r.URL.Query().Get("id")
    if idStr == "" {
        http.Error(w, "User ID is required", http.StatusBadRequest)
        return
    }

    // Convert the ID to an integer
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    // Fetch the user from the service
    user, err := userService.GetUserByID(id)
    if err != nil {
        if err.Error() == "user not found" {
            http.Error(w, "User not found", http.StatusNotFound)
        } else {
            http.Error(w, fmt.Sprintf("Failed to get user: %v", err), http.StatusInternalServerError)
        }
        return
    }

    // Return the user as JSON
    json.NewEncoder(w).Encode(user)
}