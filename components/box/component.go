package box

import (
	"github.com/gdamore/tcell"

	"retort.dev/r"
)

// Box is the basic building block for a retort app.
// Box implements the Box Model, see Properties
func Box(p r.Properties) r.Element {
	screen := r.UseScreen()

	// Get our Properties
	boxProps := p.GetProperty(
		Properties{},
		"Box requires box.Properties",
	).(Properties)

	// Get any children
	children := p.GetOptionalProperty(
		r.Children{},
	).(r.Children)

	return r.CreateScreenElement(
		calculateBlockLayout(boxProps),
		func(s tcell.Screen, blockLayout r.BlockLayout) {
			if s == nil {
				panic("Box can't render no screen")
			}

			w, h := s.Size()

			if w == 0 || h == 0 {
				panic("Box can't render on a zero size screen")
			}

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
