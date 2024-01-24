package model

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Products struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Unit  string  `json:"unit"`
}

// Handler untuk membuat produk / Jasa Laundry
func CreateProduct(c *gin.Context, db *sql.DB) {
	var product Products
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	query := "INSERT INTO products (name, price, unit) VALUES ($1, $2, $3) RETURNING id"
	var insertedID int
	err := db.QueryRow(query, product.Name, product.Price, product.Unit).Scan(&insertedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan produk"})
		return
	}

	product.ID = insertedID
	c.JSON(http.StatusCreated, gin.H{"message": "Data produk berhasil ditambahkan", "data": product})
}
