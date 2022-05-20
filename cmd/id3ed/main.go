package main

import (
	"log"
	"os"

	"github.com/lukasschwab/id3ed/pkg/meta"
	"github.com/mikkyang/id3-go"
)

func main() {
	// TODO: take settable fields as named arguments.
	filename := os.Args[1]
	file, err := id3.Open(filename)
	if err != nil {
		log.Fatalf("Parsing failed: %v", err)
	}
	defer file.Close()

	data := meta.From(file)
	updated, err := data.SolicitUpdates()
	if err != nil {
		log.Fatalf("Updates failed: %v", err)
	}

	updated.Write(file)
	log.Printf("Updated metadata.")
}
