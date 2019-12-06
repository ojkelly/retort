package component

import (
	"math"

	"github.com/gdamore/tcell"

	"retort.dev/r"
)

// Box is the basic building block for a retort app.
// Box implements the Box Model, see BoxProps
func Box(p r.Properties) r.Element {
	screen := r.UseScreen()

	// Get our BoxProps
	boxProps := p.GetProperty(
		BoxProps{},
		"Box requires BoxProps",
	).(BoxProps)

	// Get our BoxLayout
	parentBoxLayout := p.GetProperty(
		r.BoxLayout{},
		"Box requires a parent BoxLayout.",
	).(r.BoxLayout)

	// Get any children
	children := p.GetOptionalProperty(
		r.Children{},
	).(r.Children)

	// Calculate the BoxLayout of this Box
	boxLayout, innerBoxLayout := calculateBoxLayout(
		screen,
		parentBoxLayout,
		boxProps,
	)

	// Calculate the BoxLayout of any children
	childrenWithLayout := calculateBoxLayoutForChildren(
		screen,
		boxProps,
		innerBoxLayout,
		children,
	)

	return r.CreateScreenElement(
		func(s tcell.Screen) r.BoxLayout {
			if s == nil {
				panic("no screen in context")
			}

			w, h := s.Size()

			if w == 0 || h == 0 {
				panic("Box can't render on a zero size screen")
			}

			st := tcell.StyleDefault
			gl := ' '

			st = st.Foreground(boxProps.Border.Foreground)

			drawBox(
				screen,
				boxLayout.X,
				boxLayout.Y,
				boxLayout.X+boxLayout.Columns,
				boxLayout.Y+boxLayout.Rows,
				st,
				gl,
			)

			return boxLayout
		},
		childrenWithLayout,
	)
}

func drawBox(
	s tcell.Screen,
	x1, y1, x2, y2 int,
	style tcell.Style,
	r rune,
) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, style)
		s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}
	if y1 != y2 && x1 != x2 {
		// Only add corners if we need to
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		for col := x1 + 1; col < x2; col++ {
			s.SetContent(col, row, r, nil, style)
		}
	}
}

// [ BoxModel Types ]-----------------------------------------------------------

type Padding struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

type Margin struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

type Border struct {
	Style      BorderStyle
	Background tcell.Color
	Foreground tcell.Color
}

type BorderStyle int

const (
	BorderStyleNone BorderStyle = iota
	BorderStyleSingle
	BorderStyleDouble
	BorderStyleBox // Box drawing characters
)

// [ FlexBox Types ]------------------------------------------------------------

type FlexDirection int

const (
	FlexDirectionRow FlexDirection = iota
	FlexDirectionRowReverse
	FlexDirectionColumn
	FlexDirectionColumnReverse
)

type FlexBasis int

const (
	FlexBasisFill FlexBasis = iota
	FlexBasisMaxContent
	FlexBasisMinContent
	FlexBasisFitContent
)

type FlexWrapOption int

const (
	FlexWrapNone FlexWrapOption = iota
	FlexWrap
	FlexWrapReverse
)

// [ FlexBox Funcs ]------------------------------------------------------------

type BoxProps struct {
	ZIndex int

	// Flex Box

	// The flex-direction CSS property sets how flex items are placed in the flex
	// container defining the main axis and the direction (normal or reversed).
	FlexDirection FlexDirection

	// The flex-basis CSS property sets the initial main size of a flex item.
	FlexBasis FlexBasis

	// The flex-grow CSS property sets the flex grow factor of a flex item main
	// size. It specifies how much of the remaining space in the flex container
	// should be assigned to the item (the flex grow factor).
	FlexGrow int

	// The flex-shrink CSS property sets the flex shrink factor of a flex item.
	// If the size of all flex items is larger than the flex container, items
	// shrink to fit according to flex-shrink.
	FlexShrink int // TODO

	// The flex-wrap CSS property sets whether flex items are forced onto one
	// line or can wrap onto multiple lines. If wrapping is allowed, it sets the
	// direction that lines are stacked.
	FlexWrap FlexWrapOption

	// Content Box
	// If neither Width,Height or Rows,Columns are set, it will be calculated
	// automatically  When set this is the percentage width and height.
	// Ignored when Rows,Columns is not 0
	Width, Height float64 // 0 = auto

	// Set the size fixed in rows and columns.
	// Ignored if 0
	// If both Rows and Width are set Rows with be used.
	Rows, Columns int

	// Padding Box
	Padding Padding
	Margin  Margin

	// Box Sizing is Border Box only
	// Border and padding is accounted for inside the width and height, meaning
	// the Box can never be bigger than the width or height.

	// Border
	Border Border

	Background tcell.Color
	Foreground tcell.Color
}

func calculateBoxLayout(
	screen tcell.Screen,
	parentBoxLayout r.BoxLayout,
	boxProps BoxProps,
) (
	boxLayout r.BoxLayout,
	innerBoxLayout r.BoxLayout,
) {
	rows := parentBoxLayout.Rows
	columns := parentBoxLayout.Columns

	if rows == 0 && boxProps.Height != 0 {
		rows = int(
			math.Round(
				float64(parentBoxLayout.Rows) * (boxProps.Height / 100),
			),
		)
	}
	if columns == 0 && boxProps.Width != 0 {
		columns = int(
			math.Round(
				float64(parentBoxLayout.Columns) * (boxProps.Width / 100),
			),
		)
	}

	boxLayout = r.BoxLayout{
		ZIndex:  boxProps.ZIndex,
		Rows:    rows,
		Columns: columns,
		X:       parentBoxLayout.X,
		Y:       parentBoxLayout.Y,
	}
	innerBoxLayout = r.BoxLayout{
		ZIndex:  boxProps.ZIndex,
		Rows:    rows,
		Columns: columns,
		X:       parentBoxLayout.X,
		Y:       parentBoxLayout.Y,
	}

	// Calculate box size
	boxLayout.Columns = columns - boxProps.Padding.Left - boxProps.Padding.Right
	boxLayout.Rows = rows - boxProps.Padding.Top - boxProps.Padding.Bottom

	innerBoxLayout.Columns = boxLayout.Columns -
		boxProps.Padding.Left - boxProps.Padding.Right

	innerBoxLayout.Rows = boxLayout.Rows -
		boxProps.Padding.Top - boxProps.Padding.Bottom

	boxLayout.X = parentBoxLayout.X + boxProps.Margin.Left - boxProps.Margin.Right
	boxLayout.Y = parentBoxLayout.Y + boxProps.Margin.Top - boxProps.Margin.Bottom

	// Calculate padding box

	if boxProps.Padding.Top != 0 {
		innerBoxLayout.Y = innerBoxLayout.Y + boxProps.Padding.Top
		innerBoxLayout.Rows = innerBoxLayout.Rows + boxProps.Padding.Top
	}

	if boxProps.Padding.Right != 0 {
		innerBoxLayout.X = innerBoxLayout.Y - boxProps.Padding.Right
		innerBoxLayout.Columns = innerBoxLayout.Columns - boxProps.Padding.Right
	}

	if boxProps.Padding.Bottom != 0 {
		innerBoxLayout.Y = innerBoxLayout.Y - boxProps.Padding.Bottom
		innerBoxLayout.Rows = innerBoxLayout.Rows - boxProps.Padding.Bottom
	}

	if boxProps.Padding.Left != 0 {
		innerBoxLayout.Y = innerBoxLayout.Y + boxProps.Padding.Left
		innerBoxLayout.Columns = innerBoxLayout.Columns + boxProps.Padding.Left
	}

	// Border Sizing

	if boxProps.Border.Style != BorderStyleNone {
		boxLayout.Columns = boxLayout.Columns - 2 // 1 for each side
		boxLayout.Rows = boxLayout.Rows - 2       // 1 for each side

		innerBoxLayout.X = innerBoxLayout.X + 1
		innerBoxLayout.Y = innerBoxLayout.Y + 1
		innerBoxLayout.Rows = innerBoxLayout.Rows - 2
		innerBoxLayout.Columns = innerBoxLayout.Columns - 2
	}

	// Ensure the rows and cols are not below 0
	if boxLayout.Rows < 0 {
		boxLayout.Rows = 0
	}
	if boxLayout.Columns < 0 {
		boxLayout.Columns = 0
	}
	if innerBoxLayout.Rows < 0 {
		innerBoxLayout.Rows = 0
	}
	if innerBoxLayout.Columns < 0 {
		innerBoxLayout.Columns = 0
	}
	return
}

func calculateBoxLayoutForChildren(
	screen tcell.Screen,
	boxProps BoxProps,
	innerBoxLayout r.BoxLayout,
	children r.Children,
) r.Children {

	if len(children) == 0 {
		return children
	}

	propMap := map[r.Element]BoxProps{}

	colsRemaining := innerBoxLayout.Columns
	rowsRemaining := innerBoxLayout.Rows
	flexGrowCount := 0
	flexGrowDivision := 0

	for _, c := range children {
		if c == nil {
			continue
		}
		propMap[c] = c.Properties.GetOptionalProperty(
			BoxProps{},
		).(BoxProps)
	}

	// Find all children with fixed row,col sizing
	for _, props := range propMap {
		colsRemaining = colsRemaining - props.Columns
		rowsRemaining = rowsRemaining - props.Rows
		flexGrowCount = flexGrowCount + props.FlexGrow
		if props.FlexGrow == 0 {
			flexGrowCount = flexGrowCount + 1 // we force flex-grow to be at least 1
		}
	}

	switch boxProps.FlexDirection {
	case FlexDirectionRow:
		flexGrowDivision = colsRemaining / flexGrowCount
	case FlexDirectionRowReverse:
		flexGrowDivision = colsRemaining / flexGrowCount
	case FlexDirectionColumn:
		flexGrowDivision = rowsRemaining / flexGrowCount
	case FlexDirectionColumnReverse:
		flexGrowDivision = rowsRemaining / flexGrowCount

	}

	if boxProps.FlexDirection == FlexDirectionRowReverse ||
		boxProps.FlexDirection == FlexDirectionColumnReverse {
		for i := len(children)/2 - 1; i >= 0; i-- {
			opp := len(children) - 1 - i
			children[i], children[opp] = children[opp], children[i]
		}
	}

	x := innerBoxLayout.X
	y := innerBoxLayout.Y

	for i, el := range children {
		if el == nil {
			continue
		}

		props := propMap[el]

		row := 0
		c := 0
		z := boxProps.ZIndex

		switch boxProps.FlexDirection {
		case FlexDirectionRow:
			c = flexGrowDivision * props.FlexGrow
			row = innerBoxLayout.Rows
		case FlexDirectionRowReverse:
			c = flexGrowDivision * props.FlexGrow
			row = innerBoxLayout.Rows
		case FlexDirectionColumn:
			c = innerBoxLayout.Columns
			row = flexGrowDivision * props.FlexGrow
		case FlexDirectionColumnReverse:
			c = innerBoxLayout.Columns
			row = flexGrowDivision * props.FlexGrow
		}

		// Ensure r and c aren't negative
		if row < 0 {
			row = 0
		}
		if c < 0 {
			c = 0
		}

		boxLayout := r.BoxLayout{
			X:       x,
			Y:       y,
			Rows:    row,
			Columns: c,
			ZIndex:  z,
			Order:   i,
		}

		switch boxProps.FlexDirection {
		case FlexDirectionRow:
			x = x + c
		case FlexDirectionRowReverse:
			x = x + c
		case FlexDirectionColumn:
			y = y + row
		case FlexDirectionColumnReverse:
			y = y + row
		}
		el.Properties = r.ReplaceProps(el.Properties, boxLayout)
	}
	return children
}
