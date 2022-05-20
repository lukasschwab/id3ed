package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/mikkyang/id3-go"
	"github.com/tailscale/hujson"
)

// Meta stores the editable surface of an id3.Tagger.
type Meta struct {
	filename string
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Album    string `json:"album"`
	Year     string `json:"year"`
	Genre    string `json:"genre"`
	// TODO: look into Comments. Can probably just make an array, append.
}

// Extract meta fields from f.Tagger.
func From(f *id3.File) *Meta {
	return &Meta{
		Title:  f.Tagger.Title(),
		Artist: f.Tagger.Artist(),
		Album:  f.Tagger.Album(),
		Year:   f.Tagger.Year(),
		Genre:  f.Tagger.Genre(),
	}
}

type setter func(string)

// Write meta fields to f.Tagger.
func (meta *Meta) Write(f *id3.File) {
	// TODO: make it possible to clear fields. Distinguish between {"title": ""}
	// and {} (an update without the "title" field).
	maybe := func(s setter, value string) {
		if value != "" {
			s(value)
		}
	}
	maybe(f.Tagger.SetTitle, meta.Title)
	maybe(f.Tagger.SetArtist, meta.Artist)
	maybe(f.Tagger.SetAlbum, meta.Album)
	maybe(f.Tagger.SetYear, meta.Year)
	maybe(f.Tagger.SetGenre, meta.Genre)
}

func (meta *Meta) SolicitUpdates() (*Meta, error) {
	initial, err := json.MarshalIndent(meta, "", "\t")
	if err != nil {
		return nil, fmt.Errorf("error formatting initial JSON: %w", err)
	}

	data, err := meta.solicitUpdates(initial)
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

func (meta *Meta) solicitUpdates(initial []byte) ([]byte, error) {
	// getEnvDefault := func(env, def string) string {
	// 	if set := os.Getenv(env); set != "" {
	// 		return set
	// 	}
	// 	return def
	// }
	// shell := getEnvDefault("SHELL", defaultShell)
	// editor := getEnvDefault("EDITOR", defaultEditor)

	f, err := os.CreateTemp("", fmt.Sprintf("%s.*.json", meta.filename))
	if err != nil {
		return []byte{}, fmt.Errorf("error creating temp file: %w", err)
	}
	defer os.Remove(f.Name())
	if err := ioutil.WriteFile(f.Name(), initial, os.ModeAppend); err != nil {
		return []byte{}, fmt.Errorf("error initializing temp file: %w", err)
	}

	cmd := exec.Command("vim", f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error running editor: %w", err)
	}

	updated, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return nil, fmt.Errorf("error reading temp file: %w", err)
	}
	return updated, nil
}

// TODO: take settable fields as command-line arguments, and pre-set them.
func main() {
	filename := os.Args[1]
	file, err := id3.Open(filename)
	if err != nil {
		log.Fatalf("Parsing failed: %v", err)
	}
	defer file.Close()
	meta := From(file)
	updated, err := meta.SolicitUpdates()
	if err != nil {
		log.Fatalf("Updates failed: %v", err)
	}
	updated.Write(file)
	log.Printf("Updated metadata.")
}
