package main

import "net/http"

func main() {
	fileServer := http.FileServer(http.Dir("./public"))
	mux := http.NewServeMux()
	mux.Handle("/", fileServer)
	http.ListenAndServe(":8080", mux)
}
