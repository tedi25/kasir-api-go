package repositories

import (
	"database/sql"
	"kasir-api/models"
	"time"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	// 1. Mulai Database Transaction (Wajib untuk Data Keuangan)
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	// Jaga-jaga: Rollback kalau ada error di tengah jalan
	defer tx.Rollback()

	// 2. Persiapan Variable
	var totalAmount int
	var details []models.TransactionDetail // Kita tampung data lengkapnya di sini

	// 3. Logic: Cek Harga & Hitung Subtotal (Sebelum Insert)
	for _, item := range items {
		var price int
		// Ambil harga produk dari tabel products berdasarkan ID
		err := tx.QueryRow("SELECT price FROM product WHERE id = $1", item.ProductID).Scan(&price)
		if err != nil {
			// Kalau produk tidak ditemukan atau error lain
			return nil, err
		}

		// Hitung Subtotal
		subtotal := price * item.Quantity
		totalAmount += subtotal

		// Masukkan ke slice details untuk dipakai nanti
		details = append(details, models.TransactionDetail{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Subtotal:  subtotal,
		})
	}

	// 4. INSERT HEADER TRANSAKSI (Tabel transactions)
	var transactionID int
	// Masukkan Total Belanja & Waktu sekarang. Mengembalikan ID transaksi baru.
	err = tx.QueryRow("INSERT INTO transactions (total_amount, created_at) VALUES ($1, $2) RETURNING id",
		totalAmount, time.Now()).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	// ==========================================
	// 5. INSERT DETAILS (INI BAGIAN TUGAS TASK SESSION 3)
	// ==========================================
	for i := range details {
		// Assign ID Transaksi yang baru dibuat ke setiap detail
		details[i].TransactionID = transactionID

		// Query Insert ke tabel transaction_details
		_, err = tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)",
			transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)

		if err != nil {
			return nil, err
		}
	}

	// 6. Commit (Simpan Permanen jika semua lancar)
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Kembalikan struct Transaction untuk response JSON
	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		CreatedAt:   time.Now(),
		Details:     details,
	}, nil
}

func (r *TransactionRepository) GetDailyReport() (models.SalesSummary, error) {
	var summary models.SalesSummary

	// 1. Query Total Revenue & Total Transaksi HARI INI
	// COALESCE(..., 0) gunanya biar kalau ga ada sales, hasilnya 0 (bukan error/null)
	// CURRENT_DATE adalah fungsi bawaan Postgres untuk ambil tanggal hari ini
	queryStats := `
        SELECT 
            COALESCE(SUM(total_amount), 0), 
            COUNT(id) 
        FROM transactions
        WHERE created_at::date = CURRENT_DATE
    `
	err := r.db.QueryRow(queryStats).Scan(&summary.TotalRevenue, &summary.TotalTransaksi)
	if err != nil {
		return summary, err
	}

	// 2. Query Produk Terlaris HARI INI
	// Kita perlu JOIN 3 tabel: transactions -> transaction_details -> products
	queryTopProduct := `
        SELECT 
            p.name, 
            SUM(td.quantity) as total_qty
        FROM transaction_details td
        JOIN product p ON td.product_id = p.id
        JOIN transactions t ON td.transaction_id = t.id
        WHERE t.created_at::date = CURRENT_DATE
        GROUP BY p.name
        ORDER BY total_qty DESC
        LIMIT 1
    `

	// Kita pakai Scan biasa. Kalau tidak ada penjualan sama sekali, ini mungkin error sql.ErrNoRows.
	// Jadi kita handle errornya.
	err = r.db.QueryRow(queryTopProduct).Scan(&summary.ProdukTerlaris.Nama, &summary.ProdukTerlaris.QtyTerjual)

	if err != nil {
		// Kalau errornya karena data kosong, kita biarkan kosong aja (jangan return error)
		summary.ProdukTerlaris.Nama = "-"
		summary.ProdukTerlaris.QtyTerjual = 0
	}

	return summary, nil
}
