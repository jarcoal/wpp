package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()

	http.HandleFunc("/", html)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func html(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
