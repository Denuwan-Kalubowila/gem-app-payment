package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
)

func main() {

	stripe.Key = "sk_test_51ODedtJcKvObEPTHSbME7joQwRYORG37jmhGILSy9PfaBY53U0gruEV6wQqqJVjMISzM41O5IyJlksM6k8Jc6qiH00BaakBdPV"

	http.HandleFunc("/create-payment-intent", handlePaymentIntent)
	http.HandleFunc("/health", handleHealth)

	log.Println("Server is running on port :4242 ...")

	var err error = http.ListenAndServe("localhost:4242", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handlePaymentIntent(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ProductData map[string]int64 `json:"product_data"`
		FirstName   string           `json:"first_name"`
		LastName    string           `json:"last_name"`
		Address1    string           `json:"address_1"`
		Address2    string           `json:"address_2"`
		City        string           `json:"city"`
		State       string           `json:"state"`
		Zip         string           `json:"zip"`
		Country     string           `json:"country"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(claculateOrderAmount(req.ProductData)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	paymentIntent, err := paymentintent.New(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Println(paymentIntent.ClientSecret)
	var response struct {
		ClientSecret string `json:"clientSecret"`
	}
	response.ClientSecret = paymentIntent.ClientSecret

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = io.Copy(w, &buf)
	if err != nil {
		fmt.Println(err)
	}

}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	response := []byte("Server is up and running!")
	_, err := w.Write(response)
	if err != nil {
		fmt.Println(err)
		return

	}

}

func claculateOrderAmount(productID map[string]int64) int64 {
	var price int64
	for product, count := range productID {
		switch product {
		case "Blue Sapphire":
			price += 100000 * count
		case "Ruby":
			price += 200000 * count
		case "Yellow Sapphire":
			price += 300000 * count
		default:
			fmt.Printf("%s Not a Item", product)
		}
	}
	return price
}
