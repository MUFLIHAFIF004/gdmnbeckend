package controllers

import (
    "encoding/json"
    "net/http"
    "gudangmng/config"
    "gudangmng/models"
)

// GetRiwayatHandler: Menampilkan log semua barang masuk & keluar
func GetRiwayatHandler(w http.ResponseWriter, r *http.Request) {
    var listRiwayat []models.Riwayat
    config.DB.Order("created_at desc").Find(&listRiwayat)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(listRiwayat)
}

// GetSummaryHandler: Untuk tampilan dashboard total barang masuk
func GetSummaryHandler(w http.ResponseWriter, r *http.Request) {
    var total int64
    config.DB.Model(&models.Riwayat{}).Where("tipe = ?", "MASUK").Select("sum(jumlah)").Row().Scan(&total)
    
    json.NewEncoder(w).Encode(map[string]interface{}{
        "total_barang_masuk": total,
    })
}

// DeleteRiwayatHandler: Untuk menghapus catatan riwayat tertentu
func DeleteRiwayatHandler(w http.ResponseWriter, r *http.Request) {
    // Ambil ID dari query string, misal: /riwayat/delete?id=1
    id := r.URL.Query().Get("id")
    
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

    w.Write([]byte(`{"message":"Catatan riwayat berhasil dihapus"}`))
}