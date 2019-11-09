## anki-pdf

anki-pdf converts anki-decks to PDF files.

## Features

- Chose fields from your anki-cards and format them individually
- Front and back sides of the cards can be configured and formatted
- Supports images embedding
- Different strategies for overflowing text

## Installation

```bash
go get https://github.com/fbngrm/anki-pdf
```

Binaries will be located in the `anki-pdf/bin` directory. Builds use the latest commit hash of the master branch or tag.

```bash
cd anki-pdf
make build
./bin/anki-pdf --version
```

Install the program in your $GOPATH.

```bash
cd anki-pdf
make install
anki-pdf --version
```

## Generate a PDF
Two input files are required to generate a PDF file.

1. A JSON representation of your anki-deck, generated with [CrowdAnki](https://ankiweb.net/shared/info/1788670778)
2. A YAML configuration file which tells teh program how to render the PDF.

### CrowdAnki
Follow the export [instructions](https://github.com/Stvad/CrowdAnki#export) to create the JSON file.

### Configuration
The `example/` directory contains configuration files which should be used as a starting point.

```bash
anki-pdf -c path/to/anki-deck.yaml -a path/to/anki-deck.json -f path/to/fonts [-m path/to/media]
```

## Fonts
To add new fonts, follow the instructions of [gofpdf](https://github.com/jung-kurt/gofpdf#nonstandard-fonts).
If you use UFT-8 encoded text, set `utf8: true` in the configuration file of the anki-deck.

## Error reports
After the program terminated, an error report will be printed.
To optimize the formatting and reduce overflow errors, adjust the configuration for the overflowing fields or chose a different error handling strategy.
