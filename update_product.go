package model

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UpdateProduct adalah handler untuk memperbarui data produk
func UpdateProduct(c *gin.Context, db *sql.DB) {
	// Membuat variabel untuk menyimpan data produk yang diterima dari body request
	var product Products

	// Membaca data produk dari body request
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kesalahan Inputan"})
		return
	}

	// Pastikan ID produk yang ingin diperbarui sudah ada
	if product.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID produk harus disertakan"})
		return
	}

	// Update data produk di dalam database
	query := "UPDATE products SET name = $1, price = $2, unit = $3 WHERE id = $4"
	_, err := db.Exec(query, product.Name, product.Price, product.Unit, product.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data produk"})
		return
	}

	// Mengembalikan respons JSON yang menyatakan bahwa data produk telah berhasil diperbarui
	c.JSON(http.StatusOK, gin.H{"message": "Data produk berhasil diperbarui", "data": product})
}
