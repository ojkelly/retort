package box

import "github.com/gdamore/tcell"

// Properties are passed along with box.Box tocreate and configure a Box element
//
// Contents
//
// The contents of the Box is not rendered by this component
//
//
// Box Sizing
//
// Box Sizing is Border Box only
// Border and padding is accounted for inside the width and height, meaning
// the Box can never be bigger than the width or height.
type Properties struct {
	// ZIndex is the layer this Box is rendered on, with larger numbers appearing
	// on top.
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

	Align Align

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

	// Border
	Border Border

	Background tcell.Color
	Foreground tcell.Color

	Overflow Overflow

	// Title is a Label placed on the top border
	Title Label

	// Footer is a Label place on the bottom border
	Footer Label
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

type Overflow int

const (
	OverflowScroll Overflow = iota
	OverflowNone
	OverflowScrollX
	OverflowScrollY
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

type Align int

const (
	AlignAuto Align = iota
	AlignStart
	AlignCenter
	AlignEnd
)

// [ Labels ]-------------------------------------------------------------------

type LabelWrap int

const (
	LabelWrapNone LabelWrap = iota
	LabelWrapBracket
	LabelWrapBrace
	LabelWrapChevron
	LabelWrapSquareBracket
)

// Label is a decorative string that can be added to the top or bottom border
//
// Margin allows you to move the whole label around, while Padding allows you
// to define the gap between the Wrap and Value.
// If no Padding is specified a single column is still added to each side of the
// Value.
type Label struct {
	Value   string
	Wrap    LabelWrap
	Align   Align
	Margin  Margin
	Padding Padding
}
