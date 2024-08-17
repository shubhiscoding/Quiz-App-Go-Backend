package controllers

import (
    "encoding/json"
    "net/http"
    "quiz-backend/database"
    "quiz-backend/models"
)

func CreateGame(w http.ResponseWriter, r *http.Request) {
    var game models.Game
    json.NewDecoder(r.Body).Decode(&game)

    result, err := database.DB.Exec("INSERT INTO games (score, gamepoint) VALUES (?, ?)", game.Score, game.GamePoint)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    id, _ := result.LastInsertId()
    game.ID = int(id)

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
      

    json.NewEncoder(w).Encode(questions)
}
