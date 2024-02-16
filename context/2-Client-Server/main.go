package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("Request started")
	defer log.Println("Request finished")
	select {
	case <-time.After(5 * time.Second):
		// print in command line
		log.Println("Request processed with success")
		w.Write([]byte("Request processed with success"))
	case <-ctx.Done():
		log.Println("Request cancelled for client")
		//http.Error(w, "Request cancelled for client", http.StatusRequestTimeout)
	}
}
