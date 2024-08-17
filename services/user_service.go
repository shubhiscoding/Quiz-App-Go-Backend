package services

import (
    "crypto/sha256"
    "database/sql"
    "errors"
    "fmt"
    "quiz-backend/database"
    "quiz-backend/models"
)

type UserService struct{}

func NewUserService() *UserService {
    return &UserService{}
}

func (s *UserService) RegisterUser(name, email, password string) error {
    var existingUser models.User

    // Check if the user already exists
    query := `SELECT id FROM users WHERE email = ?`
    err := database.DB.QueryRow(query, email).Scan(&existingUser.ID)
    if err != nil && err != sql.ErrNoRows {
        return err
    }
    if existingUser.ID != 0 {
        return errors.New("user already exists")
    }
    query = `INSERT INTO users (name, email, password) VALUES (?, ?, ?)`
    _, err = database.DB.Exec(query, name, email, hashPassword(password))
    if err != nil {
        return err
    }

    return nil
}

func (s *UserService) AuthenticateUser(email, password string) (*models.User, error) {
    hashedPassword := hashPassword(password)

    query := `SELECT id, name, email, points FROM users WHERE email = ? AND password = ?`
    row := database.DB.QueryRow(query, email, hashedPassword)

    var user models.User
    if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Points); err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("invalid email or password")
        }
        return nil, err
    }

    return &user, nil
}

func (s *UserService) GetUsers() ([]models.User, error) {
    rows, err := database.DB.Query("SELECT id, name, email, points FROM users")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []models.User
    for rows.Next() {
        var u models.User
        if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Points); err != nil {
            return nil, err
        }
        users = append(users, u)
    }

    return users, nil
}

func (s *UserService) GetUserByID(userID int) (models.User, error) {
	var user models.User
	err := database.DB.QueryRow("SELECT id, name, email, points FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Name, &user.Email, &user.Points)
	if err != nil {
		return models.User{}, fmt.Errorf("GetUserById: failed to retrieve user: %v", err)
	}
	return user, nil
}

func (s *UserService) GetUserGames(userID int) ([]models.UserGame, error) {
    var userGame []models.UserGame
    query := `SELECT user_id, game_id, point, gamepoint FROM user_games WHERE user_id = ?`
    rows, err := database.DB.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    for rows.Next() {
        var ug models.UserGame
        if err := rows.Scan(&ug.UserID, &ug.GameID, &ug.Score, &ug.GamePoint); err != nil {
            return nil, err
        }
        userGame = append(userGame, ug)
    }
    
    if err := rows.Err(); err != nil {
        return nil, err
    }
    
    return userGame, nil
    
}


func hashPassword(password string) string {
    h := sha256.New()
    h.Write([]byte(password))
    return fmt.Sprintf("%x", h.Sum(nil))
}