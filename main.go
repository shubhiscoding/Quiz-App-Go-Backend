package main

import (
    "log"
    "net/http"
    "quiz-backend/database"
    "quiz-backend/controllers"
    "github.com/gorilla/mux"
)

func main() {
    database.InitDB()
    defer database.DB.Close()

    r := mux.NewRouter()

    r.HandleFunc("/users", controllers.GetUsers).Methods("GET")
    r.HandleFunc("/users", controllers.Register).Methods("POST")
    r.HandleFunc("/login", controllers.Login).Methods("POST")
    r.HandleFunc("/games", controllers.CreateGame).Methods("POST") 
    r.HandleFunc("/games", controllers.GetUsersGame).Methods("GET")   
    r.HandleFunc("/user-games", controllers.CreateUserGame).Methods("POST")
    r.HandleFunc("/users/user", controllers.GetUserByID).Methods("GET")
    r.HandleFunc("/user", controllers.GetUserByID).Methods("POST")
    r.HandleFunc("/game-end", controllers.EndGame).Methods("POST")

    log.Fatal(http.ListenAndServe(":8080", r))
}