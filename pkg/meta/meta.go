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
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
	Year   string `json:"year"`
	Genre  string `json:"genre"`
	// TODO: look into Comments. Can probably just make an array, append.
}

// From the existing ID3 tags on f, constructs a new Meta.
func From(f *id3.File) *Meta {
	return &Meta{
		Title:  f.Tagger.Title(),
		Artist: f.Tagger.Artist(),
		Album:  f.Tagger.Album(),
		Year:   f.Tagger.Year(),
		Genre:  f.Tagger.Genre(),
	}
}

// Write meta fields to f.Tagger.
func (meta *Meta) Write(f *id3.File) {
	// TODO: make it possible to clear fields. Distinguish between {"title": ""}
	// and {} (an update without the "title" field).
	maybe := func(set func(string), value string) {
		if value != "" {
			set(value)
		}
	}
	maybe(f.Tagger.SetTitle, meta.Title)
	maybe(f.Tagger.SetArtist, meta.Artist)
	maybe(f.Tagger.SetAlbum, meta.Album)
	maybe(f.Tagger.SetYear, meta.Year)
	maybe(f.Tagger.SetGenre, meta.Genre)
}

// SolicitUpdates to meta from the user.
func (meta *Meta) SolicitUpdates() (*Meta, error) {
	initial, err := json.MarshalIndent(meta, "", "\t")
	if err != nil {
		return nil, fmt.Errorf("error formatting initial JSON: %w", err)
	}

	data, err := editor.GetUpdates(initial)
	if err != nil {
		return nil, err
	}

	if data, err = hujson.Standardize(data); err != nil {
		return nil, fmt.Errorf("error standardizing user input: %w", err)
	}

	updated := new(Meta)
	if err := json.Unmarshal(data, updated); err != nil {
		return nil, fmt.Errorf("error parsing user input: %w", err)
	}
	return updated, nil
}
