package database

import (
	"cmp"
	"encoding/json"
	"log"
	"os"
	"slices"
	"sync"
)

type DB struct {
	path string
	mu   sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	var db DB
	db.path = path
	if err := db.ensureDB(); err != nil {
		log.Println("ensure DB failed")
		return &db, err
	}
	return &db, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStruct.Chirps) + 1
	chirp := Chirp{
		Id:   id,
		Body: body,
	}

	dbStruct.Chirps[id] = chirp
	db.writeDB(dbStruct)

	return chirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		log.Println("Couldn't load db")
		return nil, err
	}

	var chirps []Chirp
	for _, v := range dbStruct.Chirps {
		chirps = append(chirps, v)
	}

	slices.SortFunc(chirps, sortChirpSlice)
	return chirps, nil
}

func sortChirpSlice(a, b Chirp) int {
	return cmp.Compare(a.Id, b.Id)
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	if _, err := os.ReadFile(db.path); err == nil {
		return nil
	}

	var dbStruct DBStructure
	empty := make(map[int]Chirp)
	dbStruct.Chirps = empty

	bytes, err := json.Marshal(dbStruct)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, bytes, 0666)
	if err != nil {
		log.Println("couldn't create db file")
		return err
	}
	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {

	db.mu.RLock()
	defer db.mu.RUnlock()

	data, err := os.ReadFile(db.path)
	if err != nil {
		log.Println(err)
	}
	var dbStruct DBStructure
	if err := json.Unmarshal(data, &dbStruct); err != nil {
		return DBStructure{}, err
	}

	return dbStruct, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {

	db.mu.Lock()
	defer db.mu.Unlock()

	bytes, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, bytes, 0666)
	if err != nil {
		return err
	}

	return nil
}

// for use in testing
func (db *DB) zeroDB() error {
	err := os.WriteFile(db.path, []byte{}, 0666)
	if err != nil {
		return err
	}
	return nil
}
