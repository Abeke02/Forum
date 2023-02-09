package storage

import (
	"database/sql"

	"main.go/models"
)

type User interface {
	GetUserByUsername(username string) (models.User, error)
}

type UserStorage struct {
	db *sql.DB
}

func newUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (u *UserStorage) GetUserByUsername(username string) (models.User, error) {
	query := `SELECT id, email, username FROM user WHERE username=$1;`
	rows := u.db.QueryRow(query, username)
	var user models.User

	if err := rows.Scan(&user.ID, &user.Email, &user.Username); err != nil {
		return user, err
	}
	return user, nil
}
