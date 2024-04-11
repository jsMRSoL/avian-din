package database

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"sync"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}

func (db *DB) StoreChirp(body string, authorId int) (Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStruct.Chirps) + 1
	chirp := Chirp{
		Id:       id,
		Body:     body,
		AuthorId: authorId,
	}

	dbStruct.Chirps[id] = chirp
	err = db.writeDB(dbStruct)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) DeleteChirp(chirpId int) error {
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}

	delete(dbStruct.Chirps, chirpId)

	err = db.writeDB(dbStruct)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetChirps(desc bool) ([]Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStruct.Chirps))
	for _, v := range dbStruct.Chirps {
		chirps = append(chirps, v)
	}

	if desc {
		slices.SortFunc(chirps, sortChirpSliceDesc)
	} else {
		slices.SortFunc(chirps, sortChirpSliceAsc)
	}
	return chirps, nil
}

func sortChirpSliceAsc(a, b Chirp) int {
	return cmp.Compare(a.Id, b.Id)
}
func sortChirpSliceDesc(a, b Chirp) int {
	return cmp.Compare(b.Id, a.Id)
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStruct.Chirps[id]
	if !ok {
		return Chirp{}, errors.New(
			fmt.Sprintf("Database does not contain Chirp ID: %d", id),
		)
	}
	return chirp, nil
}

func (db *DB) ChirpsByAuthorID(authorId int, desc bool) ([]Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := []Chirp{}
	for _, chirp := range dbStruct.Chirps {
		if chirp.AuthorId == authorId {
			chirps = append(chirps, chirp)
		}
	}

	if desc {
		slices.SortFunc(chirps, sortChirpSliceDesc)
	} else {
		slices.SortFunc(chirps, sortChirpSliceAsc)
	}

	return chirps, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	if _, err := os.ReadFile(db.path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			dbStruct := DBStructure{
				Chirps: map[int]Chirp{},
			}
			err := db.writeDB(dbStruct)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbStruct := DBStructure{}
	data, err := os.ReadFile(db.path)
	if err != nil {
		log.Println(err)
		return dbStruct, err
	}
	if err := json.Unmarshal(data, &dbStruct); err != nil {
		log.Println(err)
		return dbStruct, err
	}

	return dbStruct, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	bytes, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, bytes, 0600)
	if err != nil {
		return err
	}
	return nil
}
