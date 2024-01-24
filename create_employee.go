package model

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Employee struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
}

// createEmployee menambahkan karyawan baru ke dalam database
func CreateEmployeeHandler(c *gin.Context, db *sql.DB) {
	var employee Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kesalahan Inputan"})
		return
	}

	// Memasukkan data karyawan ke dalam tabel employees
	query := "INSERT INTO employees (name, phone_number, address) VALUES ($1, $2, $3) RETURNING id"
	var insertedID int
	err := db.QueryRow(query, employee.Name, employee.PhoneNumber, employee.Address).Scan(&insertedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal Menambahkan Data Karyawan"})
		return
	}

	employee.ID = insertedID
	c.JSON(http.StatusCreated, gin.H{"message": "Data Karyawan Berhasil Ditambahkan", "data": employee})
}
