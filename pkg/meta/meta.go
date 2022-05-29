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

// SolicitUpdates to meta from the user.
func (meta *Meta) SolicitUpdates(partial *Partial, comment bool) (*Meta, error) {
	meta.apply(partial)
	initial, err := meta.Format()
	if err != nil {
		return nil, err
	}

	if comment {
		initial = append(
			[]byte(fmt.Sprintf("// File: %s\n", meta.filename)),
			initial...,
		)
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

func isSet(s *string) bool {
	return s != nil && *s != ""
}

func (meta *Meta) apply(partial *Partial) {
	if partial == nil {
		return
	}
	if isSet(partial.Title) {
		meta.Title = *partial.Title
	}
	if isSet(partial.Artist) {
		meta.Artist = *partial.Artist
	}
	if isSet(partial.Album) {
		meta.Album = *partial.Album
	}
	if isSet(partial.Year) {
		meta.Year = *partial.Year
	}
	if isSet(partial.Genre) {
		meta.Genre = *partial.Genre
	}
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

// Format meta as JSON.
func (meta *Meta) Format() ([]byte, error) {
	data, err := json.MarshalIndent(meta, "", "\t")
	if err != nil {
		return nil, fmt.Errorf("error formatting initial JSON: %w", err)
	}
	return data, err
}

// Close the file metadata, flushing changes.
func (meta *Meta) Close() {
	meta.file.Close()
}

// Partial metadata for pre-filling an edit struct.
type Partial struct {
	Title  *string
	Artist *string
	Album  *string
	Year   *string
	Genre  *string
}

func (p *Partial) Mask(other *Partial) {
	if other == nil {
		return
	}
	if !isSet(p.Title) {
		p.Title = other.Title
	}
	if !isSet(p.Artist) {
		p.Artist = other.Artist
	}
	if !isSet(p.Album) {
		p.Album = other.Album
	}
	if !isSet(p.Year) {
		p.Year = other.Year
	}
	if !isSet(p.Genre) {
		p.Genre = other.Genre
	}
}
