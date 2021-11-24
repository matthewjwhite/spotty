package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/matthewjwhite/spotty/db"
	"github.com/matthewjwhite/spotty/spotify"
)

const (
	success  exitStatus = 0
	argError exitStatus = 1 << iota
	spotifyError
	dbError
	homeError
)

type exitStatus int

func main() {
	os.Exit(int(run()))
}

func run() exitStatus {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Must specify command.")
		return argError
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return homeError
	}

	db, err := db.New(homeDir + "/.spotty")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return dbError
	}
	defer db.Close()

	if err = db.Init(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return dbError
	}

	switch os.Args[1] {

	case "pp":
		return playPause()

	case "next":
		return next()

	case "save":
		return save(db)

	case "play":
		return play(db)

	case "list":
		return list(db)

	case "get":
		return get()

	default:
		fmt.Fprintln(os.Stderr, "Invalid command.")
		return argError
	}
}

func playPause() exitStatus {
	err := spotify.PlayPause()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return spotifyError
	}

	return success
}

func get() exitStatus {
	track, err := spotify.Current()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return spotifyError
	}

	fmt.Println(track)

	return success
}

func next() exitStatus {
	err := spotify.Next()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return spotifyError
	}

	return success
}

func save(db db.DB) exitStatus {
	curr, err := spotify.Current()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return spotifyError
	}

	err = db.Insert(curr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return dbError
	}

	return success
}

func list(db db.DB) exitStatus {
	tracks, err := db.GetAll()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return dbError
	}

	for _, t := range tracks {
		fmt.Println(t.Data)
	}

	return success
}

func play(db db.DB) exitStatus {
	if len(os.Args) <= 2 {
		fmt.Fprintln(os.Stderr, "Must specify search.")
		return argError
	}

	track, err := db.Get(strings.Join(os.Args[2:], " "))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return dbError
	}

	err = spotify.Play(track)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return spotifyError
	}

	return success
}
