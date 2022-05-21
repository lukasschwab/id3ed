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

func editMetadata(filename string) {
	metadata := open(filename)
	defer metadata.Close()

	updated, err := metadata.SolicitUpdates()
	if err != nil {
		log.Fatalf("Updates failed: %v", err)
	}
	updated.Write()
	log.Printf("Updated metadata.")
}

func main() {
	// TODO: take settable fields as named arguments.
	inspect := flag.Bool("inspect", false, "print current file metadata")
	_ = flag.Bool("withFilename", true, "include the filename in a JWCC comment")

	flag.Parse()
	filename := flag.Args()[0]

	if *inspect {
		inspectMetadata(filename)
		return
	}
	editMetadata(filename)
}
