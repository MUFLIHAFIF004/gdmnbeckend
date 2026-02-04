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

    // 3. Simpan Data Riwayat (Gunakan variabel 'Tanggal' yang baru, bukan TglKadaluarsa)
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
        Tanggal:    req.Tanggal, // PERBAIKAN: Sekarang pakai Tangsal Transaksi (tgl 04)
    }
    config.DB.Create(&riwayat)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Barang berhasil didaftarkan"})
}
// UpdateBarangHandler: Edit data master barang
func UpdateBarangHandler(w http.ResponseWriter, r *http.Request) {
	var input models.Barang
	json.NewDecoder(r.Body).Decode(&input)

	if err := config.DB.Model(&models.Barang{}).Where("id = ?", input.ID).Updates(map[string]interface{}{
		"kode_barang":    input.KodeBarang,
		"nama_barang":    input.NamaBarang,
		"stok":           input.Stok,
		"kategori":       input.Kategori,
		"satuan":         input.Satuan,
		"deskripsi":      input.Deskripsi,
		"foto":           input.Foto,
		"tgl_kadaluarsa": input.TglKadaluarsa,
	}).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		// MENGIRIM PESAN ERROR ASLI
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Gagal update data: " + err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Data barang diperbarui"})
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
