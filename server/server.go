package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

const urlQuotation = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

type Quotation struct {
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func main() {
	http.HandleFunc("/cotacao", GetQuotation)
	http.ListenAndServe(":8000", nil)
}

func GetQuotation(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	res, err := Request(ctx)
	if err != nil {
		println("Erro ao buscar a cotação")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		println("Erro ao encodar a resposta")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Request(ctx context.Context) (*Quotation, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", urlQuotation, nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		println("Erro ao ler o body")
		return nil, err
	}
	var q Quotation
	err = json.Unmarshal(body, &q)
	if err != nil {
		println("Erro ao converter em json")
		return nil, err
	}
	return &q, nil
}
