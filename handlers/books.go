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

func (a *Auth) RentBook(c echo.Context) error {
	var rent models.RentalHistory

	query := "INSERT INTO rental_history (book_id, user_id, rental_date, rental_cost) VALUES ($1, $2, $3, $4) RETURNING *"

	if err := c.Bind(&rent); err != nil {
		return err
	}

	rent.RentalDate = time.Now()
	rent.RentalCost = 111

	if err := a.DB.Raw(query, rent.BookID, rent.UserID, rent.RentalDate, rent.RentalCost).Scan(&rent).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, rent)
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
