package spotify

import "strings"

// Track represents a Spotify track.
type Track struct {
	URI  string
	Data string // Artist, Song, Album
}

// String yields a string representation.
func (t Track) String() string {
	return strings.Join([]string{t.URI, t.Data}, ",")
}
