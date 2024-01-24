package model

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// handler untuk menampilkan jasa produk laundry bedasarkan id
func GetProductByID(c *gin.Context, db *sql.DB) {
	// Mendapatkan ID produk dari parameter URL
	id := c.Param("id")

	// Query untuk mendapatkan data produk berdasarkan ID
	query := "SELECT id, name, price, unit FROM products WHERE id = $1"
	row := db.QueryRow(query, id)

	// Membaca data produk
	var product Products
	err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Unit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan data produk"})
		return
	}

	// Mengembalikan data produk dalam respons JSON
	c.JSON(http.StatusOK, gin.H{"message": "Menampilkan Data Produk Berdasarkan Id Produk ", "data": product})
}
