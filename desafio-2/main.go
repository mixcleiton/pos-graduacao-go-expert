package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	CEP          = "01153000"
	URLBrasilApi = "https://brasilapi.com.br/api/cep/v1/" + CEP
	URLViaCep    = "http://viacep.com.br/ws/" + CEP + "/json/"
)

type RequestChanel struct {
	Body []byte
	Url  string
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	ch := make(chan RequestChanel)

	go func() {
		bodyReq, err := createRequestApi(ctx, URLViaCep)
		if err != nil {
			fmt.Printf("error to request 2. cause: %s", err)
			return
		}
		ch <- RequestChanel{
			Body: bodyReq,
			Url:  URLViaCep,
		}
	}()

	go func() {
		bodyReq, err := createRequestApi(ctx, URLBrasilApi)
		if err != nil {
			fmt.Printf("error to request 1. cause: %s", err)
			return
		}
		ch <- RequestChanel{
			Body: bodyReq,
			Url:  URLBrasilApi,
		}
	}()

	select {
	case request := <-ch:
		fmt.Println("request API: " + request.Url)
		fmt.Println("request result: " + string(request.Body))
	case <-time.After(time.Second * 1):
		fmt.Println("context canceled")
	}
}

func createRequestApi(ctx context.Context, url string) ([]byte, error) {
	cli := http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error to create request. cause: %w", err)
	}

	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error to process request. cause: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error to read response body")
	}

	return body, nil
}
