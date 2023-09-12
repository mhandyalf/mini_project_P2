package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"mini_project_p2/models"
	"net/http"
	"strings"
)

const (
	ipaymuVa  = "1179000899"                     // VA iPaymu Anda
	ipaymuKey = "QbGcoO0Qds9sQFDmY0MWg1Tq.xtuh1" // Kunci API iPaymu Anda
)

// CreateIPaymuPayment membuat permintaan pembayaran ke iPaymu
func CreateIPaymuPayment(req *models.PaymentRequest) (string, error) {
	// Mengonversi payload menjadi JSON
	postBody, _ := json.Marshal(req)

	// Menghitung hash dari body permintaan
	bodyHash := sha256.Sum256(postBody)
	bodyHashToString := hex.EncodeToString(bodyHash[:])

	// Membuat string untuk di-signature
	stringToSign := "POST:" + ipaymuVa + ":" + strings.ToLower(string(bodyHashToString)) + ":" + ipaymuKey

	// Membuat signature menggunakan HMAC-SHA256
	h := hmac.New(sha256.New, []byte(ipaymuKey))
	h.Write([]byte(stringToSign))
	signature := hex.EncodeToString(h.Sum(nil))

	// Membuat body permintaan sebagai NopCloser
	reqBody := ioutil.NopCloser(strings.NewReader(string(postBody)))

	// Membuat permintaan HTTP POST ke iPaymu
	resp, err := http.Post("https://sandbox.ipaymu.com/api/v2/payment", "application/json", reqBody)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Membaca respons dari server iPaymu
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Mengembalikan respons sebagai string
	return string(body), nil
}
