package config

import (
	"fmt"
	"gudangmng/models" // 1. Pastikan kamu import package models kamu
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	godotenv.Load()

	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, dbname)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Gagal menyambungkan ke database: " + err.Error())
	}


	fmt.Println("Koneksi Database Berhasil!")

	// 2. TAMBAHKAN INI: Otomatis membuat tabel berdasarkan Struct di Go
	err = database.AutoMigrate(&models.User{}, &models.Token{},&models.Barang{}, &models.Riwayat{},) 
	if err != nil {
		fmt.Println("Gagal migrasi tabel:", err)
	} else {
		fmt.Println("Migrasi Database Berhasil!")
	}

	DB = database
}

func SetAccessControlHeaders(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return true
	}
	return false
}
