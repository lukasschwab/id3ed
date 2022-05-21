package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/lukasschwab/id3ed/pkg/meta"
	"github.com/mikkyang/id3-go"
)

func inspectMetadata(file *id3.File) {
	data, err := meta.From(file).Format()
	if err != nil {
		log.Fatalf("Inspection failed: %v", err)
	}
	fmt.Printf("%s\n", data)
}

func editMetadata(file *id3.File) {
	data := meta.From(file)
	updated, err := data.SolicitUpdates()
	if err != nil {
		log.Fatalf("Updates failed: %v", err)
	}
	updated.Write(file)
	log.Printf("Updated metadata.")
}

func main() {
	inspect := flag.Bool("inspect", false, "print current file metadata")
	flag.Parse()

	// TODO: take settable fields as named arguments.
	filename := flag.Args()[0]
	file, err := id3.Open(filename)
	if err != nil {
		log.Fatalf("Parsing failed: %v", err)
	}
	defer file.Close()

	if *inspect {
		inspectMetadata(file)
		return
	}

	editMetadata(file)
}
