package routes

import (
	"gudangmng/config"
	"gudangmng/controllers"
	"net/http"
)

func URL(w http.ResponseWriter, r *http.Request) {
	if config.SetAccessControlHeaders(w, r) {
		return
	}

	var method, path string = r.Method, r.URL.Path

	switch {
	// === AUTHENTICATION ===
	case method == "POST" && path == "/auth/register":
		controllers.RegisterHandler(w, r)
	case method == "POST" && path == "/auth/login":
		controllers.LoginHandler(w, r)
	case method == "POST" && path == "/auth/logout":
		controllers.LogoutHandler(w, r)

	// === MANAJEMEN BARANG (CRUD) ===
	case method == "GET" && path == "/barang/all":
		controllers.GetBarangHandler(w, r)
	case method == "POST" && path == "/barang/create":
		controllers.InputBarangHandler(w, r)
	case method == "PUT" && path == "/barang/update":
		controllers.UpdateBarangHandler(w, r)
	case method == "DELETE" && path == "/barang/delete":
		controllers.DeleteBarangHandler(w, r)

	// === LOGISTIK & STOK (Fitur Utama) ===
	case method == "POST" && path == "/stok/update":
		controllers.UpdateStokHandler(w, r)

	// === RIWAYAT  ===
	case method == "GET" && path == "/riwayat/all":
		controllers.GetRiwayatHandler(w, r)
	case method == "GET" && path == "/riwayat/summary":
		controllers.GetSummaryHandler(w, r)
	case method == "DELETE" && path == "/riwayat/delete":
		controllers.DeleteRiwayatHandler(w, r)

	// === PROFILE (Opsional) ===
	case method == "GET" && path == "/profile":
		controllers.GetUserProfile(w, r)
	case method == "PUT" && path == "/profile/update":
		controllers.UpdateUserProfile(w, r)
		// controllers.UpdateProfileHandler(w, r)

	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message":"Endpoint tidak ditemukan"}`))
	}
}
