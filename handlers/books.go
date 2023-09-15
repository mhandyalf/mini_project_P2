package handlers

import (
	"mini_project_p2/models"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (a *Auth) GetAllBooks(c echo.Context) error {
	// Buat variabel untuk menyimpan hasil query
	var books []models.BookInventory

	// SQL query untuk mengambil semua buku
	query := "SELECT * FROM book_inventory"

	// Eksekusi query ke database
	if err := a.DB.Raw(query).Scan(&books).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, books)
}

func (a *Auth) UpdateBook(c echo.Context) error {
	var book models.BookInventory
	if err := c.Bind(&book); err != nil {
		return err
	}

	query := `
        UPDATE book_inventory
        SET
            name = ?,
            stock_availability = ?,
            rental_costs = ?,
            category = ?
        WHERE
            book_id = ?
    `
	result := a.DB.Exec(query, book.Name, book.StockAvailability, book.RentalCosts, book.Category, book.BookID)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Buku tidak ditemukan")
	}

	return c.JSON(http.StatusOK, book)
}

func (a *Auth) DeleteBook(c echo.Context) error {
	var book models.BookInventory
	if err := c.Bind(&book); err != nil {
		return err
	}

	query := `

		DELETE FROM book_inventory

		WHERE

			book_id = ?

	`

	result := a.DB.Exec(query, book.BookID)

	if result.Error != nil {

		return result.Error

	}

	if result.RowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Buku tidak ditemukan")

	}

	return c.JSON(http.StatusOK, book)
}

func (a *Auth) ReturnBook(c echo.Context) error {
	var returnRequest struct {
		RentalID int `json:"rental_id"`
	}

	if err := c.Bind(&returnRequest); err != nil {
		return err
	}

	userID := c.Get("user").(models.User).ID

	// Periksa apakah rental dengan ID yang diberikan ada dalam rental history pengguna
	var rental models.RentalHistory
	query := "SELECT * FROM rental_history WHERE rental_id = $1 AND user_id = $2"
	if err := a.DB.Raw(query, returnRequest.RentalID, userID).Scan(&rental).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Rental tidak ditemukan")
	}

	// Set tanggal pengembalian
	rental.ReturnDate = time.Now()

	// Perbarui tabel rental history dengan tanggal pengembalian
	query = "UPDATE rental_history SET return_date = $1 WHERE rental_id = $2"
	if err := a.DB.Exec(query, rental.ReturnDate, returnRequest.RentalID).Error; err != nil {
		return err
	}

	// Kembalikan buku ke stok
	var book models.BookInventory
	query = "SELECT * FROM book_inventory WHERE book_id = $1"
	if err := a.DB.Raw(query, rental.BookID).Scan(&book).Error; err != nil {
		return err
	}

	book.StockAvailability++
	query = "UPDATE book_inventory SET stock_availability = $1 WHERE book_id = $2"
	if err := a.DB.Exec(query, book.StockAvailability, rental.BookID).Error; err != nil {
		return err
	}

	// Mengembalikan biaya sewa ke deposit pengguna
	user := c.Get("user").(models.User)
	user.DepositAmount += rental.RentalCost

	// Perbarui nilai DepositAmount di basis data untuk pengguna yang bersangkutan
	query = "UPDATE users SET deposit_amount = $1 WHERE id = $2"
	if err := a.DB.Exec(query, user.DepositAmount, userID).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "Buku telah berhasil dikembalikan")
}
