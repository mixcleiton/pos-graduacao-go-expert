package main

import (
	"context"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	urlQuotation = "http://localhost:8080/cotacao"
)

type QuotationResponse struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", urlQuotation, nil)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		checkErrorTimeout(ctx)
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var quotation QuotationResponse
	err = json.Unmarshal(body, &quotation)
	if err != nil {
		panic(err)
	}

	file, err := os.Create("./cotacao.txt")
	if err != nil {
		panic(err)
	}
	tmp := template.New("cotacao.txt")
	tmp, _ = tmp.Parse("Dolar: {{.Bid}}")
	err = tmp.Execute(file, &quotation)
	if err != nil {
		panic(err)
	}
}

func checkErrorTimeout(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		log.Println("context cancelled")
		return true
	default:
		return false
	}
}
