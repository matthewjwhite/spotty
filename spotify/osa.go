package spotify

import (
	"strings"

	"github.com/andybrewer/mack"
)

// Play plays a track in Spotify via OSA.
func Play(track Track) (err error) {
	_, err = mack.Tell("Spotify", "play track \""+track.URI+"\"")

	return
}

// Next plays the next track in Spotify via OSA.
// If executed in the context of an album/playlist that the user explicitly
// started from, Spotify will play the next song (next song in sequence,
// or shuffled). Else, will be random. In short, the context is not preserved in the DB.
func Next() (err error) {
	_, err = mack.Tell("Spotify", "play next track")

	return
}

// PlayPause toggles play/pause in Spotify via OSA.
func PlayPause() (err error) {
	_, err = mack.Tell("Spotify", "playpause")

	return
}

// Current returns the current track playing/paused in Spotify via OSA.
func Current() (Track, error) {
	ret, err := mack.Tell("Spotify",
		"get {spotify url, artist, name, album} of current track")

	parts := strings.SplitN(ret, ", ", 2)

	if err != nil {
		return Track{}, err
	}

	return Track{URI: parts[0], Data: parts[1]}, nil
}
