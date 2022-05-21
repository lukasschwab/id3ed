package meta

import (
	"encoding/json"
	"fmt"

	"github.com/lukasschwab/id3ed/pkg/editor"
	"github.com/mikkyang/id3-go"
	"github.com/tailscale/hujson"
)

// Meta stores the editable surface of an id3.Tagger.
type Meta struct {
	filename string
	file     *id3.File

	Title  string `json:"title"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
	Year   string `json:"year"`
	Genre  string `json:"genre"`
	// TODO: look into Comments. Can probably just make an array, append.
}

// From file with filename, extracts metadata. Callers must call Close().
func From(filename string) (*Meta, error) {
	file, err := id3.Open(filename)
	if err != nil {
		// log.Fatalf("Parsing failed: %v", err)
		return nil, err
	}
	return from(filename, file), nil
}

// From the existing ID3 tags on f, constructs a new Meta.
func from(filename string, file *id3.File) *Meta {
	return &Meta{
		filename: filename,
		file:     file,
		Title:    file.Tagger.Title(),
		Artist:   file.Tagger.Artist(),
		Album:    file.Tagger.Album(),
		Year:     file.Tagger.Year(),
		Genre:    file.Tagger.Genre(),
	}
}

func (meta *Meta) Close() {
	meta.file.Close()
}

// Write meta fields to f.Tagger.
func (meta *Meta) Write() {
	// TODO: make it possible to clear fields. Distinguish between {"title": ""}
	// and {} (an update without the "title" field).
	maybe := func(set func(string), value string) {
		if value != "" {
			set(value)
		}
	}
	maybe(meta.file.Tagger.SetTitle, meta.Title)
	maybe(meta.file.Tagger.SetArtist, meta.Artist)
	maybe(meta.file.Tagger.SetAlbum, meta.Album)
	maybe(meta.file.Tagger.SetYear, meta.Year)
	maybe(meta.file.Tagger.SetGenre, meta.Genre)
}

func (meta *Meta) Format() ([]byte, error) {
	data, err := json.MarshalIndent(meta, "", "\t")
	if err != nil {
		return nil, fmt.Errorf("error formatting initial JSON: %w", err)
	}
	return data, err
}

// SolicitUpdates to meta from the user.
func (meta *Meta) SolicitUpdates() (*Meta, error) {
	initial, err := meta.Format()
	if err != nil {
		return nil, err
	}

	data, err := editor.GetUpdates(initial)
	if err != nil {
		return nil, err
	}

	if data, err = hujson.Standardize(data); err != nil {
		return nil, fmt.Errorf("error standardizing user input: %w", err)
	}

	updated := &Meta{
		filename: meta.filename,
		file:     meta.file,
	}
	if err := json.Unmarshal(data, updated); err != nil {
		return nil, fmt.Errorf("error parsing user input: %w", err)
	}
	return updated, nil
}
