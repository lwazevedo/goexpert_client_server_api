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
	req, err := http.NewRequestWithContext(ctx, "GET", urlQuotation, nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		println("Erro ao buscar a cotação")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		println("Erro ao ler o corpo da requisição")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var q Quotation
	err = json.Unmarshal(body, &q)
	if err != nil {
		println("Erro ao serializar o corpo da requisição em uma struct")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(q)
	if err != nil {
		println("Erro ao encodar a resposta")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
