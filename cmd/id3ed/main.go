package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/lukasschwab/id3ed/pkg/meta"
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

func main() {
	// TODO: take settable fields as named arguments.
	inspect := flag.Bool("inspect", false, "print current file metadata")
	// Editor options.
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

	editMetadata(filename, partial, *comment)
}
