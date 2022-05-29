package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/lukasschwab/id3ed/pkg/meta"

	"github.com/irlndts/go-discogs"
)

func open(filename string) *meta.Meta {
	data, err := meta.From(filename)
	if err != nil {
		log.Fatalf("Parsing failed: %v", err)
		return nil
	}
	return data
}

func inspectMetadata(filename string) {
	metadata := open(filename)
	defer metadata.Close()

	data, err := metadata.Format()
	if err != nil {
		log.Fatalf("Inspection failed: %v", err)
	}
	fmt.Printf("%s\n", data)
}

func editMetadata(filename string, partial *meta.Partial, comment bool) {
	metadata := open(filename)
	defer metadata.Close()

	updated, err := metadata.SolicitUpdates(partial, comment)
	if err != nil {
		log.Fatalf("Updates failed: %v", err)
	}
	updated.Write()
	log.Printf("Updated metadata.")
}

func getDiscogsMeta(releaseID string) *meta.Partial {
	id, err := strconv.Atoi(regexp.MustCompile(`\d+`).FindString(releaseID))
	if err != nil {
		log.Printf("No valid release ID in '%v'", releaseID)
		return &meta.Partial{}
	}
	client, err := discogs.New(&discogs.Options{
		UserAgent: "id3ed +github.com/lukasschwab/id3ed",
	})
	if err != nil {
		log.Printf("Error constructing Discogs client: %v", err)
		return &meta.Partial{}
	}
	release, err := client.Release(id)
	if err != nil {
		log.Printf("Couldn't get Discogs release '%v': %v", id, err)
		return &meta.Partial{}
	}
	year := strconv.Itoa(release.Year)
	genre := strings.Join(release.Genres, ", ")
	return &meta.Partial{
		Artist: &release.ArtistsSort,
		Album:  &release.Title,
		Year:   &year,
		Genre:  &genre,
	}
}

func main() {
	// TODO: take settable fields as named arguments.
	inspect := flag.Bool("inspect", false, "print current file metadata")
	// Editor options.
	discogsID := flag.String("discogs", "", "Discogs.com release (ID or URL) to pre-fill")
	partial := &meta.Partial{
		Title:  flag.String("title", "", "title to pre-fill"),
		Artist: flag.String("artist", "", "artist to pre-fill"),
		Album:  flag.String("album", "", "album to pre-fill"),
		Year:   flag.String("year", "", "year to pre-fill"),
		Genre:  flag.String("genre", "", "genre to pre-fill"),
	}
	comment := flag.Bool("comment", true, "include filename in a JWCC comment")

	flag.Parse()
	filename := flag.Args()[0]

	if *inspect {
		inspectMetadata(filename)
		return
	}

	if discogsID != nil {
		discogsMeta := getDiscogsMeta(*discogsID)
		partial.Mask(discogsMeta)
	}

	editMetadata(filename, partial, *comment)
}
