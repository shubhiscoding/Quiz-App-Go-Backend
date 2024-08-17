package models

type UserGame struct {
    UserID int `json:"user_id"`
    GameID int `json:"game_id"`
    Score  int `json:"point"`
    GamePoint int `json:"gamepoint"`
}