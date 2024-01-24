package model

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler untuk menghapus data pelanggan berdasarkan ID
func DeleteCustomer(c *gin.Context, db *sql.DB) {
	customerID := c.Param("id")

	query := "DELETE FROM customers WHERE id = $1"
	_, err := db.Exec(query, customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data pelanggan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data pelanggan berhasil dihapus"})
}
