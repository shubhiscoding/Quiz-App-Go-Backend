package controllers

import (
    "encoding/json"
    "net/http"
    "quiz-backend/models"
    "quiz-backend/services"
    "fmt"
)

var gameService = services.NewGameService()

func CreateGame(w http.ResponseWriter, r *http.Request) {
    var game models.Game
    if err := json.NewDecoder(r.Body).Decode(&game); err != nil {
        http.Error(w, fmt.Sprintf("CreateGame: invalid request: %v", err), http.StatusBadRequest)
        return
    }

    game, err := gameService.CreateGame(game)
    if err != nil {
        http.Error(w, fmt.Sprintf("CreateGame: failed to create game: %v", err), http.StatusInternalServerError)
        return
    }

    questions := []struct {
        QuestionText        string   `json:"questionText"`
        Options             []string `json:"options"`
        CorrectAnswerIndex  int      `json:"correctAnswerIndex"`
    }{
        {
            QuestionText:       "What is the capital of France?",
            Options:            []string{"London", "Berlin", "Paris", "Madrid"},
            CorrectAnswerIndex: 2,
        },
        {
            QuestionText:       "Which planet is known as the Red Planet?",
            Options:            []string{"Mars", "Jupiter", "Venus", "Saturn"},
            CorrectAnswerIndex: 0,
        },
        {
            QuestionText:       "What is the largest mammal in the world?",
            Options:            []string{"Elephant", "Blue Whale", "Giraffe", "Hippopotamus"},
            CorrectAnswerIndex: 1,
        },
    }

    retObj := struct {
        Game      models.Game `json:"game"`
        Questions []struct {
            QuestionText        string   `json:"questionText"`
            Options             []string `json:"options"`
            CorrectAnswerIndex  int      `json:"correctAnswerIndex"`
        } `json:"questions"`
    }{
        Game: game,
        Questions: questions,
    }

    json.NewEncoder(w).Encode(retObj)
}

func CreateUserGame(w http.ResponseWriter, r *http.Request) {
    var userGame models.UserGame
    if err := json.NewDecoder(r.Body).Decode(&userGame); err != nil {
        http.Error(w, fmt.Sprintf("CreateUserGame: invalid request: %v", err), http.StatusBadRequest)
        return
    }

    _, errr := gameService.GetGame(userGame.GameID)
    if errr != nil {
        http.Error(w, fmt.Sprintf("CreateUserGame: failed to get game: %v", errr), http.StatusInternalServerError)
        return
    }

    _, existErr := gameService.GetUserGame(userGame.GameID, userGame.UserID)
    if existErr == nil {
        http.Error(w, fmt.Sprintf("CreateUserGame: Game already exsits"), http.StatusInternalServerError)
        return
    }

    err := gameService.CreateUserGame(userGame)
    if err != nil {
        http.Error(w, fmt.Sprintf("CreateUserGame: failed to create user game: %v", err), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(userGame)
}

func EndGame(w http.ResponseWriter, r *http.Request) {
    var userGame struct {
        UserID    int `json:"user_id"`
        GameID    int `json:"game_id"`
        Score     int `json:"point"`
        GamePoint int `json:"gamePoint"`
    }

    if err := json.NewDecoder(r.Body).Decode(&userGame); err != nil {
        http.Error(w, fmt.Sprintf("EndGame: invalid request: %v", err), http.StatusBadRequest)
        return
    }

    _, err := gameService.GetUserGame(userGame.GameID, userGame.UserID)
    if err != nil {
        http.Error(w, fmt.Sprintf("EndGame: failed to get game: %v", err), http.StatusInternalServerError)
        return
    }

    updatedUserGame := models.UserGame{
        UserID:    userGame.UserID,
        GameID:    userGame.GameID,
        Score:     userGame.Score,
        GamePoint: userGame.GamePoint,
    }

    err = gameService.UpdateUserGame(updatedUserGame)
    if err != nil {
        http.Error(w, fmt.Sprintf("EndGame: failed to update user game: %v", err), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(userGame)
}
