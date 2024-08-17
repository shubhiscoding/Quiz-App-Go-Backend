package models

type Game struct {
    ID     int `json:"id"`
    Score  int `json:"score"`
    GamePoint int `json:"points"`
}