package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Password  string    `gorm:"not null" json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type Token struct {
	ID        uint      `gorm:"primaryKey"`
	Token     string    `gorm:"index"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}


// response and request structures

type RegisterInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token"`
	User    User   `json:"user"`
}