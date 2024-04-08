package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) upgradeUser(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Event string
		Data  map[string]int
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusOK)
		return
	}

	log.Println("--> /api/polka/webhooks: received user.upgraded event ")
	data := params.Data
	userId, ok := data["user_id"]
	if !ok {
		log.Printf("--> /api/polka/webhooks: no user_id obtained from %v", data)
		w.WriteHeader(http.StatusOK)
		return
	}

	err = cfg.userDB.UpgradeUser(userId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Something went wrong")
		return
	}

	w.WriteHeader(http.StatusOK)
}
