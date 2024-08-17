package controllers

import (
    "encoding/json"
    "net/http"
    "quiz-backend/models"
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
