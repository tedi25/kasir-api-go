package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
	// Import layer-layer aplikasi Anda
)

// LANGKAH 2: Define Struct Config
// Ditaruh di luar function main agar jadi tipe data global
type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONNECTION"`
}

func main() {
	// LANGKAH 4: Setup Viper & Baca .env
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Cek apakah file .env ada menggunakan os.Stat
	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig() // underscore (_) artinya errornya diabaikan/tidak ditampung
	}

	// LANGKAH 5: Assign ke Struct
	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONNECTION"),
	}

	// Jaga-jaga jika PORT kosong (opsional, tapi disarankan)
	if config.Port == "" {
		config.Port = "8080"
	}

	// LANGKAH 6: Database Connection
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		fmt.Println("gagal konek ke database:", err)
		return
	}
	defer db.Close()

	// Product Wiring
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Category Wiring
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Transaction
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// ==========================================
	// ROUTING APP
	// ==========================================

	// Route Produk
	http.HandleFunc("/api/produk", productHandler.HandleProducts)
	http.HandleFunc("/api/produk/", productHandler.HandleProductByID)

	// Route Categories
	http.HandleFunc("/api/kategori", categoryHandler.HandleCategories)
	http.HandleFunc("/api/kategori/", categoryHandler.HandleCategoryByID)

	// ROute Transactions
	http.HandleFunc("/api/checkout", transactionHandler.HandleCheckout)

	// Route Report Harian
	http.HandleFunc("/api/report/hari-ini", transactionHandler.HandleDailyReport)

	// Route Health Check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	// LANGKAH 6: Implement ke Server
	addr := "0.0.0.0:" + config.Port

	fmt.Println("Server running di", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("gagal running server", err)
	}
}
