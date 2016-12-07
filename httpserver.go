package main

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func main() {
	http.ListenAndServe(
		":8000",
		handlers.LoggingHandler(
			os.Stderr,
			http.FileServer(http.Dir("."))))
}
