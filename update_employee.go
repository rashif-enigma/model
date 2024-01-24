package model

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// updateEmployee memperbarui data karyawan di dalam database
func UpdateEmployee(c *gin.Context, db *sql.DB) {
	// Mendapatkan ID karyawan dari parameter URL
	employeeID := c.Param("id")

	// Membaca data karyawan yang akan diperbarui dari body permintaan
	var updatedEmployee Employee
	if err := c.ShouldBindJSON(&updatedEmployee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data karyawan tidak valid"})
		return
	}

	// Query SQL untuk memperbarui data karyawan berdasarkan ID
	query := `
		UPDATE employees
		SET name = $1, phone_number = $2, address = $3
		WHERE id = $4
	`

	// Menjalankan query SQL untuk memperbarui data karyawan
	_, err := db.Exec(query, updatedEmployee.Name, updatedEmployee.PhoneNumber, updatedEmployee.Address, employeeID)
	if err != nil {
		// Jika terjadi kesalahan, mengirimkan response internal server error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data karyawan"})
		return
	}

	// Mengirimkan response sukses
	c.JSON(http.StatusOK, gin.H{"message": "Data karyawan berhasil diperbarui"})
}
