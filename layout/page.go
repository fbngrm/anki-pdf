package layout

type Rect struct {
	W float64
	H float64
}

type Layout struct {
	O        Orientation
	PageSize Rect
	CardSize Rect // not DIN, we want fill the whole page
}

func New(size string) *Layout {
	s := cardSizes[size]
	return &Layout{
		O:        s.o,
		PageSize: pageSize[s.o],
		CardSize: s.size,
	}
}

type Orientation int

const (
	Landscape Orientation = iota
	Portrait
)

// orientation depends on card size
var pageSize = map[Orientation]Rect{
	Landscape: {
		W: 297.,
		H: 210.,
	},
	Portrait: {
		W: 210.,
		H: 297.,
	},
}

// always landscape
var cardSizes = map[string]struct {
	o    Orientation
	size Rect
}{
	"A4": {
		o: Landscape,
		size: Rect{
			W: 297.,
			H: 210.,
		},
	},
	"A5": {
		o: Portrait,
		size: Rect{
			W: 210.,
			H: 148.5,
		},
	},
	"A6": {
		o: Landscape,
		size: Rect{
			W: 148.5,
			H: 105.,
		},
	},
	"A7": {
		o: Portrait,
		size: Rect{
			W: 105.,
			H: 74.25,
		},
	},
	"A8": {
		o: Landscape,
		size: Rect{
			W: 74.25,
			H: 52.5,
		},
	},
}
