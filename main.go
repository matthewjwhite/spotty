package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/matthewjwhite/spotty/db"
	"github.com/matthewjwhite/spotty/spotify"
)

const (
	argError = 1 << iota
	spotifyError
	dbError
	homeError
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Must specify command.")
		os.Exit(argError)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(homeError)
	}

	db, err := db.New(homeDir + "/.spotty")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(dbError)
	}
	defer db.Close()

	if err = db.Init(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(dbError)
	}

	switch os.Args[1] {

	case "pp":
		err := spotify.PlayPause()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(spotifyError)
		}

	case "next":
		err := spotify.Next()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(spotifyError)
		}

	case "save":
		curr, err := spotify.Current()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(spotifyError)
		}

		err = db.Insert(curr)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(dbError)
		}

	case "play":
		if len(os.Args) <= 2 {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(argError)
		}

		track, err := db.Get(strings.Join(os.Args[2:], " "))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(dbError)
		}

		err = spotify.Play(track)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(spotifyError)
		}

	case "list":
		tracks, err := db.GetAll()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(dbError)
		}

		for _, t := range tracks {
			fmt.Println(t.Data)
		}

	default:
		fmt.Fprintln(os.Stderr, "Invalid command.")
		os.Exit(argError)
	}
}
