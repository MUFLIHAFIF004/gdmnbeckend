package controllers

import (
	"encoding/json"
	"net/http"
	"gudangmng/config"
	"gudangmng/models"
)

// GetUserProfile: Mengambil data profile user yang sedang login
func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	// Mengambil ID dari parameter URL, misal: /profile?id=1
	id := r.URL.Query().Get("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "User tidak ditemukan"})
		return
	}

	// Keamanan: Password dikosongkan agar tidak terlihat di Flutter
	user.Password = ""
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateUserProfile: Logika untuk edit profile (nama, username, dll)
func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	var input models.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Format data salah"})
		return
	}

	// Pastikan ID user dikirim dari Flutter agar kita tahu siapa yang di-update
	if input.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "ID User diperlukan"})
		return
	}

	// Gunakan Map untuk update agar GORM tidak mengabaikan field kosong (opsional)
	// Tapi dengan Updates(input) sudah cukup untuk kebutuhan Tubes kamu
	if err := config.DB.Model(&input).Where("id = ?", input.ID).Updates(input).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Gagal update profile"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Profile berhasil diperbarui",
	})
}

// DeleteUserAccount: Tambahkan ini untuk jaga-jaga jika ada fitur hapus akun
func DeleteUserAccount(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if err := config.DB.Delete(&models.User{}, id).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Gagal hapus akun"})
		return
	}
	w.Write([]byte(`{"message":"Akun berhasil dihapus"}`))
}