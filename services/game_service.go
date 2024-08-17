package services

import (
    "quiz-backend/database"
    "quiz-backend/models"
	"fmt"
)

type GameService struct{}

func NewGameService() *GameService {
	return &GameService{}
}

func (s *GameService) CreateGame(game models.Game) (models.Game, error) {
	result, err := database.DB.Exec("INSERT INTO games (level) VALUES (?)", game.Level)
	if err != nil {
        return models.Game{}, fmt.Errorf("CreateGame: failed to insert game: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Game{}, fmt.Errorf("CreateGame: failed to retrieve last insert ID: %v", err)
	}

	game.ID = int(id)
	return game, nil
}

func (s *GameService) CreateUserGame(userGame models.UserGame) error {
	_, err := database.DB.Exec("INSERT INTO user_games (user_id, game_id, point, gamepoint) VALUES (?, ?, ?, ?)", userGame.UserID, userGame.GameID, userGame.Score, userGame.GamePoint)
	if err != nil {
		return fmt.Errorf("CreateUserGame: failed to insert user game: %v", err)
	}
	return nil
}

func (s *GameService) GetUserGame(gameID int, userID int) (models.Game, error) {
	var game models.Game
	err := database.DB.QueryRow("SELECT g.id, g.level FROM games g JOIN user_games u ON g.id = u.game_id WHERE g.id = ? and u.user_id = ?", gameID, userID).Scan(&game.ID, &game.Level)
	if err != nil {
		return models.Game{}, fmt.Errorf("GetGame: failed to retrieve game: %v", err)
	}
	return game, nil
}

func (s *GameService) GetGame(gameID int) (models.Game, error) {
	var game models.Game
	err := database.DB.QueryRow("SELECT g.id, g.level FROM games g WHERE g.id = ?", gameID).Scan(&game.ID, &game.Level)
	if err != nil {
		return models.Game{}, fmt.Errorf("GetGame: failed to retrieve game: %v", err)
	}
	return game, nil
}

func (s *GameService) UpdateUserGame(game models.UserGame) error {
	var currentPoints int
	err := database.DB.QueryRow("SELECT points FROM users WHERE id = ?", game.UserID).Scan(&currentPoints)
	if err != nil {
		return fmt.Errorf("UpdateUserGame: failed to retrieve user points: %v", err)
	}

	_, err = database.DB.Exec("UPDATE users SET points = ? WHERE id = ?", game.Score + currentPoints, game.UserID)
	if err != nil {
		return fmt.Errorf("UpdateUserGame: failed to update user points: %v", err)
	}
	
	_, err = database.DB.Exec("UPDATE user_games SET point = ?, gamepoint = ? WHERE user_id = ? AND game_id = ?", game.Score + currentPoints, game.GamePoint, game.UserID, game.GameID)
	if err != nil {
		return fmt.Errorf("UpdateUserGame: failed to update user game: %v", err)
	}

	return nil
}

