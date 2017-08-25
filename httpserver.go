package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func main() {
	port := flag.Int("port", 8000, "listen on which port")
	dir := flag.String("dir", ".", "directory to serve")
	flag.Parse()
	fmt.Printf("listening on port: %d\nserving directory: %s\n", *port, *dir)
	http.ListenAndServe(
		fmt.Sprintf(":%d", *port),
		handlers.LoggingHandler(
			os.Stderr,
			http.FileServer(http.Dir(*dir))))
}
