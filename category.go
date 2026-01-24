package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var categories = []Category{
	{ID: 1, Name: "Makanan"},
	{ID: 2, Name: "Minuman"},
}

// Handler utama untuk routing /categories dan /categories/{id}
func handleCategories(w http.ResponseWriter, r *http.Request) {
	// Cek apakah ada ID di URL (misal: /categories/1)
	path := strings.TrimPrefix(r.URL.Path, "/categories")

	// Jika path kosong atau cuma "/", berarti ambil SEMUA atau POST baru
	if path == "" || path == "/" {
		switch r.Method {
		case "GET":
			getAllCategories(w, r)
		case "POST":
			createCategory(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// Jika ada ID (path tidak kosong), lanjut ke logic by ID
	// Hapus slash di depan jika ada (misal "/1" jadi "1")
	pathID := strings.TrimPrefix(path, "/")
	id, err := strconv.Atoi(pathID)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		getCategoryByID(w, r, id)
	case "PUT":
		updateCategory(w, r, id)
	case "DELETE":
		deleteCategory(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getAllCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func createCategory(w http.ResponseWriter, r *http.Request) {
	var newCat Category
	if err := json.NewDecoder(r.Body).Decode(&newCat); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	newCat.ID = len(categories) + 1 // Auto ID sederhana
	categories = append(categories, newCat)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCat)
}

func getCategoryByID(w http.ResponseWriter, r *http.Request, id int) {
	for _, c := range categories {
		if c.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(c)
			return
		}
	}
	http.Error(w, "Category not found", http.StatusNotFound)
}

func updateCategory(w http.ResponseWriter, r *http.Request, id int) {
	var updateData Category
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	for i, c := range categories {
		if c.ID == id {
			categories[i].Name = updateData.Name // Update nama saja

			// Return object yang sudah diupdate
			categories[i].ID = id
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(categories[i])
			return
		}
	}
	http.Error(w, "Category not found", http.StatusNotFound)
}

func deleteCategory(w http.ResponseWriter, r *http.Request, id int) {
	for i, c := range categories {
		if c.ID == id {
			categories = append(categories[:i], categories[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "Category deleted"})
			return
		}
	}
	http.Error(w, "Category not found", http.StatusNotFound)
}
