package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
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
		ID         string `json:"id"`
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
type Output struct {
	Bid float64 `json:"bid"`
}

func main() {
	db = initializeDatabase()
	defer db.Close()
	http.HandleFunc("/cotacao", GetQuotation)
	http.ListenAndServe(":8080", nil)
}

func GetQuotation(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	res, err := quotationRequest(ctx)
	if err != nil {
		println("Erro ao buscar a cotação")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	q, err := saveQuotation(r, res)
	if err != nil {
		println("Erro ao tentar salvar a cotação")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	out, err := mapperOutput(q)
	if err != nil {
		println("Erro ao tentar mapear o bid para output")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(out)
	if err != nil {
		println("Erro ao encodar a resposta")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func quotationRequest(ctx context.Context) (*Quotation, error) {
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

func initializeDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "quotation.db")
	if err != nil {
		println("Erro ao conectar no banco de dados")
		panic(err)
	}
	_, err = db.Exec(initScriptDatabase)
	if err != nil {
		println("Erro ao executar script inicial do banco de dados")
		panic(err)
	}
	return db
}

func saveQuotation(r *http.Request, q *Quotation) (*Quotation, error) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Millisecond)
	defer cancel()
	stmt, err := db.Prepare("insert into quotation(id, code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		println("Erro ao tentar sanitizar os dados via stmt")
		return nil, err
	}
	q.Usdbrl.ID = uuid.New().String()
	_, err = stmt.ExecContext(ctx,
		q.Usdbrl.ID,
		q.Usdbrl.Code,
		q.Usdbrl.Codein,
		q.Usdbrl.Name,
		q.Usdbrl.High,
		q.Usdbrl.Low,
		q.Usdbrl.VarBid,
		q.Usdbrl.PctChange,
		q.Usdbrl.Bid,
		q.Usdbrl.Ask,
		q.Usdbrl.Timestamp,
		q.Usdbrl.CreateDate)
	if err != nil {
		println("Erro ao tentar executar insert do quotation")
		return nil, err
	}
	return q, nil
}

func mapperOutput(q *Quotation) (*Output, error) {
	bid, err := strconv.ParseFloat(q.Usdbrl.Bid, 64)
	if err != nil {
		return nil, err
	}
	out := Output{Bid: bid}
	return &out, nil
}
