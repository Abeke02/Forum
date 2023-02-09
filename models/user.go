package models

import "time"

type User struct {
	ID             int
	Email          string
	Username       string
	Password       string
	VerifyPassword string
	Token          string
	ExpiresAt      time.Time
}
