package service

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"main.go/internal/storage"
	"main.go/models"
)

var (
	ErrUserNotFound      = errors.New("user does not exist or password incorrect")
	ErrInvalidUserName   = errors.New("invalid username")
	ErrInvalidEmail      = errors.New("invalid email")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrPasswordDontMatch = errors.New("password didn't match")
	ErrUserExist         = errors.New("user exist")
)

type Auth interface {
	CreateUser(user models.User) error
	GenerateSessionToken(login, password string) (string, time.Time, error)
	ParseSessionToken(token string) (models.User, error)
	DeleteSessionToken(token string) error
}

type AuthService struct {
	storage storage.Auth
}

func newAuthService(storage storage.Auth) *AuthService {
	return &AuthService{
		storage: storage,
	}
}

func (s *AuthService) CreateUser(user models.User) error {
	var err error
	if _, err := s.storage.GetUserByLogin(user.Username); err == nil {
		return fmt.Errorf("service: create user: %w: username already used", ErrUserExist)
	}
	if _, err := s.storage.GetUserByEmail(user.Email); err == nil {
		return fmt.Errorf("service: create user: %w: email already used", ErrUserExist)
	}
	if err := validUser(user); err != nil {
		return fmt.Errorf("service: create user: %w", err)
	}

	user.Password, err = generateHashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("service: create user: %w", err)
	}
	return s.storage.CreateUser(user)
}

func (s *AuthService) GenerateSessionToken(username, password string) (string, time.Time, error) {
	user, err := s.storage.GetUserByLogin(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", time.Time{}, fmt.Errorf("service: generate session token: %w", ErrUserNotFound)
		}
		return "", time.Time{}, fmt.Errorf("service: generate session token: %w", err)
	}
	if err := compareHashAndPassword(user.Password, password); err != nil {
		return "", time.Time{}, fmt.Errorf("service: generate session token: %w", err)
	}
	token := uuid.NewString()
	expiresAt := time.Now().Add(time.Second * 300)
	if err := s.storage.SaveSessionToken(user.Username, token, expiresAt); err != nil {
		return "", time.Time{}, fmt.Errorf("service: generate session token: %w", err)
	}
	return token, expiresAt, nil
}

func (s *AuthService) ParseSessionToken(token string) (models.User, error) {
	user, err := s.storage.GetUserByToken(token)
	if err != nil {
		return user, fmt.Errorf("service: patse session token: %w", err)
	}
	return user, nil
}

func (s *AuthService) DeleteSessionToken(token string) error {
	err := s.storage.DeleteSessionToken(token)
	if err != nil {
		return fmt.Errorf("service: delete session token: %w", err)
	}
	return nil
}

func generateHashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func compareHashAndPassword(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return ErrUserNotFound
	}
	return nil
}

func validUser(user models.User) error {
	for _, char := range user.Username {
		if char <= 32 || char >= 127 {
			return ErrInvalidUserName
		}
	}
	validEmail, err := regexp.MatchString(`[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`, user.Email)
	if err != nil {
		return err
	}
	if !validEmail {
		return ErrInvalidEmail
	}
	if len(user.Username) <= 3 || len(user.Username) >= 36 {
		return ErrInvalidUserName
	}

	if user.Password != user.VerifyPassword {
		return ErrPasswordDontMatch
	}
	return nil
}
