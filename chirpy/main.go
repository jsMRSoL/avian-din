package main

import (
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	filepathRoot := "."
	fs := http.FileServer(http.Dir(filepathRoot))

	mux.Handle("/", fs)

	corsMux := middlewareCors(mux)
	port := ":8080"

	var srv http.Server
	srv.Handler = corsMux
	srv.Addr = port

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().
			Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
