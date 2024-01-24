package model

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// handler menampilkan data pelanggan bedasarkan id pelanggan
func GetCustomerByID(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	query := "SELECT id, name, phone_number, address FROM customers WHERE id = $1"
	row := db.QueryRow(query, id)

	var customer Customers
	err := row.Scan(&customer.ID, &customer.Name, &customer.PhoneNumber, &customer.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan data pelanggan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menampilkan Data Pelanggan Bedasarkan Id Pelanggan", "data": customer})
}
