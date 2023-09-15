package handlers

import (
	"mini_project_p2/middleware"
	"mini_project_p2/models"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

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

	// Ambil nilai DepositAmount dari pengguna
	user := c.Get("user").(models.User)
	depositAmount := user.DepositAmount

	// Kurangkan DepositAmount sesuai dengan biaya sewa buku
	depositAmount -= rent.RentalCost

	// Perbarui nilai DepositAmount di basis data untuk pengguna yang bersangkutan
	query = `
		UPDATE users SET deposit_amount = $1 WHERE id = $2
	`
	if err := a.DB.Exec(query, depositAmount, userID).Error; err != nil {
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

func (a *Auth) DeleteRentalHistory(c echo.Context) error {
	// Dapatkan ID pengguna dari token JWT
	userID := c.Get("user").(models.User).ID

	var deleteRequest struct {
		RentalID int `json:"rental_id"`
	}

	if err := c.Bind(&deleteRequest); err != nil {
		return err
	}

	// Periksa apakah rental dengan ID yang diberikan ada dalam rental history pengguna
	var rental models.RentalHistory
	query := "SELECT * FROM rental_history WHERE rental_id = $1 AND user_id = $2"
	if err := a.DB.Raw(query, deleteRequest.RentalID, userID).Scan(&rental).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Rental tidak ditemukan")
	}

	// Hapus rental history dari basis data
	query = "DELETE FROM rental_history WHERE rental_id = $1 AND user_id = $2"
	if err := a.DB.Exec(query, deleteRequest.RentalID, userID).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "Rental history telah berhasil dihapus")
}
