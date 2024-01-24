package model

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListEmployees(c *gin.Context, db *sql.DB) {
	// Mengambil data karyawan dari tabel employees
	rows, err := db.Query("SELECT id, name, phone_number, address FROM employees")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal Mengambil Data Karyawan"})
		return
	}
	defer rows.Close()

	var employees []Employee

	// Mengiterasi hasil query dan menambahkan data karyawan ke dalam slice employees
	for rows.Next() {
		var employee Employee
		err := rows.Scan(&employee.ID, &employee.Name, &employee.PhoneNumber, &employee.Address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal Membaca Data Karyawan"})
			return
		}
		employees = append(employees, employee)
	}

	// Memeriksa apakah terjadi error selama iterasi
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal Membaca Data Karyawan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Semua Data Karyawan ", "data": employees})
}
