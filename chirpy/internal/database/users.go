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

type UserDB struct {
	path string
	mu   *sync.RWMutex
}

type UserDBStructure struct {
	Users map[int]User `json:"users"`
}

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

// NewUserDB creates a new database connection
// and creates the database file if it doesn't exist
func NewUserDB(path string) (*UserDB, error) {
	db := &UserDB{
		path: path,
		mu:   &sync.RWMutex{},
	}
	err := db.ensureUserDB()
	return db, err
}

func (db *UserDB) AddUser(body string) (User, error) {
	dbStruct, err := db.loadUserDB()
	if err != nil {
		return User{}, err
	}

	id := len(dbStruct.Users) + 1
	user := User{
		Id:    id,
		Email: body,
	}

	dbStruct.Users[id] = user
	err = db.writeUserDB(dbStruct)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *UserDB) GetUsers() ([]User, error) {
	dbStruct, err := db.loadUserDB()
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(dbStruct.Users))
	for _, v := range dbStruct.Users {
		users = append(users, v)
	}

	slices.SortFunc(users, sortUserSlice)
	return users, nil
}

func sortUserSlice(a, b User) int {
	return cmp.Compare(a.Id, b.Id)
}

func (db *UserDB) GetUser(id int) (User, error) {
	dbStruct, err := db.loadUserDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStruct.Users[id]
	if !ok {
		return User{}, errors.New(
			fmt.Sprintf("Database does not contain User ID: %d", id),
		)
	}
	return user, nil
}

// ensureUserDB creates a new database file if it doesn't exist
func (db *UserDB) ensureUserDB() error {
	if _, err := os.ReadFile(db.path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			dbStruct := UserDBStructure{
				Users: map[int]User{},
			}
			err := db.writeUserDB(dbStruct)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// loadUserDB reads the database file into memory
func (db *UserDB) loadUserDB() (UserDBStructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbStruct := UserDBStructure{}
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

func (db *UserDB) writeUserDB(dbStructure UserDBStructure) error {
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
