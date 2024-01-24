package model

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// getEmployeeByID mendapatkan data karyawan berdasarkan ID dari database
func GetEmployeeByID(c *gin.Context, db *sql.DB) {
	// Mendapatkan ID karyawan dari parameter URL
	employeeID := c.Param("id")

	// Query SQL untuk mendapatkan data karyawan berdasarkan ID
	query := `
		SELECT id, name, phone_number, address
		FROM employees
		WHERE id = $1
	`

	// Menjalankan query SQL
	row := db.QueryRow(query, employeeID)

	// Membaca data karyawan dari hasil query
	var employee Employee
	err := row.Scan(
		&employee.ID,
		&employee.Name,
		&employee.PhoneNumber,
		&employee.Address,
	)

	if err != nil {
		// Jika karyawan tidak ditemukan, mengirimkan response not found
		c.JSON(http.StatusNotFound, gin.H{"error": "Data Karyawan tidak ditemukan"})
		return
	}

	// Mengirimkan response dengan data karyawan
	c.JSON(http.StatusOK, gin.H{"message": "Data Karyawan Bedasarkan Id Karyawan", "data": employee})
}
