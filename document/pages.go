package document

import (
	"github.com/fgrimme/anki-pdf/layout"
)

type row []Card
type page []row
type Document []page

// New orders the given cards in pages and rows so that we can easily access
// them in the render step by iterating the pages. Page and card size is used
// to determine the count of cards in a row and rows on a page.
func New(p, c layout.Rect, cards []Card) Document {
	cellsPerRow := int(p.W / c.W)
	rowsPerPage := int(p.H / c.H)

	// order all cards in rows
	var rows []row
	for cellsPerRow < len(cards) {
		cards, rows = cards[cellsPerRow:], append(rows, cards[0:cellsPerRow:cellsPerRow])
	}
	rows = append(rows, cards)

	// order all rows in pages
	var pages []page
	for rowsPerPage < len(rows) {
		rows, pages = rows[rowsPerPage:], append(pages, rows[0:rowsPerPage:rowsPerPage])
	}
	pages = append(pages, rows)
	return pages
}
