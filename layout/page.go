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

func New(cardSize DIN) *Layout {
	o := Orientation(cardSize % 2)
	return &Layout{
		O:        o,
		PageSize: pageSize[o],
		CardSize: cardSizes[cardSize],
	}
}

type DIN int

const (
	A8 DIN = iota
	A7
	A6
	A5
	A4
)

type Orientation int

const (
	landscape Orientation = iota
	portrait
)

// orientation depends on card size
var pageSize = map[Orientation]Rect{
	landscape: {
		W: 297.,
		H: 210.,
	},
	portrait: {
		W: 210.,
		H: 297.,
	},
}

// always landscape
var cardSizes = map[DIN]Rect{
	A4: {
		W: 297.,
		H: 210.,
	},
	A5: {
		W: 210.,
		H: 148.5,
	},
	A6: {
		W: 148.5,
		H: 10.5,
	},
	A7: {
		W: 10.5,
		H: 74.25,
	},
	A8: {
		W: 74.25,
		H: 5.25,
	},
}
