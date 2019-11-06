package main

import (
	"github.com/fgrimme/ankiPDF/anki"
	"github.com/fgrimme/ankiPDF/config"
	"github.com/fgrimme/ankiPDF/document"
	"github.com/fgrimme/ankiPDF/layout"
)

func main() {
	cfg, err := config.FromFile("./cfg.yaml")
	if err != nil {
		panic(err)
	}

	// we load the anki deck from file
	deck, err := anki.New("./01_NihongoShark.com__Kanji/01_NihongoShark.json")
	if err != nil {
		panic(err)
	}

	// cards have configured fields for front and back
	cards, err := anki.Cards(cfg, deck)
	if err != nil {
		panic(err)
	}

	// layout
	l := layout.New(cfg.CardSize)

	// document
	doc := document.New(l.PageSize, l.CardSize, cards)
	for _, page := range doc {
		for _, row := range page {
			for _, cell := range row {
				// fmt.Printf("%+v\n", cell)
				_ = cell
			}
		}
	}

	// pdf := gofpdf.New(gofpdf.OrientationLandscape, "mm", "A4", "./font")
	// pdf.AddUTF8Font("NotoSansSC-Regular", "", "NotoSansSC-Regular.ttf")
	// pdf.SetFont("NotoSansSC-Regular", "", 16)
	// pdf.SetMargins(0, 0, 0)
	// pdf.SetAutoPageBreak(false, 0)
	// // pdf.SetFillColor(00, 00, 200)

	// s := "中文，你好！"

	// margin := 5.
	// lHeight := 5.5

	// var x, y float64
	// h := float64(cf.h)
	// w := float64(cf.w)

	// for _, page := range doc {
	// 	pdf.AddPage()
	// 	y = 0
	// 	for _, row := range page {
	// 		x = 0
	// 		for _, cell := range row {
	// 			pdf.Rect(x, y, h, w, "D")
	// 			pdf.SetXY(x+margin, y+margin)
	// 			pdf.MultiCell(h-2*margin, lHeight, string(cell)+s, "0", "CM", false)
	// 			x += h
	// 		}
	// 		y += w
	// 	}
	// }

	// err := pdf.OutputFileAndClose("Fpdf_AddPage.pdf")
	// if err != nil {
	// 	panic(err)
	// }

}
