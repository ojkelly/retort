package text

import (
	"github.com/gdamore/tcell"
	"retort.dev/components/box"
	"retort.dev/r"
)

func calculateBlockLayout(
	props box.Properties,
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
		return

	}
}
