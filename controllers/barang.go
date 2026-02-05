package controllers

import (
	"encoding/json"
	"gudangmng/config"
	"gudangmng/models"
	"net/http"
)

// GetBarangHandler: Mengambil semua daftar barang
func GetBarangHandler(w http.ResponseWriter, r *http.Request) {
	var barangs []models.Barang
	config.DB.Find(&barangs)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(barangs)
}

// InputBarangHandler: Menambah barang baru ke sistem
func InputBarangHandler(w http.ResponseWriter, r *http.Request) {
    // 1. Buat struct sementara untuk menangkap "tanggal" transaksi dari Flutter
    var req struct {
        models.Barang
        Tanggal string `json:"tanggal"` // Ini untuk menampung transDateStr (tgl 04) dari Flutter
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"message": "Format data tidak valid"})
        return
    }

    // 2. Simpan Data Master Barang (Ini akan menyimpan tgl_kadaluarsa tgl 28)
    if err := config.DB.Create(&req.Barang).Error; err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"message": "Gagal input barang"})
        return
    }
    ketRiwayat := req.Deskripsi
    if ketRiwayat == "" {
        ketRiwayat = "Pendaftaran Barang Baru"
    }

    riwayat := models.Riwayat{
        BarangID:   req.Barang.ID,
        NamaBarang: req.Barang.NamaBarang,
        Tipe:       "MASUK",
        Jumlah:     req.Barang.Stok,
        Keterangan: ketRiwayat + " (MASUK)",
        Tanggal:    req.Tanggal, 
    }
    config.DB.Create(&riwayat)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Barang berhasil didaftarkan"})
}
// UpdateBarangHandler: Edit data master barang
func UpdateBarangHandler(w http.ResponseWriter, r *http.Request) {
    // 1. Gunakan struct penampung agar bisa menangkap "tanggal" (masuk) dan "tgl_kadaluarsa"
    var req struct {
        models.Barang
        Tanggal string `json:"tanggal"` // Menangkap transDateStr baru dari Flutter
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    // 2. UPDATE TABEL BARANG (Master Data & Expired)
    if err := config.DB.Model(&models.Barang{}).Where("id = ?", req.Barang.ID).Updates(map[string]interface{}{
        "kode_barang":    req.Barang.KodeBarang,
        "nama_barang":    req.Barang.NamaBarang,
        "stok":           req.Barang.Stok,
        "kategori":       req.Barang.Kategori,
        "satuan":         req.Barang.Satuan,
        "deskripsi":      req.Barang.Deskripsi,
        "foto":           req.Barang.Foto,
        "tgl_kadaluarsa": req.Barang.TglKadaluarsa, // Tgl Expired baru
    }).Error; err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    // 3. UPDATE TABEL RIWAYAT (Tanggal Masuk)
    // Cari riwayat bertipe 'MASUK' yang paling pertama milik barang ini, lalu ubah tanggalnya
    if req.Tanggal != "" {
        config.DB.Model(&models.Riwayat{}).
            Where("barang_id = ? AND tipe = ?", req.Barang.ID, "MASUK").
            Order("created_at asc"). // Ambil yang paling pertama dibuat
            Limit(1).
            Update("tanggal", req.Tanggal) // Update ke Tgl Masuk baru
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Data barang dan riwayat diperbarui"})
}

// DeleteBarangHandler: Menghapus data barang
func DeleteBarangHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	// Hapus barang berdasarkan ID
	if err := config.DB.Delete(&models.Barang{}, id).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Gagal hapus barang"})
		return
	}

	// Opsional: Hapus juga riwayat yang terkait agar database bersih
	config.DB.Where("barang_id = ?", id).Delete(&models.Riwayat{})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Barang dan riwayat berhasil dihapus"})
}

// UpdateStokHandler: Logika Masuk/Keluar stok
func UpdateStokHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID         uint   `json:"id"`
		Jumlah     int    `json:"jumlah"`
		Tipe       string `json:"tipe"`
		Keterangan string `json:"keterangan"`
		Tanggal    string `json:"tanggal"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Format data tidak valid"})
		return
	}

	var barang models.Barang
	if err := config.DB.First(&barang, req.ID).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Barang tidak ditemukan"})
		return
	}

	// LOGIKA MUTASI
	barang.Status = req.Tipe
	config.DB.Save(&barang)

	finalKeterangan := req.Keterangan + " (" + req.Tipe + ")"

	config.DB.Create(&models.Riwayat{
		BarangID:   barang.ID,
		NamaBarang: barang.NamaBarang,
		Tipe:       req.Tipe,
		Jumlah:     req.Jumlah,
		Keterangan: finalKeterangan,
		Tanggal:    req.Tanggal,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Status barang berhasil diupdate menjadi " + req.Tipe,
	})
}
