// Package db provides helper functionality for interacting with the database.
package db

import (
	"database/sql"
	"fmt"

	"github.com/matthewjwhite/spotty/spotify"
	_ "github.com/mattn/go-sqlite3"
)

const (
	tracksTable = "tracks"
	uriCol      = "uri"
	dataCol     = "data"
)

// DB is a wrapper for *sql.DB, for the purpose of "inheriting" methods like Close
// that can be used alongside the custom ones here.
type DB struct {
	*sql.DB
}

// New returns DB struct, given a valid path to the local database file.
func New(path string) (DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return DB{}, err
	}

	return DB{db}, nil
}

// Init initializes the database with the tracks table.
func (d DB) Init() error {
	statement, err := d.Prepare(
		fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s "+
			"(id INTEGER PRIMARY KEY, "+
			"%s TEXT NOT NULL UNIQUE, %s TEXT NOT NULL UNIQUE)",
			tracksTable, uriCol, dataCol))
	if err != nil {
		return err
	}

	_, err = statement.Exec()
	if err != nil {
		return err
	}

	return nil
}

// Insert inserts a track into the database.
func (d DB) Insert(t spotify.Track) error {
	statement, err := d.Prepare(
		fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES (?, ?)",
			tracksTable, uriCol, dataCol))
	if err != nil {
		return err
	}

	_, err = statement.Exec(t.URI, t.Data)
	if err != nil {
		return err
	}

	return nil
}

// Get retrieves the first track matching the specified pattern.
func (d DB) Get(match string) (spotify.Track, error) {
	row := d.QueryRow(
		fmt.Sprintf("SELECT %s, %s FROM %s WHERE %s LIKE ?",
			uriCol, dataCol, tracksTable, dataCol),
		"%"+match+"%")

	var uri string
	var data string
	err := row.Scan(&uri, &data)
	if err != nil {
		return spotify.Track{}, err
	}

	return spotify.Track{URI: uri, Data: data}, nil
}

// GetAll gets all tracks.
func (d DB) GetAll() ([]spotify.Track, error) {
	row, err := d.Query(
		fmt.Sprintf("SELECT %s, %s FROM %s",
			uriCol, dataCol, tracksTable))
	if err != nil {
		return nil, err
	}
	defer row.Close()

	tracks := make([]spotify.Track, 0)

	for row.Next() {
		var uri string
		var data string

		err = row.Scan(&uri, &data)
		if err != nil {
			return nil, err
		}

		tracks = append(tracks, spotify.Track{URI: uri, Data: data})
	}

	return tracks, nil
}
