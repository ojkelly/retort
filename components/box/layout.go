package box

import (
	"math"

	"github.com/gdamore/tcell"
	"retort.dev/r"
)

func calculateBlockLayout(
	props Properties,
) r.CalculateLayout {
	return func(
		s tcell.Screen,
		stage r.CalculateLayoutStage,
		parentBlockLayout r.BlockLayout,
		childrenBlockLayouts *r.BlockLayouts,
	) (outerBlockLayout r.BlockLayout, innerBlockLayout r.BlockLayout) {
		switch stage {
		case r.CalculateLayoutStageInitial:
			// if any widths or heights are explicitly set, set them here
			// otherwise inherit from the parentBlockLayout

			rows := parentBlockLayout.Rows
			columns := parentBlockLayout.Columns

			if rows == 0 && props.Height != 0 {
				rows = int(
					math.Round(
						float64(parentBlockLayout.Rows) * (props.Height / 100),
					),
				)
			}
			if columns == 0 && props.Width != 0 {
				columns = int(
					math.Round(
						float64(parentBlockLayout.Columns) * (props.Width / 100),
					),
				)
			}

			outerBlockLayout := r.BlockLayout{
				ZIndex:  props.ZIndex,
				Rows:    rows,
				Columns: columns,
				X:       parentBlockLayout.X,
				Y:       parentBlockLayout.Y,
			}

			// Calculate box size
			outerBlockLayout.Columns = columns
			outerBlockLayout.Rows = rows

			// Calculate margin
			outerBlockLayout.X = parentBlockLayout.X + props.Margin.Left
			outerBlockLayout.Columns = outerBlockLayout.Columns - props.Margin.Right
			outerBlockLayout.Y = parentBlockLayout.Y + props.Margin.Top
			outerBlockLayout.Rows = outerBlockLayout.Rows - props.Margin.Bottom

			innerBlockLayout = r.BlockLayout{
				ZIndex:  props.ZIndex,
				Rows:    outerBlockLayout.Rows,
				Columns: outerBlockLayout.Columns,
				X:       outerBlockLayout.X,
				Y:       outerBlockLayout.Y,
			}

			innerBlockLayout.Columns = outerBlockLayout.Columns -
				props.Padding.Left - props.Padding.Right

			innerBlockLayout.Rows = outerBlockLayout.Rows -
				props.Padding.Top - props.Padding.Bottom

			// Calculate padding box
			innerBlockLayout.Y = innerBlockLayout.Y + props.Padding.Top
			innerBlockLayout.Columns = innerBlockLayout.Columns - props.Padding.Right
			innerBlockLayout.Rows = innerBlockLayout.Rows - props.Padding.Bottom
			innerBlockLayout.X = innerBlockLayout.X + props.Padding.Left

			// Border Sizing

			if props.Border.Style != BorderStyleNone {
				outerBlockLayout.Columns = outerBlockLayout.Columns - 2 // 1 for each side
				outerBlockLayout.Rows = outerBlockLayout.Rows - 2       // 1 for each side

				innerBlockLayout.X = innerBlockLayout.X + 1
				innerBlockLayout.Y = innerBlockLayout.Y + 1
				innerBlockLayout.Rows = innerBlockLayout.Rows - 1
				innerBlockLayout.Columns = innerBlockLayout.Columns - 2
			}

			// Ensure the rows and cols are not below 0
			if outerBlockLayout.Rows < 0 {
				outerBlockLayout.Rows = 0
			}
			if outerBlockLayout.Columns < 0 {
				outerBlockLayout.Columns = 0
			}
			if innerBlockLayout.Rows < 0 {
				innerBlockLayout.Rows = 0
			}
			if innerBlockLayout.Columns < 0 {
				innerBlockLayout.Columns = 0
			}

		case r.CalculateLayoutStageWithChildren:
			if childrenBlockLayouts == nil {
				return
			}

			// children := *childrenBlockLayouts

			// propMap := map[r.Element]Properties{}

			// colsRemaining := innerBlockLayout.Columns
			// rowsRemaining := innerBlockLayout.Rows

		case r.CalculateLayoutStageFinal:

		}

		return
	}
}

func calculateOldBlockLayout(
	screen tcell.Screen,
	parentBlockLayout r.BlockLayout,
	boxProps Properties,
) (
	blockLayout r.BlockLayout,
	innerBlockLayout r.BlockLayout,
) {
	rows := parentBlockLayout.Rows
	columns := parentBlockLayout.Columns

	if rows == 0 && boxProps.Height != 0 {
		rows = int(
			math.Round(
				float64(parentBlockLayout.Rows) * (boxProps.Height / 100),
			),
		)
	}
	if columns == 0 && boxProps.Width != 0 {
		columns = int(
			math.Round(
				float64(parentBlockLayout.Columns) * (boxProps.Width / 100),
			),
		)
	}

	blockLayout = r.BlockLayout{
		ZIndex:  boxProps.ZIndex,
		Rows:    rows,
		Columns: columns,
		X:       parentBlockLayout.X,
		Y:       parentBlockLayout.Y,
	}

	// Calculate box size
	blockLayout.Columns = columns
	blockLayout.Rows = rows

	// Calculate margin
	blockLayout.X = parentBlockLayout.X + boxProps.Margin.Left
	blockLayout.Columns = blockLayout.Columns - boxProps.Margin.Right
	blockLayout.Y = parentBlockLayout.Y + boxProps.Margin.Top
	blockLayout.Rows = blockLayout.Rows - boxProps.Margin.Bottom

	innerBlockLayout = r.BlockLayout{
		ZIndex:  boxProps.ZIndex,
		Rows:    blockLayout.Rows,
		Columns: blockLayout.Columns,
		X:       blockLayout.X,
		Y:       blockLayout.Y,
	}

	innerBlockLayout.Columns = blockLayout.Columns -
		boxProps.Padding.Left - boxProps.Padding.Right

	innerBlockLayout.Rows = blockLayout.Rows -
		boxProps.Padding.Top - boxProps.Padding.Bottom

	// Calculate padding box
	innerBlockLayout.Y = innerBlockLayout.Y + boxProps.Padding.Top
	innerBlockLayout.Columns = innerBlockLayout.Columns - boxProps.Padding.Right
	innerBlockLayout.Rows = innerBlockLayout.Rows - boxProps.Padding.Bottom
	innerBlockLayout.X = innerBlockLayout.X + boxProps.Padding.Left

	// Border Sizing

	if boxProps.Border.Style != BorderStyleNone {
		blockLayout.Columns = blockLayout.Columns - 2 // 1 for each side
		blockLayout.Rows = blockLayout.Rows - 2       // 1 for each side

		innerBlockLayout.X = innerBlockLayout.X + 1
		innerBlockLayout.Y = innerBlockLayout.Y + 1
		innerBlockLayout.Rows = innerBlockLayout.Rows - 1
		innerBlockLayout.Columns = innerBlockLayout.Columns - 2
	}

	// Ensure the rows and cols are not below 0
	if blockLayout.Rows < 0 {
		blockLayout.Rows = 0
	}
	if blockLayout.Columns < 0 {
		blockLayout.Columns = 0
	}
	if innerBlockLayout.Rows < 0 {
		innerBlockLayout.Rows = 0
	}
	if innerBlockLayout.Columns < 0 {
		innerBlockLayout.Columns = 0
	}
	return
}

func calculateOldBlockLayoutForChildren(
	screen tcell.Screen,
	boxProps Properties,
	innerBlockLayout r.BlockLayout,
	children r.Children,
) r.Children {
	// if len(children) == 0 {
	// 	return children
	// }

	// propMap := map[r.Element]Properties{}

	// colsRemaining := innerBlockLayout.Columns
	// rowsRemaining := innerBlockLayout.Rows
	// flexGrowCount := 0
	// flexGrowDivision := 0

	// for _, c := range children {
	// 	if c == nil {
	// 		continue
	// 	}

	// 	propMap[c] = c.Properties.GetOptionalProperty(
	// 		Properties{},
	// 	).(Properties)
	// }

	// // Find all children with fixed row,col sizing
	// for _, props := range propMap {
	// 	colsRemaining = colsRemaining - props.Columns
	// 	rowsRemaining = rowsRemaining - props.Rows
	// 	flexGrowCount = flexGrowCount + props.FlexGrow

	// 	if props.FlexGrow == 0 {
	// 		flexGrowCount = flexGrowCount + 1 // we force flex-grow to be at least 1
	// 	}
	// }

	// switch boxProps.FlexDirection {
	// case FlexDirectionRow:
	// 	flexGrowDivision = colsRemaining / flexGrowCount
	// case FlexDirectionRowReverse:
	// 	flexGrowDivision = colsRemaining / flexGrowCount
	// case FlexDirectionColumn:
	// 	flexGrowDivision = rowsRemaining / flexGrowCount
	// case FlexDirectionColumnReverse:
	// 	flexGrowDivision = rowsRemaining / flexGrowCount

	// }

	// // Reverse the slices if needed
	// if boxProps.FlexDirection == FlexDirectionRowReverse ||
	// 	boxProps.FlexDirection == FlexDirectionColumnReverse {
	// 	for i := len(children)/2 - 1; i >= 0; i-- {
	// 		opp := len(children) - 1 - i
	// 		children[i], children[opp] = children[opp], children[i]
	// 	}
	// }

	// x := innerBlockLayout.X
	// y := innerBlockLayout.Y

	// for i, el := range children {
	// 	if el == nil {
	// 		continue
	// 	}

	// 	props := propMap[el]
	// 	flexGrow := props.FlexGrow

	// 	if props.FlexGrow == 0 {
	// 		flexGrow = flexGrow + 1 // we force flex-grow to be at least 1
	// 	}

	// 	rows := 0
	// 	columns := 0

	// 	switch boxProps.FlexDirection {
	// 	case FlexDirectionRow:
	// 		columns = flexGrowDivision * flexGrow
	// 		rows = innerBlockLayout.Rows
	// 	case FlexDirectionRowReverse:
	// 		columns = flexGrowDivision * flexGrow
	// 		rows = innerBlockLayout.Rows
	// 	case FlexDirectionColumn:
	// 		columns = innerBlockLayout.Columns
	// 		rows = flexGrowDivision * flexGrow
	// 	case FlexDirectionColumnReverse:
	// 		columns = innerBlockLayout.Columns
	// 		rows = flexGrowDivision * flexGrow
	// 	}

	// 	// Ensure rows and columns aren't negative
	// 	if rows < 0 {
	// 		rows = 0
	// 	}
	// 	if columns < 0 {
	// 		columns = 0
	// 	}

	// 	if props.MinHeight != 0 {
	// 		rows = intmath.Min(rows, props.MinHeight)
	// 	}
	// 	if props.MinWidth != 0 {
	// 		columns = intmath.Min(columns, props.MinWidth)
	// 	}

	// 	blockLayout := r.BlockLayout{
	// 		X:       x,
	// 		Y:       y,
	// 		Rows:    rows,
	// 		Columns: columns,
	// 		ZIndex:  boxProps.ZIndex,
	// 		Order:   i,
	// 	}

	// 	switch boxProps.FlexDirection {
	// 	case FlexDirectionRow:
	// 		x = x + columns
	// 	case FlexDirectionRowReverse:
	// 		x = x + columns
	// 	case FlexDirectionColumn:
	// 		y = y + rows
	// 	case FlexDirectionColumnReverse:
	// 		y = y + rows
	// 	}

	// 	el.Properties = r.ReplaceProps(el.Properties, blockLayout)
	// }
	return children
}

// calculateSizeOfChildren recurses down the tree until it finds
// a single of multiple boxes, and calculates their size
// func calculateSizeOfChildren(el r.Element) r.BlockLayout {

// }
