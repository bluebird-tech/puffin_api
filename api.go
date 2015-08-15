package main

import (
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello!"))
}
