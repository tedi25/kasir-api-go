package repositories

import (
	"database/sql"
	"kasir-api/models"
)

type CategoryRepository struct {
	DB *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	rows, err := r.DB.Query("SELECT id, name FROM category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *CategoryRepository) GetByID(id int) (*models.Category, error) {
	var c models.Category
	err := r.DB.QueryRow("SELECT id, name FROM category WHERE id = $1", id).
		Scan(&c.ID, &c.Name)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) Create(c *models.Category) error {
	err := r.DB.QueryRow("INSERT INTO category (name) VALUES ($1) RETURNING id", c.Name).Scan(&c.ID)
	return err
}

func (r *CategoryRepository) Update(c *models.Category) error {
	_, err := r.DB.Exec("UPDATE category SET name=$1 WHERE id=$2", c.Name, c.ID)
	return err
}

func (r *CategoryRepository) Delete(id int) error {
	_, err := r.DB.Exec("DELETE FROM category WHERE id=$1", id)
	return err
}
