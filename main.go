package main

import "net/http"

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from lommeulken"))
	})

	http.ListenAndServe(":8080", mux)
}
