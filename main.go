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
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Must specify command.")
		os.Exit(int(argError))
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(int(homeError))
	}

	db, err := db.New(homeDir + "/.spotty")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(int(dbError))
	}
	defer db.Close()

	if err = db.Init(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(int(dbError))
	}

	switch os.Args[1] {

	case "pp":
		os.Exit(int(playPause()))

	case "next":
		os.Exit(int(next()))

	case "save":
		os.Exit(int(save(db)))

	case "play":
		os.Exit(int(play(db)))

	case "list":
		os.Exit(int(list(db)))

	default:
		fmt.Fprintln(os.Stderr, "Invalid command.")
		os.Exit(int(argError))
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
