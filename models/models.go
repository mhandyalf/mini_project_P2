package models

import "time"

// Struct untuk tabel "users"
type User struct {
	ID            int     `json:"id"`
	Email         string  `json:"email"`
	Username      string  `json:"username"`
	Password      string  `json:"password"`
	DepositAmount float64 `json:"deposit_amount"`
}

// Struct untuk tabel "book_inventory"
type BookInventory struct {
	BookID            int     `json:"book_id"`
	Name              string  `json:"name"`
	StockAvailability int     `json:"stock_availability"`
	RentalCosts       float64 `json:"rental_costs"`
	Category          string  `json:"category"`
}

// Struct untuk tabel "rental_history"
type RentalHistory struct {
	RentalID   int       `json:"rental_id"`
	UserID     int       `json:"user_id"`
	BookID     int       `json:"book_id"`
	RentalDate time.Time `json:"rental_date"`
	ReturnDate time.Time `json:"return_date,omitempty"`
	RentalCost float64   `json:"rental_cost"`
}
