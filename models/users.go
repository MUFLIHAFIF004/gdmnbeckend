package models

import "time"

type User struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Nama       string    `json:"nama"`        // Dari 'Nama Lengkap'
	Username   string    `gorm:"unique;not null" json:"username"`
	Email      string    `gorm:"unique" json:"email"` // Dari 'Email'
	Telepon    string    `json:"telepon"`     // Dari 'No. Telepon'
	IDKaryawan string    `json:"id_karyawan"` // Dari 'ID Karyawan'
	Password   string    `gorm:"not null" json:"password"`
	Foto       string    `gorm:"type:longtext" json:"foto"`
	CreatedAt  time.Time `json:"created_at"`
}

type Token struct {
	ID        uint      `gorm:"primaryKey"`
	Token     string    `gorm:"index"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}


// response and request structures

type RegisterInput struct {
	Nama       string `json:"nama"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Telepon    string `json:"telepon"`
	IDKaryawan string `json:"id_karyawan"`
	Password   string `json:"password"`
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