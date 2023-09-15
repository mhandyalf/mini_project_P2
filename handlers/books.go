package handlers

import (
	"mini_project_p2/middleware"
	"mini_project_p2/models"
	"net/http"
	"strconv"
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

	userID := c.Get("user").(models.User).ID

	query := `
		INSERT INTO rental_history (book_id, user_id, rental_date, rental_cost)
		VALUES ($1, $2, $3, $4)
		RETURNING *
	`

	if err := c.Bind(&rent); err != nil {
		return err
	}

	rent.RentalDate = time.Now()

	var book models.BookInventory

	query = `
		SELECT * FROM book_inventory WHERE book_id = $1
	`
	if err := a.DB.Raw(query, rent.BookID).Scan(&book).Error; err != nil {
		return err
	}

	rent.RentalCost = book.RentalCosts

	if book.StockAvailability > 0 {
		book.StockAvailability--

		query = `
			UPDATE book_inventory SET stock_availability = $1 WHERE book_id = $2
		`
		if err := a.DB.Exec(query, book.StockAvailability, rent.BookID).Error; err != nil {
			return err
		}
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, "Stok buku habis")
	}

	query = `
		INSERT INTO rental_history (book_id, user_id, rental_date, rental_cost)
		VALUES ($1, $2, $3, $4)
		RETURNING *
	`
	if err := a.DB.Raw(query, rent.BookID, userID, rent.RentalDate, rent.RentalCost).Scan(&rent).Error; err != nil {
		return err
	}

	paymentData := models.PaymentData{
		Product:     []string{"Book Rental"},
		Qty:         []int8{1},
		Price:       []float64{rent.RentalCost},
		ReturnURL:   "http://your-website/thank-you-page",
		CancelURL:   "http://your-website/cancel-page",
		NotifyURL:   "http://your-website/callback-url",
		ReferenceID: "RENTAL-" + strconv.Itoa(rent.RentalID),
		BuyerName:   "test",
		BuyerEmail:  "test@mail.com",
		BuyerPhone:  "08123456789",
	}

	err := middleware.SendPaymentRequest(paymentData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
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
