package model

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DeleteProduct adalah handler untuk menghapus data produk
func DeleteProduct(c *gin.Context, db *sql.DB) {
	// Mendapatkan ID produk dari path parameter
	productID := c.Param("id")

	// Menjalankan query untuk menghapus data produk dari database
	query := "DELETE FROM products WHERE id = $1"
	_, err := db.Exec(query, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data produk"})
		return
	}

	// Mengembalikan respons JSON yang menyatakan bahwa data produk telah berhasil dihapus
	c.JSON(http.StatusOK, gin.H{"message": "Data produk berhasil dihapus"})
}
