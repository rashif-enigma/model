package model

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListCustomers adalah handler untuk menampilkan semua data pelanggan
func ListCustomers(c *gin.Context, db *sql.DB) {
	query := "SELECT id, name, phone_number, address FROM customers"
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pelanggan"})
		return
	}
	defer rows.Close()

	var customers []Customers
	for rows.Next() {
		var customer Customers
		err := rows.Scan(&customer.ID, &customer.Name, &customer.PhoneNumber, &customer.Address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data pelanggan"})
			return
		}
		customers = append(customers, customer)
	}

	c.JSON(http.StatusOK, gin.H{"message": "List Data Pelanggan ", "data": customers})
}
