package repository

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type T struct {
	Success bool `json:"success"`
	Query   struct {
		From   string `json:"from"`
		To     string `json:"to"`
		Amount int    `json:"amount"`
	} `json:"query"`
	Info struct {
		Timestamp int     `json:"timestamp"`
		Rate      float64 `json:"rate"`
	} `json:"info"`
	Date   string  `json:"date"`
	Result float64 `json:"result"`
}

func (r *Repository) BalanceInCurrency(id int, c string) (float64, error) {
	accessKey := os.Getenv("ACCESS_KEY")

	amount, err := r.UserBalance(id)
	if err != nil {
		return 0, err
	}

	client := http.Client{}

	url := fmt.Sprintf("https://api.apilayer.com/exchangerates_data/convert?to=%v&from=RUB&amount=%v", c, amount.Balance)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Add("apikey", accessKey)

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	var body T

	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return 0, err
	}

	return body.Result, nil
}
