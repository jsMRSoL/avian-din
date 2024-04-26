package database

import (
	"os"
	"reflect"
	"sync"
	"testing"
)

func TestNewDB(t *testing.T) {
	type args struct {
		path string
	}
	// var db DB
	// db.path = "./database.db"
	db := DB{
		path: "./database.db",
		mu:   &sync.RWMutex{},
	}

	tests := []struct {
		name    string
		args    args
		want    *DB
		wantErr bool
	}{
		// test cases.
		{
			name: "can create database",
			args: args{
				path: "./database.db",
			},
			want:    &db,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDB(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func setup(t *testing.T) (*DB, error) {
	_, err := os.Stat("./database.db")
	if err == nil {
		err := os.Remove("./database.db")
		if err != nil {
			t.Error("couldn't remove db file")
			return nil, err
		}
	}
	db, err := NewDB("./database.db")
	if err != nil {
		t.Error("choked on setup function!")
		return nil, err
	}
	err = db.ensureDB()
	if err != nil {
		t.Error("couldn't initialize db structure")
		return nil, err
	}
	return db, nil
}

func TestDB_StoreChirp(t *testing.T) {
	type args struct {
		authorId int
		body     string
	}
	tests := []struct {
		name    string
		args    args
		want    Chirp
		wantErr bool
	}{
		// test cases.
		{
			name: "can create chirp 1",
			args: args{
				authorId: 1,
				body:     "This is a test!",
			},
			want: Chirp{
				Id:       1,
				Body:     "This is a test!",
				AuthorId: 1,
			},
			wantErr: false,
		},
		{
			name: "can create chirp 2",
			args: args{
				authorId: 2,
				body:     "This is a second test!",
			},
			want: Chirp{
				Id:       2,
				Body:     "This is a second test!",
				AuthorId: 2,
			},
			wantErr: false,
		},
	}
	// zero db
	db, err := setup(t)
	if err != nil {
		t.Error("choked on setup function!")
		return
	}
	//
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.StoreChirp(tt.args.body, tt.args.authorId)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"DB.StoreChirp() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DB.StoreChirp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_GetChirpByID(t *testing.T) {

	chirps := []Chirp{
		{Id: 1, Body: "This is the first one", AuthorId: 1},
		{Id: 2, Body: "This is the second one", AuthorId: 2},
		{Id: 3, Body: "This is the third one", AuthorId: 3},
		{Id: 4, Body: "This is the fourth one", AuthorId: 4},
		{Id: 5, Body: "This is the fifth one", AuthorId: 5},
	}
	tests := []struct {
		name    string
		msg     string
		id      int
		want    Chirp
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "Can get chirp 1",
			want:    Chirp{Id: 1, Body: "This is the first one", AuthorId: 1},
			id:      1,
			msg:     "This is the first one",
			wantErr: false,
		},
		{
			name:    "Can get chirp 2",
			want:    Chirp{Id: 2, Body: "This is the second one", AuthorId: 2},
			id:      2,
			msg:     "This is the second one",
			wantErr: false,
		},
	}
	// zero db
	db, err := setup(t)
	if err != nil {
		return
	}
	//
	for _, ch := range chirps {
		db.StoreChirp(ch.Body, ch.Id)
	}
	// done
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.GetChirp(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.GetChirp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DB.GetChirps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_GetChirps(t *testing.T) {
	type fields struct {
		path string
	}
	chirps := []Chirp{
		{Id: 1, Body: "This is the first one", AuthorId: 1},
		{Id: 2, Body: "This is the second one", AuthorId: 2},
		{Id: 3, Body: "This is the third one", AuthorId: 3},
		{Id: 4, Body: "This is the fourth one", AuthorId: 4},
		{Id: 5, Body: "This is the fifth one", AuthorId: 5},
	}
	tests := []struct {
		name    string
		fields  fields
		want    []Chirp
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Can get five chirps",
			fields: fields{
				path: "./database.db",
			},
			want:    chirps,
			wantErr: false,
		},
	}
	// zero db
	db, err := setup(t)
	if err != nil {
		t.Error("choked on setup function!")
		return
	}
	//
	// set up db contents
	for _, ch := range chirps {
		db.StoreChirp(ch.Body, ch.Id)
	}
	// done
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.GetChirps(false)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.GetChirps() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DB.GetChirps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_ensureDB(t *testing.T) {
	// test case 1
	// no db file exists
	_, err := os.Stat("./database.db")
	if err == nil {
		err := os.Remove("./database.db")
		if err != nil {
			t.Error("couldn't remove db file")
			return
		}
	}
	// set up variables
	path := "./database.db"
	db, err := NewDB(path)
	if err != nil {
		t.Error("choked on setting up new db (never got to ensure_DB!)")
		return
	}
	// test 1 begins
	if err := db.ensureDB(); err != nil {
		t.Error("Test 1 (dbfile does not exist): ensure DB returned an error")
		return
	}

	// test case 2: file already exists
	if err := db.ensureDB(); err != nil {
		t.Error("Test 2 (dbfile already exists): ensure DB returned an error")
		return
	}
}

func TestDB_loadDB(t *testing.T) {

	tests := []struct {
		name    string
		want    DBStructure
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "can load db",
			want:    DBStructure{map[int]Chirp{}},
			wantErr: false,
		},
	}
	// zero db
	db, err := setup(t)
	if err != nil {
		t.Error("choked on setup function!")
		return
	}
	//
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.loadDB()
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.loadDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DB.loadDB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_writeDB(t *testing.T) {
	type args struct {
		dbStructure DBStructure
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Can write to empty db",
			args: args{
				dbStructure: DBStructure{
					map[int]Chirp{
						1: {
							Id:   1,
							Body: "This is the first one!",
						},
						2: {
							Id:   1,
							Body: "This is the second one!",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	// zero db
	db, err := setup(t)
	if err != nil {
		return
	}
	// New empty db is now setup...
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.writeDB(tt.args.dbStructure); (err != nil) != tt.wantErr {
				t.Errorf("DB.writeDB() error = %v, wantErr %v", err, tt.wantErr)
			}
			dbStructFromFile, _ := db.loadDB()
			if !reflect.DeepEqual(dbStructFromFile, tt.args.dbStructure) {
				t.Errorf(
					"loadDB got %v\nexpected: %v",
					dbStructFromFile,
					tt.args.dbStructure,
				)
			}
		})
	}
}
