package text

import (
	"github.com/gdamore/tcell"
	"retort.dev/components/box"
	"retort.dev/r"
	"retort.dev/r/debug"
)

func calculateBlockLayout(
	textProps Properties,
	boxProps box.Properties,
) r.CalculateLayout {
	return func(
		s tcell.Screen,
		stage r.CalculateLayoutStage,
		parentBlockLayout r.BlockLayout,
		children r.BlockLayouts,
	) (
		outerBlockLayout r.BlockLayout,
		innerBlockLayout r.BlockLayout,
		childrenBlockLayouts r.BlockLayouts,
	) {
		outerBlockLayout = parentBlockLayout
		innerBlockLayout = parentBlockLayout
		childrenBlockLayouts = children
		// debug.Spew(stage, outerBlockLayout)

		switch stage {
		case r.CalculateLayoutStageInitial:

			lines := breakLines(textProps, innerBlockLayout)
			rows := len(lines)

			debug.Spew("rows", rows, textProps, innerBlockLayout)

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
