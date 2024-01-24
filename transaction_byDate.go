package model

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func GetTransactionsByDate(c *gin.Context, db *sql.DB) {
	// Ambil nilai parameter date dari query string
	dateParam := c.Query("date")

	// Validasi format tanggal
	parsedDate, err := time.Parse("02/01/2006", dateParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal tidak valid. Gunakan format DD/MM/YYYY"})
		return
	}

	// Query ke database untuk mendapatkan transaksi pada tanggal tertentu
	transactions, err := getTransactionsFromDatabase(parsedDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data transaksi dari database"})
		return
	}

	// Berikan response dengan data transaksi
	c.JSON(http.StatusOK, gin.H{
		"message": "Data Transaksi Berdasarkan Tanggal",
		"date":    parsedDate.Format("02 January 2006"),
		"data":    transactions,
	})
}

func getTransactionsFromDatabase(date time.Time) ([]gin.H, error) {
	// Gantilah dengan query sesuai dengan struktur database Anda
	query := `
		SELECT 
			bill_date,
			total_bill,
			customer_id,
			employee_id
		FROM transactions
		WHERE bill_date = $1
	`

	rows, err := db.Query(query, date.Format("02/01/2006"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []gin.H
	for rows.Next() {
		var transaction = make(map[string]interface{})
		var billDate string
		var totalBill float64
		var customerID, employeeID int
		err := rows.Scan(
			&billDate,
			&totalBill,
			&customerID,
			&employeeID,
		)
		if err != nil {
			return nil, err
		}

		transaction["bill_date"] = billDate
		transaction["total_bill"] = totalBill
		transaction["customer_id"] = customerID
		transaction["employee_id"] = employeeID

		// Tambahan logika atau query lainnya sesuai kebutuhan

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
