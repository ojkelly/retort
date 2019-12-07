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
		"Box requires Properties",
	).(Properties)

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
				panic("Box can't render no screen")
			}

			w, h := s.Size()

			if w == 0 || h == 0 {
				panic("Box can't render on a zero size screen")
			}

			render(
				screen,
				boxProps,
				boxLayout,
			)

			return boxLayout
		},
		childrenWithLayout,
	)
}
