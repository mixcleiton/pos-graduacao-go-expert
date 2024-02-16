package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	URLQuotationService = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	ServerPort          = ":8080"
	QuotationEndPoint   = "/cotacao"
	pathDatabase        = "./quotation.db"
)

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

type QuotationResponse struct {
	Bid string `json:"bid"`
}

func main() {
	os.Remove(pathDatabase)

	db, err := sql.Open("sqlite3", pathDatabase)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
		create table quotation
		(
			code       text,
			codein     text,
			name       text,
			high       text,
			low        text,
			varbid     text,
			pctchange  text,
			bid        text,
			ask        text,
			timestamp  text,
			createdate text
		);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	http.HandleFunc(QuotationEndPoint, ProcessQuotation)
	http.ListenAndServe(ServerPort, nil)
}

func ProcessQuotation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	quotation, err := getQuotation(ctx)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = saveQuotation(ctx, quotation); err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := QuotationResponse{
		Bid: quotation.Usdbrl.Bid,
	}

	if !checkErrorTimeout(ctx) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
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

func saveQuotation(ctx context.Context, quotation *Quotation) error {
	db, err := sql.Open("sqlite3", pathDatabase)
	if err != nil {
		return fmt.Errorf("error to open conection with the database. cause: %w", err)
	}
	defer db.Close()

	stmt, err := db.Prepare("insert into main.quotation (code, codein, name, high, low, varbid, pctchange, bid, ask," +
		" \"timestamp\", createdate) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);")
	if err != nil {
		return fmt.Errorf("error to create statement. cause: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	usdToBrl := &quotation.Usdbrl

	_, err = stmt.ExecContext(ctx, usdToBrl.Code, usdToBrl.Codein, usdToBrl.Name, usdToBrl.High, usdToBrl.Low, usdToBrl.VarBid,
		usdToBrl.PctChange, usdToBrl.Bid, usdToBrl.Ask, usdToBrl.Timestamp, usdToBrl.CreateDate)
	if err != nil {
		checkErrorTimeout(ctx)
		return fmt.Errorf("error to execute statement. cause: %w", err)
	}

	return nil
}

func getQuotation(ctx context.Context) (*Quotation, error) {
	client := &http.Client{}

	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", URLQuotationService, nil)
	if err != nil {
		return nil, fmt.Errorf("error to consuming quotation service. cause: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		checkErrorTimeout(ctx)
		return nil, fmt.Errorf("error to execute request. cause: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error to read body. cause: %w", err)
	}

	var quotation Quotation
	err = json.Unmarshal(body, &quotation)
	if err != nil {
		return nil, fmt.Errorf("error to unmarshal body. cause: %w", err)
	}

	return &quotation, nil
}
