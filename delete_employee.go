package model

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// deleteEmployee menghapus data karyawan dari database
func DeleteEmployee(c *gin.Context, db *sql.DB) {
	// Mendapatkan ID karyawan dari parameter URL
	employeeID := c.Param("id")

	// Query SQL untuk menghapus data karyawan berdasarkan ID
	query := "DELETE FROM employees WHERE id = $1"

	// Menjalankan query SQL untuk menghapus data karyawan
	_, err := db.Exec(query, employeeID)
	if err != nil {
		// Jika terjadi kesalahan, mengirimkan response internal server error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data karyawan"})
		return
	}

	// Mengirimkan response sukses
	c.JSON(http.StatusOK, gin.H{"message": "Data karyawan berhasil dihapus"})
}
