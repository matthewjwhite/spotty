package spotify

// Track represents a Spotify track.
type Track struct {
	URI  string
	Data string // Artist, Song, Album
}

// String yields a string representation.
func (t Track) String() string {
	return t.Data
}
