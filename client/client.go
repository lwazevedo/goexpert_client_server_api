package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const baseUrl = "http://localhost:8080"
const maxTime = 300
const fileName = "cotacao.txt"

type Output struct {
	Bid float64 `json:"bid"`
}

func main() {
	out := GetQuotation()
	saveQuotationFile(out)

}
func GetQuotation() *Output {
	ctx, cancel := context.WithTimeout(context.TODO(), maxTime*time.Millisecond)
	defer cancel()
	url := fmt.Sprintf("%s%s", baseUrl, "/cotacao")
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		println("Não foi criar a requisição para o server")
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		println("Não foi possível realizar a requisição para o server")
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		println("Não foi possível ler o corpo da requisição do server")
		panic(err)
	}
	var out Output
	err = json.Unmarshal(body, &out)
	if err != nil {
		println("Não foi possível converter o corpo da requisição do server")
		panic(err)
	}
	return &out
}

func saveQuotationFile(out *Output) {
	file, err := os.Create(fileName)
	if err != nil {
		println("Não foi possível criar o arquivo cotacao.txt")
		panic(err)
	}
	defer file.Close()
	_, err = file.Write([]byte(fmt.Sprintf("Dólar: %.2f", out.Bid)))
	if err != nil {
		println("Não foi possível escrever no arquivo cotacao.txt")
		panic(err)
	}
	fmt.Println("Arquivo de cotação criado com sucesso.")
}
