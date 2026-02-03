package controllers

import (
    "encoding/json"
    "net/http"
    "gudangmng/config"
    "gudangmng/models"
)

// GetRiwayatHandler: Sudah Aman
func GetRiwayatHandler(w http.ResponseWriter, r *http.Request) {
    var listRiwayat []models.Riwayat
    config.DB.Order("created_at desc").Find(&listRiwayat)
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(listRiwayat)
}

// GetSummaryHandler: Perbaikan pada penanganan nilai Kosong (Null)
func GetSummaryHandler(w http.ResponseWriter, r *http.Request) {
    var total int64
    err := config.DB.Model(&models.Riwayat{}).Where("tipe = ?", "MASUK").Select("COALESCE(sum(jumlah), 0)").Row().Scan(&total)
    
    if err != nil {
        total = 0 
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "total_barang_masuk": total,
    })
}

// DeleteRiwayatHandler: Perbaikan pada Header Response
func DeleteRiwayatHandler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    
    w.Header().Set("Content-Type", "application/json") // Tambahkan ini

    if id == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"message": "ID riwayat tidak ditemukan"})
        return
    }

    if err := config.DB.Delete(&models.Riwayat{}, id).Error; err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"message": "Gagal menghapus riwayat"})
        return
    }

    // Gunakan Encode agar format JSON konsisten dengan handler lainnya
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Catatan riwayat berhasil dihapus",
    })
}