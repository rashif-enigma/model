package model

import (
	"api-laundry-enigma/config"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// Handler untuk koneksi ke database config
var db *sql.DB

func init() {
	var err error
	db, err = config.ConnectDB()
	if err != nil {
		panic(err)
	}
}

// TransactionDetails adalah struktur data untuk transaksi dengan detail
type TransactionDetails struct {
	ID              int          `json:"id"`
	BillDate        string       `json:"bill_date"`
	EntryDate       string       `json:"entry_date"`
	FinishDate      string       `json:"finish_date"`
	EmployeeID      int          `json:"employee_id"`
	CustomerID      int          `json:"customer_id"`
	Address         string       `json:"address"`
	TotalBill       float64      `json:"total_bill"`
	EmployeeName    string       `json:"employee_name"`
	EmployeePhone   string       `json:"employee_phone"`
	EmployeeAddress string       `json:"employee_address"`
	CustomerName    string       `json:"customer_name"`
	CustomerPhone   string       `json:"customer_phone"`
	CustomerAddress string       `json:"customer_address"`
	BillDetails     []BillDetail `json:"bill_details"`
}

// BillDetail adalah struktur data untuk detail tagihan
type BillDetail struct {
	ID           int     `json:"id"`
	BillID       int     `json:"bill_id"`
	ProductID    int     `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductUnit  string  `json:"product_unit"`
	ProductPrice float64 `json:"product_price"`
	Qty          int     `json:"qty"`
}

// handler untuk membuat data transaksi berdasarkan id pelanggan
// fungsi ini memuat 2 fungsi: otomatis memasukkan data ke dalam tabel bill_details dan Transactions
func CreateTransaction(c *gin.Context, db *sql.DB) {
	var request struct {
		Data struct {
			BillDate    string `json:"bill_date"`
			EntryDate   string `json:"entry_date"`
			FinishDate  string `json:"finish_date"`
			EmployeeID  int    `json:"employee_id"`
			CustomerID  int    `json:"customer_id"`
			Address     string `json:"address"`
			BillDetails []struct {
				BillID       int `json:"bill_id"`
				ProductID    int `json:"product_id"`
				ProductPrice int `json:"product_price"`
				Qty          int `json:"qty"`
			} `json:"bill_details"`
			TotalBill int `json:"total_bill"`
		} `json:"data"`
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Println("Error parsing JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kesalahan Inputan"})
		return
	}

	// Menambahkan data transaksi ke dalam tabel transactions
	query := "INSERT INTO transactions (bill_date, entry_date, finish_date, employee_id, customer_id, address, total_bill) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	var insertedID int
	err := db.QueryRow(query, request.Data.BillDate, request.Data.EntryDate, request.Data.FinishDate, request.Data.EmployeeID, request.Data.CustomerID, request.Data.Address, request.Data.TotalBill).Scan(&insertedID)
	if err != nil {
		fmt.Println("Error inserting transaction:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal Menambahkan Data Transaksi"})
		return
	}

	// Menambahkan detail tagihan ke dalam tabel bill_details
	for _, billDetail := range request.Data.BillDetails {
		query := "INSERT INTO bill_details (bill_id, product_id, product_price, qty) VALUES ($1, $2, $3, $4) RETURNING id"
		var billDetailID int
		err := db.QueryRow(query, billDetail.BillID, billDetail.ProductID, billDetail.ProductPrice, billDetail.Qty).Scan(&billDetailID)
		if err != nil {
			fmt.Println("Error inserting bill detail:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal Menambahkan Detail Transaksi"})
			return
		}
	}

	// Menetapkan ID transaksi yang baru dibuat ke dalam objek transaction
	response := struct {
		Data    gin.H  `json:"data"`
		Message string `json:"message"`
	}{
		gin.H{
			"bill_date":    request.Data.BillDate,
			"entry_date":   request.Data.EntryDate,
			"finish_date":  request.Data.FinishDate,
			"employee_id":  request.Data.EmployeeID,
			"customer_id":  request.Data.CustomerID,
			"address":      request.Data.Address,
			"bill_details": request.Data.BillDetails,
			"total_bill":   request.Data.TotalBill,
		},
		"Data Transaksi Berhasil di tambahkan",
	}

	c.JSON(http.StatusCreated, response)
}
