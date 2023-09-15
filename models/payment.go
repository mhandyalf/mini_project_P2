package models

// PaymentData struct untuk menyimpan data pembayaran
type PaymentData struct {
	Product     []string  `json:"product"`
	Qty         []int8    `json:"qty"`
	Price       []float64 `json:"price"`
	ReturnURL   string    `json:"returnUrl"`
	CancelURL   string    `json:"cancelUrl"`
	NotifyURL   string    `json:"notifyUrl"`
	ReferenceID string    `json:"referenceId"`
	BuyerName   string    `json:"buyerName"`
	BuyerEmail  string    `json:"buyerEmail"`
	BuyerPhone  string    `json:"buyerPhone"`
}
