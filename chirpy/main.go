package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/jsMRSoL/avian-din/internal/database"
)

func main() {

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if *dbg == true {
		log.Println("In debug mode.................")
		os.Remove("storage.db")
		os.Remove("users.db")
	}

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

	/// Get env variable
	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")
	polkaApikey := os.Getenv("POLKA_APIKEY")

	apiConfig := apiConfig{
		chirpsDB:    chirpsDB,
		userDB:      userDB,
		secret:      jwtSecret,
		polkaApikey: polkaApikey,
	}

	mux := http.NewServeMux()

	filepathRoot := "."
	fsHandle := http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/*", apiConfig.middlewareMetricsInc(fsHandle))

	mux.HandleFunc("POST /api/users", apiConfig.addUser)
	mux.HandleFunc("PUT /api/users", apiConfig.updateUser)
	mux.HandleFunc("POST /api/login", apiConfig.loginUser)
	mux.HandleFunc("POST /api/refresh", apiConfig.refreshAccessToken)
	mux.HandleFunc("POST /api/revoke", apiConfig.revokeRefreshToken)
	mux.HandleFunc("GET /api/chirps", apiConfig.getChirps)
	mux.HandleFunc("POST /api/chirps", apiConfig.postChirp)
	mux.HandleFunc("DELETE /api/chirps/{ID}", apiConfig.deleteChirp)
	mux.HandleFunc("GET /api/chirps/{ID}", apiConfig.getChirpByID)

	mux.HandleFunc("POST /api/polka/webhooks", apiConfig.upgradeUser)

	mux.HandleFunc("GET /api/healthz", healthEndPoint)
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
