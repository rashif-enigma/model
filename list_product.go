package model

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListProducts adalah handler untuk menampilkan semua jasa laundry / produk
func ListProducts(c *gin.Context, db *sql.DB) {
	query := "SELECT id, name, price, unit FROM products"
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Gagal membaca data produk: %v", err)})
		return
	}
	defer rows.Close()

	var products []Products
	for rows.Next() {
		var product Products
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Unit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Gagal membaca data produk: %v", err)})
			return
		}
		products = append(products, product)
	}

	c.JSON(http.StatusOK, gin.H{"message": "List Semua Produk Jasa", "data": products})
}
