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
	var db DB
	db.path = "./database.db"

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

func setup() error {
	os.Remove("./database.db")
	_, err := NewDB("./database.db")
	if err != nil {
		return err
	}
	return nil
}

func TestDB_CreateChirp(t *testing.T) {
	type fields struct {
		path string
		mu   sync.RWMutex
	}
	type args struct {
		body string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Chirp
		wantErr bool
	}{
		// test cases.
		{
			name: "can create a chirp",
			fields: fields{
				path: "./database.db",
				mu:   sync.RWMutex{},
			},
			args: args{
				body: "This is a test!",
			},
			want: Chirp{
				Id:   1,
				Body: "This is a test!",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		// zero db
		err := setup()
		if err != nil {
			t.Error("choked on setup function!")
			return
		}
		//
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{
				path: tt.fields.path,
				mu:   tt.fields.mu,
			}
			got, err := db.CreateChirp(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"DB.CreateChirp() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DB.CreateChirp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_GetChirps(t *testing.T) {
	type fields struct {
		path string
	}
	chirps := []Chirp{
		{Id: 1, Body: "This is the first one"},
		{Id: 2, Body: "This is the second one"},
		{Id: 3, Body: "This is the third one"},
		{Id: 4, Body: "This is the fourth one"},
		{Id: 5, Body: "This is the fifth one"},
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
	for _, tt := range tests {
		// zero db
		err := setup()
		if err != nil {
			t.Error("choked on setup function!")
			return
		}
		//
		t.Run(tt.name, func(t *testing.T) {
			// set up db contents
			db := &DB{
				path: tt.fields.path,
			}
			db.ensureDB()
			for _, ch := range chirps {
				db.CreateChirp(ch.Body)
			}
			// done
			got, err := db.GetChirps()
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
	err := os.Remove("./database.db")
	if err != nil {
		t.Error("SETUP: Couldn't remove db file.")
		return
	}
	// set up variables
	path := "./database.db"
	var db DB
	db.path = path
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

	type fields struct {
		path string
		mu   sync.RWMutex
	}
	tests := []struct {
		name    string
		fields  fields
		want    DBStructure
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "can load db",
			fields: fields{
				path: "./database.db",
			},
			want:    DBStructure{map[int]Chirp{}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		// zero db
		err := setup()
		if err != nil {
			t.Error("choked on setup function!")
			return
		}
		//
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{
				path: tt.fields.path,
			}
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
	type fields struct {
		path string
		mu   sync.RWMutex
	}
	type args struct {
		dbStructure DBStructure
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Can write to empty db",
			fields: fields{
				path: "./database.db",
			},
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
	for _, tt := range tests {
		// zero db
		err := setup()
		if err != nil {
			t.Error("choked on setup function!")
			return
		}
		// New empty db is now setup...
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{
				path: tt.fields.path,
			}
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
