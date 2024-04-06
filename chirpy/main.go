package main

import (
	"github.com/jsMRSoL/avian-din/internal/database"
	"log"
	"net/http"
)

func main() {

	path := "storage.db"
	chirpsDB, err := database.NewDB(path)
	if err != nil {
		log.Printf("Error creating DB: %s", err)
		return
	}

	userDB_path := "users.db"
	userDB, err := database.NewUserDB(userDB_path)
	if err != nil {
		log.Printf("Error creating DB: %s", err)
		return
	}

	apiConfig := apiConfig{
		chirpsDB: chirpsDB,
		userDB:   userDB,
	}

	mux := http.NewServeMux()

	filepathRoot := "."
	fsHandle := http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/*", apiConfig.middlewareMetricsInc(fsHandle))

	mux.HandleFunc("GET /api/healthz", healthEndPoint)
	mux.HandleFunc("POST /api/chirps", apiConfig.postChirp)
	mux.HandleFunc("POST /api/users", apiConfig.addUser)
	mux.HandleFunc("GET /api/chirps", apiConfig.getChirps)
	mux.HandleFunc("GET /api/chirps/{ID}", apiConfig.getChirpByID)
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
