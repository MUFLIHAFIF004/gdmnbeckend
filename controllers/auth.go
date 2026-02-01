package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"gudangmng/config" 
	"gudangmng/models"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("RAHASIA_GUDANG_714230012")

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Data tidak valid", http.StatusBadRequest)
		return
	}

	newUser := models.User{
		Username:  req.Username,
		Password:  req.Password,
		CreatedAt: time.Now(),
	}

	if err := config.DB.Create(&newUser).Error; err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Gagal daftar atau username sudah ada"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Registrasi berhasil"})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginInput
	json.NewDecoder(r.Body).Decode(&req)

	var user models.User
	err := config.DB.Where("username = ? AND password = ?", req.Username, req.Password).First(&user).Error
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Username atau Password salah"})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"username": user.Username,
		"exp":      expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)

	newToken := models.Token{
		Token:     tokenString,
		Username:  user.Username,
		CreatedAt: time.Now(),
	}
	config.DB.Create(&newToken)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.LoginResponse{
		Status:  true,
		Message: "Login Berhasil",
		Token:   tokenString,
		User:    user,
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tokenHeader := r.Header.Get("Authorization")
	if tokenHeader == "" {
		http.Error(w, "Token diperlukan", http.StatusUnauthorized)
		return
	}

	config.DB.Where("token = ?", tokenHeader).Delete(&models.Token{})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logout berhasil"})
}