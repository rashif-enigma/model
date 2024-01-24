package model

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Customers adalah struktur data untuk tabel pelanggan
type Customers struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
}

// CreateCustomer adalah handler untuk membuat pelanggan baru
func CreateCustomer(c *gin.Context, db *sql.DB) {
	var customer Customers
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	query := "INSERT INTO customers (name, phone_number, address) VALUES ($1, $2, $3) RETURNING id"
	var insertedID int
	err := db.QueryRow(query, customer.Name, customer.PhoneNumber, customer.Address).Scan(&insertedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal Menambahkan Pelanggan"})
		return
	}

	customer.ID = insertedID
	c.JSON(http.StatusCreated, gin.H{"message": "Data Pelanggan Berhasil Ditambahkan", "data": customer})
}
