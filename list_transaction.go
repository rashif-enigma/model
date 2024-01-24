package model

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// handler menampilkan semua data transaksi
func ListTransactions(c *gin.Context, db *sql.DB) {
	query := `
        SELECT 
            t.id, 
            t.bill_date, 
            t.entry_date, 
            t.finish_date, 
            t.employee_id, 
            t.customer_id, 
            t.total_bill,
            c.address AS customer_address,
            c.name AS customer_name,
            c.phone_number AS customer_phone,
            e.name AS employee_name,
            e.phone_number AS employee_phone,
            e.address AS employee_address
        FROM transactions t
        JOIN customers c ON t.customer_id = c.id
        JOIN employees e ON t.employee_id = e.id
    `
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data transaksi"})
		return
	}
	defer rows.Close()

	transactionsMap := make(map[int]TransactionDetails)

	for rows.Next() {
		var transaction TransactionDetails
		err := rows.Scan(
			&transaction.ID,
			&transaction.BillDate,
			&transaction.EntryDate,
			&transaction.FinishDate,
			&transaction.EmployeeID,
			&transaction.CustomerID,
			&transaction.TotalBill,
			&transaction.CustomerAddress,
			&transaction.CustomerName,
			&transaction.CustomerPhone,
			&transaction.EmployeeName,
			&transaction.EmployeePhone,
			&transaction.EmployeeAddress,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data transaksi"})
			return
		}

		transactionsMap[transaction.ID] = transaction
	}

	// Mendapatkan detail tagihan (bill details) dan menambahkannya ke map transactionsMap
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
    `

	billDetailsRows, err := db.Query(billDetailsQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil detail tagihan"})
		return
	}
	defer billDetailsRows.Close()

	for billDetailsRows.Next() {
		var transactionID, productID, billDetailsID, qty int
		var productName, productUnit string
		var productPrice float64

		err := billDetailsRows.Scan(
			&billDetailsID,
			&transactionID,
			&productID,
			&productName,
			&productUnit,
			&productPrice,
			&qty,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca detail tagihan"})
			return
		}

		transaction, ok := transactionsMap[transactionID]
		if !ok {
			// Skip jika transaksi tidak ditemukan (seharusnya tidak terjadi)
			continue
		}

		billDetails := BillDetail{
			ID:           billDetailsID,
			BillID:       transaction.ID,
			ProductID:    productID,
			ProductName:  productName,
			ProductUnit:  productUnit,
			ProductPrice: productPrice,
			Qty:          qty,
		}

		transaction.BillDetails = append(transaction.BillDetails, billDetails)
		transactionsMap[transactionID] = transaction
	}

	// Membuat response sesuai format
	var response []gin.H
	for _, transaction := range transactionsMap {
		var billDetailsResponse []gin.H
		for _, detail := range transaction.BillDetails {
			billDetailsResponse = append(billDetailsResponse, gin.H{
				"id":            detail.ID,
				"product":       gin.H{"id": detail.ProductID, "name": detail.ProductName, "price": detail.ProductPrice, "unit": detail.ProductUnit},
				"product_price": detail.ProductPrice,
				"qty":           detail.Qty,
			})
		}

		response = append(response, gin.H{
			"bill_date":    transaction.BillDate,
			"bill_details": billDetailsResponse,
			"customer_id": gin.H{
				"id":           transaction.CustomerID,
				"address":      transaction.CustomerAddress,
				"name":         transaction.CustomerName,
				"phone_number": transaction.CustomerPhone,
			},
			"employee_id": gin.H{
				"id":           transaction.EmployeeID,
				"address":      transaction.EmployeeAddress,
				"name":         transaction.EmployeeName,
				"phone_number": transaction.EmployeePhone,
			},
			"entry_date":  transaction.EntryDate,
			"finish_date": transaction.FinishDate,
			"bill_id":     transaction.ID,
			"total_bill":  transaction.TotalBill,
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "List Semua Data Transaksi", "data": response})
}
