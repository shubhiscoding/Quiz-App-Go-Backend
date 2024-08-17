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
        Explanation         string   `json:"explanation"`
    }{
        {
            QuestionText:       "What is a stock?",
            Options:            []string{"A type of bond", "A share in the ownership of a company", "A loan given to the government", "A type of commodity"},
            CorrectAnswerIndex: 1,
            Explanation:        "A stock represents a share in the ownership of a company. When you buy a stock, you own a part of that company and are entitled to a portion of its profits.",
        },
        {
            QuestionText:       "What is a bear market?",
            Options:            []string{"A market characterized by rising prices", "A market characterized by falling prices", "A market with no significant changes", "A market that is closed for trading"},
            CorrectAnswerIndex: 1,
            Explanation:        "A bear market is characterized by falling prices, typically in the context of securities or commodities markets. It indicates widespread pessimism and a negative sentiment among investors.",
        },
        {
            QuestionText:       "What is the primary purpose of an Initial Public Offering (IPO)?",
            Options:            []string{"To raise capital for the company", "To pay off company debt", "To distribute profits to shareholders", "To buy back shares from investors"},
            CorrectAnswerIndex: 0,
            Explanation:        "The primary purpose of an IPO is to raise capital for the company. By going public, a company can raise funds from a large pool of investors to finance its growth and operations.",
        },
        {
            QuestionText:       "Which of the following is a type of derivative instrument?",
            Options:            []string{"Stocks", "Bonds", "Options", "Real Estate"},
            CorrectAnswerIndex: 2,
            Explanation:        "Options are a type of derivative instrument. Derivatives are financial contracts whose value is derived from the performance of underlying assets, indices, or interest rates.",
        },
        {
            QuestionText:       "What is market capitalization?",
            Options:            []string{"The total value of a company's outstanding bonds", "The total amount of profit a company has earned", "The total value of a company's outstanding shares", "The total amount of debt a company has"},
            CorrectAnswerIndex: 2,
            Explanation:        "Market capitalization is the total value of a company's outstanding shares of stock. It is calculated by multiplying the current market price of one share by the total number of shares outstanding.",
        },
    }
    

    retObj := struct {
        Game      models.Game `json:"game"`
        Questions []struct {
            QuestionText        string   `json:"questionText"`
            Options             []string `json:"options"`
            CorrectAnswerIndex  int      `json:"correctAnswerIndex"`
            Explanation         string   `json:"explanation"`
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
