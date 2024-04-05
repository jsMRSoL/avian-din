package main

import "strings"

func cleanChirp(chirp string) (clean_chirp string) {

	var cleaned []string

	words := strings.Fields(chirp)
	for _, wd := range words {
		cleaned = append(cleaned, clean(wd))
	}
	return strings.Join(cleaned, " ")
}

func clean(word string) string {
	lowered := strings.ToLower(word)
	if lowered == "kerfuffle" || lowered == "sharbert" || lowered == "fornax" {
		return "****"
	}
	return word
}
