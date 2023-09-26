package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const urlQuotation = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
const initScriptDatabase = `
CREATE TABLE IF NOT EXISTS quotation (
	id varchar(255) NOT NULL PRIMARY KEY,
	code varchar(255),
	codein varchar(255),
	name varchar(255),
	high varchar(255),
	low varchar(255),
	varBid varchar(255),
	pctChange varchar(255),
	bid varchar(255),
	ask varchar(255),
	timestamp varchar(255),
	create_date varchar(255));
`

var db *sql.DB

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
	db = InitializeDatabase()
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

func InitializeDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "quotation.db")
	if err != nil {
		println("Erro ao conectar no banco de dados")
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec(initScriptDatabase)
	if err != nil {
		println("Erro ao executar script inicial do banco de dados")
		panic(err)
	}
	return db
}
