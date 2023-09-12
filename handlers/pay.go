package handlers

import (
	"mini_project_p2/middleware"
	"mini_project_p2/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Auth) CreatePayment(c echo.Context) error {
	// Membaca data permintaan dari client
	req := new(models.PaymentRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, "Format data tidak valid")
	}

	// Membuat permintaan pembayaran ke iPaymu
	response, err := middleware.CreateIPaymuPayment(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Mengembalikan respons ke client
	return c.JSON(http.StatusOK, response)
}
