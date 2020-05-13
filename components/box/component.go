package box

import (
	"github.com/gdamore/tcell"
	"retort.dev/debug"

	"retort.dev/r"
)

// Box is the basic building block for a retort app.
// Box implements the Box Model, see Properties
func Box(p r.Properties) r.Element {
	screen := r.UseScreen()

	// Get our Properties
	boxProps := p.GetProperty(
		Properties{},
		"Box requires Properties",
	).(Properties)

	// Get any children
	children := p.GetOptionalProperty(
		r.Children{},
	).(r.Children)

	// // Calculate the BlockLayout of this Box
	// blockLayout, innerBlockLayout := calculateBlockLayout(
	// 	screen,
	// 	parentBlockLayout,
	// 	boxProps,
	// )

	// // Calculate the BlockLayout of any children
	// childrenWithLayout := calculateBlockLayoutForChildren(
	// 	screen,
	// 	boxProps,
	// 	innerBlockLayout,
	// 	children,
	// )

	return r.CreateScreenElement(
		calculateBlockLayout(boxProps),
		func(s tcell.Screen, blockLayout r.BlockLayout) {
			// debug.Spew(p)

			if s == nil {
				panic("Box can't render no screen")
			}

			w, h := s.Size()

			if w == 0 || h == 0 {
				panic("Box can't render on a zero size screen")
			}

			debug.Spew("CreateScreenElement Box", w, h, blockLayout)
			render(
				screen,
				boxProps,
				blockLayout,
			)

		},
		r.Properties{},
		children,
	)
}
