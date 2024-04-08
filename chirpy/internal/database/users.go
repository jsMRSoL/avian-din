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
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserDB struct {
	path string
	mu   *sync.RWMutex
}

type UserDBStructure struct {
	Users         map[int]RegisteredUser `json:"users"`
	Addrs         map[string]int         `json:"addrs"`
	RevokedTokens map[string]time.Time   `json:"revoked_tokens"`
}

type RegisteredUser struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	HashedPw    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type User struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type SignedUser struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	AccessToken  string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	IsChirpyRed  bool   `json:"is_chirpy_red"`
}

func (u *User) ToSignedUser(accessToken, refreshToken string) SignedUser {
	return SignedUser{
		Id:           u.Id,
		Email:        u.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		IsChirpyRed:  u.IsChirpyRed,
	}
}

func (rg *RegisteredUser) toUser() User {
	return User{
		Id:          rg.Id,
		Email:       rg.Email,
		IsChirpyRed: rg.IsChirpyRed,
	}
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

func (db *UserDB) AddUser(body string, passwd string) (User, error) {
	dbStruct, err := db.loadUserDB()
	if err != nil {
		return User{}, err
	}

	// Check if user already registered
	_, registered := dbStruct.Addrs[body]
	if registered {
		return User{}, errors.New("User is already registered")
	}

	// Hash password
	pw, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	id := len(dbStruct.Users) + 1
	user := RegisteredUser{
		Id:          id,
		Email:       body,
		HashedPw:    string(pw),
		IsChirpyRed: false,
	}

	dbStruct.Users[id] = user
	dbStruct.Addrs[body] = id
	err = db.writeUserDB(dbStruct)
	if err != nil {
		return User{}, err
	}

	return user.toUser(), nil
}

func (db *UserDB) UpdateUser(id int, email, passwd string) (User, error) {
	dbStruct, err := db.loadUserDB()
	if err != nil {
		return User{}, err
	}

	//////////////////////////
	// Delete old entries first
	// delete old email from Addrs db
	oldUserData, ok := dbStruct.Users[id]
	if !ok {
		return User{}, errors.New(
			fmt.Sprintf("Database does not contain User ID: %d", id),
		)
	}
	delete(dbStruct.Addrs, oldUserData.Email)
	/////////////////////////
	// Hash password
	pw, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	// Should I be copying the token across?
	user := RegisteredUser{
		Id:       id,
		Email:    email,
		HashedPw: string(pw),
	}
	/////////////////////////
	// Overwrite old entry in Users
	dbStruct.Users[id] = user
	dbStruct.Addrs[email] = id
	err = db.writeUserDB(dbStruct)
	if err != nil {
		return User{}, err
	}

	return user.toUser(), nil
}

func (db *UserDB) UpgradeUser(userId int) error {

	dbStruct, err := db.loadUserDB()
	if err != nil {
		log.Println("--> DB: Could not load db")
		return err
	}

	user := dbStruct.Users[userId]
	user.IsChirpyRed = true
	dbStruct.Users[userId] = user

	err = db.writeUserDB(dbStruct)
	if err != nil {
		log.Println("--> DB: Could not write db")
		return err
	}

	log.Printf("--> DB: upgrading user %d to Chirpy Red", user.Id)
	return nil
}

func (db *UserDB) GetUsers() ([]User, error) {
	dbStruct, err := db.loadUserDB()
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(dbStruct.Users))
	for _, v := range dbStruct.Users {
		users = append(users, v.toUser())
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
	return user.toUser(), nil
}

func (db *UserDB) GetUserId(email string) (int, error) {
	dbStruct, err := db.loadUserDB()
	if err != nil {
		return 0, err
	}

	id, ok := dbStruct.Addrs[email]
	if !ok {
		return 0, errors.New("Could not get userID")
	}

	return id, nil
}

func (db *UserDB) GetUserDetails(
	userId int,
) (email, hashedPW string, isChirpyRed bool, err error) {
	dbStruct, err := db.loadUserDB()
	if err != nil {
		return "", "", false, err
	}

	user, ok := dbStruct.Users[userId]
	if !ok {
		return "", "", false, errors.New("Could not get userID")
	}

	return user.Email, user.HashedPw, user.IsChirpyRed, nil
}

func (db *UserDB) GetUserPassword(id int) (string, error) {
	dbStruct, err := db.loadUserDB()
	if err != nil {
		return "", err
	}

	user, ok := dbStruct.Users[id]
	if !ok {
		return "", errors.New("Could not get userID")
	}

	return user.HashedPw, nil
}

func (db *UserDB) AuthenticateUser(email string, password string) (User, error) {
	userID, err := db.GetUserId(email)
	if err != nil {
		return User{}, err
	}

	hash, err := db.GetUserPassword(userID)
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return User{}, err
	}

	_, _, isChirpyRed, err := db.GetUserDetails(userID)
	if err != nil {
		return User{}, err
	}
	user := User{
		Id:          userID,
		Email:       email,
		IsChirpyRed: isChirpyRed,
	}
	return user, nil
}

func (db *UserDB) ensureUserDB() error {
	if _, err := os.ReadFile(db.path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			dbStruct := UserDBStructure{
				Users:         map[int]RegisteredUser{},
				Addrs:         map[string]int{},
				RevokedTokens: map[string]time.Time{},
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

func (db *UserDB) AddRevokedToken(refreshToken string) error {
	dbStruct, err := db.loadUserDB()
	if err != nil {
		return err
	}

	dbStruct.RevokedTokens[refreshToken] = time.Now()

	err = db.writeUserDB(dbStruct)
	if err != nil {
		return err
	}

	return nil
}

func (db *UserDB) IsRevoked(refreshToken string) (bool, error) {
	dbStruct, err := db.loadUserDB()
	if err != nil {
		return false, err
	}

	_, ok := dbStruct.RevokedTokens[refreshToken]
	if ok {
		return true, nil
	}

	return false, nil
}
