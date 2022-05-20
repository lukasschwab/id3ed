package editor

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

// GetUpdates from the user. Open initial in an editor as a temporary file, and
// return whatever the file contains when the editor is closed.
func GetUpdates(initial []byte) ([]byte, error) {
	f, err := os.CreateTemp("", "meta.*.json")
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
