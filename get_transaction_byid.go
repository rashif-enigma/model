package model

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Mendapatkan ID transaksi dari URL parameter
func GetTransactionByID(c *gin.Context, db *sql.DB) {

	transactionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID transaksi tidak valid"})
		return
	}

	// Query untuk mendapatkan data transaksi berdasarkan ID
	query := `
		SELECT
			t.id,
			t.bill_date,
			t.entry_date,
			t.finish_date,
			t.employee_id,
			t.customer_id,
			t.address,
			t.total_bill,
			e.name AS employee_name,
			e.phone_number AS employee_phone,
			e.address AS employee_address,
			c.name AS customer_name,
			c.phone_number AS customer_phone,
			c.address AS customer_address
		FROM transactions t
		JOIN employees e ON t.employee_id = e.id
		JOIN customers c ON t.customer_id = c.id
		WHERE t.id = $1
	`

	var transaction TransactionDetails
	err = db.QueryRow(query, transactionID).Scan(
		&transaction.ID,
		&transaction.BillDate,
		&transaction.EntryDate,
		&transaction.FinishDate,
		&transaction.EmployeeID,
		&transaction.CustomerID,
		&transaction.Address,
		&transaction.TotalBill,
		&transaction.EmployeeName,
		&transaction.EmployeePhone,
		&transaction.EmployeeAddress,
		&transaction.CustomerName,
		&transaction.CustomerPhone,
		&transaction.CustomerAddress,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data transaksi"})
		return
	}

	// Query untuk mendapatkan detail tagihan
	billDetailsQuery := `
		SELECT
			bd.id AS bill_details_id,
			bd.bill_id,
			bd.product_id,
			p.name AS product_name,
			p.unit AS product_unit,
			bd.product_price,
			bd.qty
		FROM bill_details bd
		JOIN products p ON bd.product_id = p.id
		WHERE bd.bill_id = $1
	`

	rows, err := db.Query(billDetailsQuery, transactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil detail tagihan"})
		return
	}
	defer rows.Close()

	// Menyiapkan data transaksi
	var data = make(map[string]interface{})
	data["message"] = "Data Transaksi Berdasarkan ID"
	data["data"] = make([]map[string]interface{}, 0)

	// Menyiapkan detail tagihan
	var billDetails = make([]map[string]interface{}, 0)

	// Loop through bill details
	for rows.Next() {
		var billDetail BillDetail
		err := rows.Scan(
			&billDetail.ID,
			&billDetail.BillID,
			&billDetail.ProductID,
			&billDetail.ProductName,
			&billDetail.ProductUnit,
			&billDetail.ProductPrice,
			&billDetail.Qty,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca detail tagihan"})
			return
		}

		// Menambahkan detail tagihan ke slice
		billDetails = append(billDetails, map[string]interface{}{
			"id":     billDetail.ID,
			"billid": billDetail.BillID,
			"product": map[string]interface{}{
				"id":            billDetail.ProductID,
				"name":          billDetail.ProductName,
				"product_price": billDetail.ProductPrice,
				"qty":           billDetail.Qty,
				"unit":          billDetail.ProductUnit,
			},
		})
	}

	// Menambahkan data transaksi ke hasil akhir
	data["data"] = append(data["data"].([]map[string]interface{}), map[string]interface{}{
		"bill_date":   transaction.BillDate,
		"entry_date":  transaction.EntryDate,
		"finish_date": transaction.FinishDate,
		"employee": map[string]interface{}{
			"id":           transaction.EmployeeID,
			"name":         transaction.EmployeeName,
			"phone_number": transaction.EmployeePhone,
			"address":      transaction.EmployeeAddress,
		},
		"customer": map[string]interface{}{
			"id":           transaction.CustomerID,
			"name":         transaction.CustomerName,
			"phone_number": transaction.CustomerPhone,
			"address":      transaction.CustomerAddress,
		},
		"bill_details": billDetails,
		"total_bill":   transaction.TotalBill,
	})

	// Mengembalikan response
	c.JSON(http.StatusOK, data)
}
