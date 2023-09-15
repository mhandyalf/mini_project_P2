package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mini_project_p2/models"
	"net/http"
	"os"
	"strings"
)

// Define your iPaymu credentials
var ipaymuVa = os.Getenv("va")      // Your iPaymu VA
var ipaymuKey = os.Getenv("ipaymu") // Your iPaymu API key

func SendPaymentRequest(paymentData models.PaymentData) error {
	postBody, err := json.Marshal(paymentData)
	if err != nil {
		return err
	}

	bodyHash := sha256.Sum256(postBody)
	bodyHashToString := hex.EncodeToString(bodyHash[:])
	stringToSign := "POST:" + ipaymuVa + ":" + strings.ToLower(string(bodyHashToString)) + ":" + ipaymuKey

	h := hmac.New(sha256.New, []byte(ipaymuKey))
	h.Write([]byte(stringToSign))
	signature := hex.EncodeToString(h.Sum(nil))

	urlStr := "https://sandbox.ipaymu.com/api/v2/payment"
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(string(postBody)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("va", ipaymuVa)
	req.Header.Set("signature", signature)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("iPaymu Response: %s\n", string(responseBody))

	return nil
}
