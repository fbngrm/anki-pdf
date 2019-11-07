package main

import (
	"html"
	"path/filepath"
	"regexp"

	"github.com/fgrimme/ankiPDF/anki"
	"github.com/fgrimme/ankiPDF/config"
	"github.com/fgrimme/ankiPDF/document"
	"github.com/fgrimme/ankiPDF/layout"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/jung-kurt/gofpdf"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	version  = "unkown"
	cfgpath  = kingpin.Flag("cfg-path", "path to config file").Short('c').Required().String()
	ankipath = kingpin.Flag("anki-path", "path to anki deck JSON file").Short('a').Required().String()
)

func main() {
	kingpin.Version(version)
	kingpin.Parse()

	cfg, err := config.FromFile(*cfgpath)
	if err != nil {
		panic(err)
	}

	// we load the anki deck from file
	// name := "./anki-decks/01_NihongoShark.com__Kanji/01_NihongoShark"
	// name := "./anki-decks/Katakana_with_stroke_diagrams_and_audio/Katakana_with_stroke_diagrams_and_audio"
	deck, err := anki.New(*ankipath)
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

	// render
	orientations := map[layout.Orientation]string{
		layout.Landscape: gofpdf.OrientationLandscape,
		layout.Portrait:  gofpdf.OrientationPortrait,
	}

	pdf := gofpdf.New(orientations[l.O], "mm", "A4", "./font")
	pdf.SetMargins(0, 0, 0)
	pdf.SetAutoPageBreak(false, 0)

	var x, y float64
	w := l.CardSize.W
	h := l.CardSize.H

	margin := cfg.Margin

	// default layout for front pages
	font := cfg.Front.Layout.Font
	size := cfg.Front.Layout.Size
	height := cfg.Front.Layout.Height
	align := cfg.Front.Layout.Align
	color := cfg.Front.Layout.Color

	// remove duplicate whitespaces
	space := regexp.MustCompile(`\s+`)

	// render front pages
	for _, page := range doc {
		pdf.AddPage()
		y = 0
		for _, row := range page {
			x = 0
			for _, card := range row {
				pdf.SetDrawColor(220, 220, 220)
				pdf.Rect(x, y, w, h, "D")
				pdf.SetXY(x+margin, y+margin)
				for _, field := range cfg.Front.Fields {
					// line-break
					if field == "break" {
						pdf.Ln(cfg.FieldLayouts["break"].Height)
						pdf.SetXY(x+margin, pdf.GetY())
						continue
					}
					// optional field fromatting from config
					if cfg.FieldLayouts[field].Size > 0 {
						size = cfg.FieldLayouts[field].Size
					}
					if cfg.FieldLayouts[field].Height > 0 {
						height = cfg.FieldLayouts[field].Height
					}
					if len(cfg.FieldLayouts[field].Font) > 0 {
						font = cfg.FieldLayouts[field].Font
					}
					if len(cfg.FieldLayouts[field].Align) > 0 {
						align = cfg.FieldLayouts[field].Align
					}
					if len(cfg.FieldLayouts[field].Color) > 0 {
						color = cfg.FieldLayouts[field].Color
					}
					// set formatting
					pdf.AddUTF8Font(font, "", font+".ttf")
					pdf.SetFont(font, "", size)
					pdf.SetTextColor(color[0], color[1], color[2])
					// render
					txt := card.Front[field]
					if cfg.StripHTML {
						txt = strip.StripTags(txt)
					}
					if cfg.TrimSpace {
						txt = space.ReplaceAllString(txt, " ")
					}
					txt = html.UnescapeString(txt)
					pdf.MultiCell(w-2*margin, height, txt, "0", align, false)
					pdf.SetXY(x+margin, pdf.GetY())
				}
				x += w
			}
			y += h
		}
		pdf.AddPage()
		y = 0
		// render back pages
		for _, row := range page {
			// draw from right to left
			x = l.PageSize.W - w
			// iterate cards in row from right to left
			for _, card := range row {
				pdf.SetDrawColor(220, 220, 220)
				pdf.Rect(x, y, w, h, "D")
				pdf.SetXY(x+margin, y+margin)

				for _, field := range cfg.Back.Fields {
					// optional field fromatting from config
					if cfg.FieldLayouts[field].Size > 0 {
						size = cfg.FieldLayouts[field].Size
					}
					if cfg.FieldLayouts[field].Height > 0 {
						height = cfg.FieldLayouts[field].Height
					}
					if len(cfg.FieldLayouts[field].Font) > 0 {
						font = cfg.FieldLayouts[field].Font
					}
					if len(cfg.FieldLayouts[field].Align) > 0 {
						align = cfg.FieldLayouts[field].Align
					}
					if len(cfg.FieldLayouts[field].Color) > 0 {
						color = cfg.FieldLayouts[field].Color
					}
					// set formatting
					pdf.AddUTF8Font(font, "", font+".ttf")
					pdf.SetFont(font, "", size)
					pdf.SetTextColor(color[0], color[1], color[2])
					// render
					txt := card.Back[field]
					if cfg.StripHTML {
						txt = strip.StripTags(txt)
					}
					if cfg.TrimSpace {
						txt = space.ReplaceAllString(txt, " ")
					}
					txt = html.UnescapeString(txt)
					// render
					pdf.MultiCell(w-2*margin, height, txt, "0", align, false)
					pdf.SetXY(x+margin, pdf.GetY())
				}
				x -= w
			}
			y += h
		}
	}

	// default layout for back pages
	font = cfg.Back.Layout.Font
	size = cfg.Back.Layout.Size
	height = cfg.Back.Layout.Height
	align = cfg.Back.Layout.Align
	color = cfg.Back.Layout.Color

	outpath := *ankipath
	outpath = outpath[0 : len(outpath)-len(filepath.Ext(outpath))]
	err = pdf.OutputFileAndClose(outpath + ".pdf")
	if err != nil {
		panic(err)
	}

}
