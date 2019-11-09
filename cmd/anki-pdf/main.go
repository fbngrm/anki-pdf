package main

import (
	"fmt"
	"html"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/fgrimme/anki-pdf/anki"
	"github.com/fgrimme/anki-pdf/config"
	"github.com/fgrimme/anki-pdf/document"
	"github.com/fgrimme/anki-pdf/layout"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/jung-kurt/gofpdf"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	version   = "unkown"
	cfgpath   = kingpin.Flag("cfg-path", "path to config file").Short('c').Required().String()
	ankipath  = kingpin.Flag("anki-path", "path to anki deck JSON file").Short('a').Required().String()
	fontpath  = kingpin.Flag("font-path", "path to directory containing font files").Short('f').Required().String()
	mediapath = kingpin.Flag("media-path", "path to directory containing media files").Short('m').String()
)

func main() {
	kingpin.Version(version)
	kingpin.Parse()

	// deck specific configuration
	c, err := config.FromFile(*cfgpath)
	if err != nil {
		panic(err)
	}

	// load the anki deck from file
	deck, err := anki.New(*ankipath)
	if err != nil {
		panic(err)
	}

	// create cards from the anki deck
	cards, err := document.Cards(c, deck)
	if err != nil {
		panic(err)
	}

	// calc sizes and orientation
	l := layout.New(c.CardSize)

	// create a document with cards ordered in pages and rows
	d := document.New(l.PageSize, l.CardSize, cards)

	render(c, l, d)
}

func render(c *config.Config, l *layout.Layout, d document.Document) {
	orientations := map[layout.Orientation]string{
		layout.Landscape: gofpdf.OrientationLandscape,
		layout.Portrait:  gofpdf.OrientationPortrait,
	}
	pdf := gofpdf.New(orientations[l.O], "mm", "A4", "./font")
	pdf.SetMargins(0, 0, 0)
	pdf.SetAutoPageBreak(false, 0)

	// remove duplicate whitespaces
	space := regexp.MustCompile(`\s+`)

	// position on the page
	var x, y float64

	w := l.CardSize.W
	h := l.CardSize.H
	margin := c.Margin
	paddingTop := 1.

	// error report
	hText := 0.0
	hCard := h - 2*margin
	errs := make(map[string][]string, 0)

	// render loops
	for _, page := range d {
		// front page
		pdf.AddPage()
		y = 0
		for _, row := range page {
			x = 0
			for _, card := range row {
				pdf.SetDrawColor(220, 220, 220)
				pdf.Rect(x, y, w, h, "D")
				pdf.SetXY(x+margin, y+margin+paddingTop)
				for _, field := range c.Front.Fields {
					// do not render fields if we are already out of bounds when trimming
					if hText > hCard && c.ErrorStrat == "trim" {
						continue
					}
					// default layout
					font := c.Front.Layout.Font
					checkFontPath(font) // panics if font is not accessible
					size := c.Front.Layout.Size
					height := c.Front.Layout.Height
					align := c.Front.Layout.Align
					color := c.Front.Layout.Color
					// optional field fromatting from config
					if c.FieldLayouts[field].Size > 0 {
						size = c.FieldLayouts[field].Size
					}
					if c.FieldLayouts[field].Height > 0 {
						height = c.FieldLayouts[field].Height
					}
					if len(c.FieldLayouts[field].Font) > 0 {
						font = c.FieldLayouts[field].Font
						checkFontPath(font) // panics if font is not accessible
					}
					if len(c.FieldLayouts[field].Align) > 0 {
						align = c.FieldLayouts[field].Align
					}
					if len(c.FieldLayouts[field].Color) > 0 {
						color = c.FieldLayouts[field].Color
					}
					// render image field; supports landscape orientation only
					if c.FieldLayouts[field].Image {
						img := card.Front[field]
						img = img[len("<img src=\"") : len(img)-len("\" />")] // needs fix
						path := filepath.Join(*mediapath, img)
						pdf.Image(path, x+margin, y, w-2*margin, 0, true, "", 0, "")
						// calc height of image to add it to the text height used
						// in overflow calculation
						info := pdf.RegisterImage(path, "")
						ratio := info.Width() / info.Height()
						ih := (w - 2*margin) / ratio
						hText += ih // increase text height for line-breaks
						pdf.SetXY(x+margin, pdf.GetY()+margin)
						continue
					}
					// set formatting
					pdf.AddUTF8Font(font, "", font+".ttf")
					pdf.SetFont(font, "", size)
					pdf.SetTextColor(color[0], color[1], color[2])
					// render
					if field == "break" {
						pdf.MultiCell(w-2*margin, height, "", "0", align, false)
						pdf.SetXY(x+margin, pdf.GetY())
						hText += height // increase text height for line-breaks
						continue
					}
					txt := card.Front[field]
					if c.StripHTML {
						txt = strip.StripTags(txt)
					}
					if c.TrimSpace {
						txt = space.ReplaceAllString(txt, " ")
					}
					txt = html.UnescapeString(txt)
					// check height
					var lines []string
					if c.UTF8 {
						lines = pdf.SplitText(txt, w-2*margin)
						hText += float64(len(lines)) * height
					} else {
						lns := pdf.SplitLines([]byte(txt), w-2*margin)
						hText += float64(len(lines)) * height
						for _, l := range lns {
							lines = append(lines, string(l))
						}
					}
					if len(lines) == 0 {
						continue
					}
					// check for error/text reaches bottom of card
					if hText > hCard {
						// do not render the field
						if c.ErrorStrat == "skip" {
							continue
						}
						// trim the field
						if c.ErrorStrat == "trim" {
							// line-height units overflowing card boundary
							outOfBounds := hText - hCard
							// lines overflowing card boundary
							linesOOB := math.Ceil(outOfBounds / height)
							// trim cells to render
							lines = lines[:len(lines)-int(linesOOB)]
						}
					}
					// render
					for _, line := range lines {
						pdf.CellFormat(w-2*margin, height, line, "", 0, align, false, 0, "")
						pdf.SetXY(x+margin, pdf.GetY()+height)
					}
					pdf.SetXY(x+margin, pdf.GetY())
				}
				// check height
				if hText > hCard {
					errs["front"] = append(errs["front"], card.ID)
				}
				hText = 0.0
				x += w
			}
			y += h
		}
		// back page
		pdf.AddPage()
		y = 0
		// render
		for _, row := range page {
			// draw from right to left
			x = l.PageSize.W - w
			// iterate cards in row from right to left
			for _, card := range row {
				pdf.SetDrawColor(220, 220, 220)
				pdf.Rect(x, y, w, h, "D")
				pdf.SetXY(x+margin, y+margin+paddingTop)
				for _, field := range c.Back.Fields {
					// do not render fields if we are already out of bounds when trimming
					if hText > hCard && c.ErrorStrat == "trim" {
						continue
					}
					// default layout
					font := c.Back.Layout.Font
					checkFontPath(font) // panics if font is not accessible
					size := c.Back.Layout.Size
					height := c.Back.Layout.Height
					align := c.Back.Layout.Align
					color := c.Back.Layout.Color
					// optional field fromatting from config
					if c.FieldLayouts[field].Size > 0 {
						size = c.FieldLayouts[field].Size
					}
					if c.FieldLayouts[field].Height > 0 {
						height = c.FieldLayouts[field].Height
					}
					if len(c.FieldLayouts[field].Font) > 0 {
						font = c.FieldLayouts[field].Font
						checkFontPath(font) // panics if font is not accessible
					}
					if len(c.FieldLayouts[field].Align) > 0 {
						align = c.FieldLayouts[field].Align
					}
					if len(c.FieldLayouts[field].Color) > 0 {
						color = c.FieldLayouts[field].Color
					}
					// render image field; supports landscape orientation only
					if c.FieldLayouts[field].Image {
						img := card.Back[field]
						img = img[len("<img src=\"") : len(img)-len("\" />")]
						path := filepath.Join(*mediapath, img)
						pdf.Image(path, x+margin, y, w-2*margin, 0, true, "", 0, "")
						// calc height of image to add it to the text height used
						// in overflow calculation
						info := pdf.RegisterImage(path, "")
						ratio := info.Width() / info.Height()
						ih := (w - 2*margin) / ratio
						hText += ih // increase text height for line-breaks
						pdf.SetXY(x+margin, pdf.GetY()+margin)
						continue
					}
					// set formatting
					pdf.AddUTF8Font(font, "", font+".ttf")
					pdf.SetFont(font, "", size)
					pdf.SetTextColor(color[0], color[1], color[2])
					// render
					if field == "break" {
						pdf.MultiCell(w-2*margin, height, "", "0", align, false)
						pdf.SetXY(x+margin, pdf.GetY())
						hText += height // increase text height for line-breaks
						continue
					}
					txt := card.Back[field]
					if c.StripHTML {
						txt = strip.StripTags(txt)
					}
					if c.TrimSpace {
						txt = space.ReplaceAllString(txt, " ")
					}
					txt = html.UnescapeString(txt)
					// check height
					var lines []string
					if c.UTF8 {
						lines = pdf.SplitText(txt, w-2*margin)
						hText += float64(len(lines)) * height
					} else {
						lns := pdf.SplitLines([]byte(txt), w-2*margin)
						hText += float64(len(lines)) * height
						for _, l := range lns {
							lines = append(lines, string(l))
						}
					}
					if len(lines) == 0 {
						continue
					}
					// check for error/text reaches bottom of card
					if hText > hCard {
						// do not render the field
						if c.ErrorStrat == "skip" {
							continue
						}
						// trim the field
						if c.ErrorStrat == "trim" {
							// line-height units overflowing card boundary
							outOfBounds := hText - hCard
							// lines overflowing card boundary
							linesOOB := math.Ceil(outOfBounds / height)
							// trim cells to render
							lines = lines[:len(lines)-int(linesOOB)]
						}
					}
					// render
					for _, line := range lines {
						pdf.CellFormat(w-2*margin, height, line, "", 0, align, false, 0, "")
						pdf.SetXY(x+margin, pdf.GetY()+height)
					}
					pdf.SetXY(x+margin, pdf.GetY())
				}
				// error reporting
				if hText > hCard {
					errs["back"] = append(errs["back"], card.ID)
				}
				hText = 0.0
				x -= w
			}
			y += h
		}
	}

	outpath := *ankipath
	outpath = outpath[0 : len(outpath)-len(filepath.Ext(outpath))]
	err := pdf.OutputFileAndClose(outpath + ".pdf")
	if err != nil {
		panic(err)
	}
	// error report
	fmt.Printf("%s fields front %d: %v\n", c.ErrorStrat, len(errs["front"]), errs["front"])
	fmt.Printf("%s fields back %d: %v\n", c.ErrorStrat, len(errs["back"]), errs["back"])

}

// tests if a file is accessible
func accessable(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// checks if a font is accessible
func checkFontPath(font string) {
	p, err := filepath.Abs(*fontpath)
	if err != nil {
		panic(err)
	}
	path := path.Join(p, font+".ttf")
	if !accessable(path) {
		log.Panicf("cannot find font: %s", path)
	}
}
