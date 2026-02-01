package controllers

import (
	"encoding/json"
	"net/http"
	"gudangmng/config"
	"gudangmng/models"
)

// InputBarangHandler: Menambah barang baru ke sistem (First Entry)
func InputBarangHandler(w http.ResponseWriter, r *http.Request) {
	var input models.Barang
	json.NewDecoder(r.Body).Decode(&input)

	if err := config.DB.Create(&input).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Gagal input barang"})
		return
	}

	// TAMBAHAN: Catat riwayat awal agar saat pertama input tidak kosong
	riwayat := models.Riwayat{
		BarangID:   input.ID,
		NamaBarang: input.NamaBarang,
		Tipe:       "MASUK",
		Jumlah:     input.Stok,
		Keterangan: "Pendaftaran Barang Baru (Stok Awal)",
	}
	config.DB.Create(&riwayat)

	w.Write([]byte(`{"message":"Barang berhasil didaftarkan"}`))
}

// GetBarangHandler & DeleteBarangHandler tetap sama...
func GetBarangHandler(w http.ResponseWriter, r *http.Request) {
	var barangs []models.Barang
	config.DB.Find(&barangs)
	json.NewEncoder(w).Encode(barangs)
}

// UpdateBarangHandler: Edit data barang (untuk perbaikan human error)
func UpdateBarangHandler(w http.ResponseWriter, r *http.Request) {
	var input models.Barang
	json.NewDecoder(r.Body).Decode(&input)

	if err := config.DB.Model(&input).Where("id = ?", input.ID).Updates(input).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Gagal update data"})
		return
	}

	// MODIFIKASI: Sinkronkan riwayat lama agar datanya ikut berubah (tidak double)
	config.DB.Model(&models.Riwayat{}).
		Where("barang_id = ? AND keterangan LIKE ?", input.ID, "%Pendaftaran%").
		Updates(models.Riwayat{
			NamaBarang: input.NamaBarang,
			Jumlah:     input.Stok,
		})

	w.Write([]byte(`{"message":"Data barang diperbarui"}`))
}

// DeleteBarangHandler tetap sama
func DeleteBarangHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if err := config.DB.Delete(&models.Barang{}, id).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Gagal hapus barang"})
		return
	}
	w.Write([]byte(`{"message":"Barang berhasil dihapus"}`))
}

// UpdateStokHandler: Logika Masuk/Keluar & Pencatatan Otomatis ke Riwayat
func UpdateStokHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID         uint   `json:"id"`
		Jumlah     int    `json:"jumlah"`
		Tipe       string `json:"tipe"` // "MASUK" atau "KELUAR"
		Keterangan string `json:"keterangan"`
	}
	json.NewDecoder(r.Body).Decode(&input)

	var barang models.Barang
	config.DB.First(&barang, input.ID)

	if input.Tipe == "MASUK" {
		barang.Stok += input.Jumlah
	} else {
		if barang.Stok < input.Jumlah {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Stok tidak cukup!"})
			return
		}
		barang.Stok -= input.Jumlah
	}
	barang.Status = input.Tipe
	config.DB.Save(&barang)

	riwayat := models.Riwayat{
		BarangID:   barang.ID,
		NamaBarang: barang.NamaBarang,
		Tipe:       input.Tipe,
		Jumlah:     input.Jumlah,
		Keterangan: input.Keterangan,
	}
	config.DB.Create(&riwayat)
	w.Write([]byte(`{"message":"Mutasi stok berhasil dicatat"}`))
}