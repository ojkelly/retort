package box

import (
	"math"

	"github.com/gdamore/tcell"
	"retort.dev/r"
)

func calculateBoxLayout(
	screen tcell.Screen,
	parentBoxLayout r.BoxLayout,
	boxProps Properties,
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

	// Calculate box size
	boxLayout.Columns = columns
	boxLayout.Rows = rows
	// Calculate margin

	boxLayout.X = parentBoxLayout.X + boxProps.Margin.Left
	boxLayout.Columns = boxLayout.Columns - boxProps.Margin.Right
	boxLayout.Y = parentBoxLayout.Y + boxProps.Margin.Top
	boxLayout.Rows = boxLayout.Rows - boxProps.Margin.Bottom

	innerBoxLayout = r.BoxLayout{
		ZIndex:  boxProps.ZIndex,
		Rows:    rows,
		Columns: columns,
		X:       boxLayout.X,
		Y:       boxLayout.Y,
	}

	innerBoxLayout.Columns = boxLayout.Columns -
		boxProps.Padding.Left - boxProps.Padding.Right

	innerBoxLayout.Rows = boxLayout.Rows -
		boxProps.Padding.Top - boxProps.Padding.Bottom

	// Calculate padding box

	if boxProps.Padding.Top != 0 {
		innerBoxLayout.Y = innerBoxLayout.Y + boxProps.Padding.Top
	}

	if boxProps.Padding.Right != 0 {
		innerBoxLayout.Columns = innerBoxLayout.Columns - boxProps.Padding.Right
	}

	if boxProps.Padding.Bottom != 0 {
		innerBoxLayout.Rows = innerBoxLayout.Rows - boxProps.Padding.Bottom
	}

	if boxProps.Padding.Left != 0 {
		innerBoxLayout.X = innerBoxLayout.X + boxProps.Padding.Left
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
	boxProps Properties,
	innerBoxLayout r.BoxLayout,
	children r.Children,
) r.Children {

	if len(children) == 0 {
		return children
	}

	propMap := map[r.Element]Properties{}

	colsRemaining := innerBoxLayout.Columns
	rowsRemaining := innerBoxLayout.Rows
	flexGrowCount := 0
	flexGrowDivision := 0

	for _, c := range children {
		if c == nil {
			continue
		}
		propMap[c] = c.Properties.GetOptionalProperty(
			Properties{},
		).(Properties)
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
