package main

import (
	"net/http"
	"os"
)

func main() {
	http.Handle("/", helloHandler())
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

func helloHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello!"))
	})
}
