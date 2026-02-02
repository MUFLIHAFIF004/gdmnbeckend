package models

import "time"

type Barang struct {
    ID            uint      `gorm:"primaryKey" json:"id"`
    KodeBarang    string    `gorm:"unique;not null" json:"kode_barang"`
    NamaBarang    string    `gorm:"not null" json:"nama_barang"`
    Stok          int       `gorm:"default:0" json:"stok"`
    Kategori      string    `json:"kategori"`
    Satuan        string    `json:"satuan"`
    Status        string    `json:"status"`
    Deskripsi     string    `json:"deskripsi"`
    Foto          string    `json:"foto"`
    TglKadaluarsa string `json:"tgl_kadaluarsa"`
    UpdatedAt     time.Time `json:"updated_at"`
}

type Riwayat struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    BarangID    uint      `json:"barang_id"`
    NamaBarang  string    `json:"nama_barang"`
    Tipe        string    `json:"tipe"` 
    Jumlah      int       `json:"jumlah"`
    Keterangan  string    `json:"keterangan"`
    CreatedAt   time.Time `json:"created_at"`
}