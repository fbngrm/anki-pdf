package document

import (
	"github.com/fgrimme/anki-pdf/anki"
	"github.com/fgrimme/anki-pdf/config"
)

// Card holds all fields for front and back side of a card by model name.
type Card struct {
	Front map[string]string
	Back  map[string]string
}

func Cards(c *config.Config, d *anki.Deck) ([]Card, error) {
	// notes is a map of all notes in the deck, with their fields mapped by
	// note-model name
	notes, err := makeNotes(d)
	if err != nil {
		return nil, err
	}
	// sort by front and back page and substitue empty fields
	return makeCards(c, notes), nil
}

// notes contain fields by note-model name
type notes []map[string]string

// makeNotes is an intermdiate processing step to map fields by their note-model
// name. Thus, we can easily access fields and get their formatting configuration
// by note-model name in the rendering step.
func makeNotes(d *anki.Deck) (notes, error) {
	// index of note-model names by id
	models := make(map[string][]string)
	for _, m := range d.NoteModels {
		names := make([]string, len(m.Fields))
		for i, f := range m.Fields {
			names[i] = f.Name
		}
		models[m.ID] = names
	}
	// map fields by model name
	nts := make(notes, len(d.Notes))
	for i, n := range d.Notes {
		fields := make(map[string]string)
		for y, f := range n.Fields {
			modelName := models[n.NoteModelID][y]
			fields[modelName] = f
		}
		nts[i] = fields
	}
	return nts, nil
}

// makeCards creates a card with front and back pages and substitutes empty fields.
func makeCards(c *config.Config, n notes) []Card {
	cards := make([]Card, len(n))
	for i, note := range n {
		// get fields for front page
		front := make(map[string]string)
		for _, field := range c.Front.Fields {
			n := note[field]
			if len(n) == 0 { // substitute empty field
				n = note[c.Empty[field]]
			}
			front[field] = n
		}
		// get fields for back page
		back := make(map[string]string)
		for _, field := range c.Back.Fields {
			n := note[field]
			if len(n) == 0 { // substitute empty field
				n = note[c.Empty[field]]
			}
			back[field] = n
		}
		cards[i] = Card{
			Front: front,
			Back:  back,
		}
	}
	return cards
}
