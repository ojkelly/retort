package text

import (
	"github.com/gdamore/tcell"
	"retort.dev/components/box"
	"retort.dev/r"
)

func calculateBlockLayout(
	textProps Properties,
	boxProps box.Properties,
) r.CalculateLayout {
	return func(
		s tcell.Screen,
		stage r.CalculateLayoutStage,
		parentBlockLayout r.BlockLayout,
		children []r.BlockLayoutWithProperties,
	) (
		outerBlockLayout r.BlockLayout,
		innerBlockLayout r.BlockLayout,
		childrenBlockLayouts []r.BlockLayoutWithProperties,
	) {
		outerBlockLayout = parentBlockLayout
		innerBlockLayout = parentBlockLayout
		childrenBlockLayouts = children
		// debug.Spew(stage, outerBlockLayout)

		switch stage {
		case r.CalculateLayoutStageInitial:

			lines := breakLines(textProps, innerBlockLayout)
			rows := len(lines)

			outerBlockLayout.Rows = rows
			outerBlockLayout.FixedRows = true
			innerBlockLayout.Rows = rows
			innerBlockLayout.FixedRows = true
		case r.CalculateLayoutStageWithChildren:
		case r.CalculateLayoutStageFinal:

		}

		return
	}
}
