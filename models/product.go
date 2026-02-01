package models

type Product struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Price      int    `json:"price"`
	Stock      int    `json:"stock"`
	CategoryID int    `json:"category_id"`
	// Field khusus untuk Output (Join Result)
	// "omitempty" artinya kalau kosong tidak usah ditampilkan di JSON
	CategoryName string `json:"category_name,omitempty"`
}
