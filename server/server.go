package main

import "net/http"

func main() {
	http.HandleFunc("/cotacao", GetQuotation)
	http.ListenAndServe(":8000", nil)
}

func GetQuotation(w http.ResponseWriter, r *http.Request) {}
