package document

import (
	"github.com/fgrimme/ankiPDF/anki"
	"github.com/fgrimme/ankiPDF/layout"
)

type row []anki.Card
type page []row
type document []page

func New(p, c layout.Rect, cards []anki.Card) document {
	cellsPerRow := int(p.W / c.W)
	rowsPerPage := int(p.H / c.H)

	var rows []row
	for cellsPerRow < len(cards) {
		cards, rows = cards[cellsPerRow:], append(rows, cards[0:cellsPerRow:cellsPerRow])
	}
	rows = append(rows, cards)

	var pages []page
	for rowsPerPage < len(rows) {
		rows, pages = rows[rowsPerPage:], append(pages, rows[0:rowsPerPage:rowsPerPage])
	}
	pages = append(pages, rows)
	return pages
}
