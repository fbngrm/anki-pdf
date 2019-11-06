package anki

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
)

type Deck struct {
	Name       string `json:"name"`
	NoteModels []struct {
		ID     string `json:"crowdanki_uuid"`
		Fields []struct {
			Name string `json:"name"`
		} `json:"flds"`
	} `json:"note_models"`
	Notes []struct {
		NoteModelID string   `json:"note_model_uuid"`
		Fields      []string `json:"fields"`
	} `json:"notes"`
}

// New loads an anki deck from file.
func New(path string) (*Deck, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return load(f)
}

// load loads an anki deck from an io.Reader.
func load(in io.Reader) (*Deck, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(in)
	if err != nil {
		return nil, err
	}

	var d Deck
	if err := json.Unmarshal(buf.Bytes(), &d); err != nil {
		return nil, err
	}
	return &d, nil
}
