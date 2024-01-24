package model

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// handler untuk mengupdate data pelanggan
func UpdateCustomer(c *gin.Context, db *sql.DB) {
	var customer Customers
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	query := "UPDATE customers SET name = $1, phone_number = $2, address = $3 WHERE id = $4"
	_, err := db.Exec(query, customer.Name, customer.PhoneNumber, customer.Address, customer.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data pelanggan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data pelanggan berhasil diperbarui", "data": customer})
}
