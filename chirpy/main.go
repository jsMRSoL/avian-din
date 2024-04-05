package main

import (
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	filepathRoot := "."
	fsHandle := http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))
	var apiConfig apiConfig
	mux.Handle("/app/*", apiConfig.middlewareMetricsInc(fsHandle))

	mux.HandleFunc("GET /api/healthz", healthEndPoint)
	mux.HandleFunc("POST /api/validate_chirp", validateChirp)
	mux.HandleFunc("GET /admin/metrics", apiConfig.getFsHits)
	mux.HandleFunc("/api/reset", apiConfig.resetFsHits)

	corsMux := middlewareCors(mux)
	port := ":8080"

	var srv http.Server
	srv.Handler = corsMux
	srv.Addr = port

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
