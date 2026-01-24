package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	// --- ROUTING ---

	// 1. Route untuk List & Create
	// Memanggil fungsi handleProduk yang ada di file produk.go
	http.HandleFunc("/api/produk", handleProduk)

	// 2. Route untuk Detail (GetByID, Update, Delete)
	// Perhatikan tanda slash "/" di akhir, artinya menangkap sub-path (misal /api/produk/1)
	http.HandleFunc("/api/produk/", handleProdukDetail)

	// Route untuk Category
	http.HandleFunc("/categories", handleCategories)
	http.HandleFunc("/categories/", handleCategories)

	// 3. Health Check (Generic, boleh ditaruh di sini)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	// --- SERVER ---
	fmt.Println("Server running di localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("gagal running server:", err)
	}
}
