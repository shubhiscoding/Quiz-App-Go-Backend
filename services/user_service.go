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

func (s *UserService) RegisterUser(name, email, password string) error {
    hashedPassword := hashPassword(password)

    query := `INSERT INTO users (name, email, password) VALUES (?, ?, ?)`
    _, err := database.DB.Exec(query, name, email, hashedPassword)
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

func hashPassword(password string) string {
    h := sha256.New()
    h.Write([]byte(password))
    return fmt.Sprintf("%x", h.Sum(nil))
}
