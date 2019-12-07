package text

import (
	"github.com/gdamore/tcell"
	"retort.dev/debug"
	"retort.dev/r"
)

// Text is the basic building block for a retort app.
// Text implements the Text Model, see Properties
func Text(p r.Properties) r.Element {
	screen := r.UseScreen()

	// Get our Properties
	textProps := p.GetProperty(
		Properties{},
		"Text requires Properties",
	).(Properties)

	// Get our BoxLayout
	parentBoxLayout := p.GetProperty(
		r.BoxLayout{},
		"Text requires a parent BoxLayout.",
	).(r.BoxLayout)
	debug.Spew("text parentBoxLayout", parentBoxLayout)
	// Get any children
	children := p.GetOptionalProperty(
		r.Children{},
	).(r.Children)
	if len(children) != 0 {
		panic("Text cannot have children")
	}

	// Calculate the BoxLayout of this Text
	boxLayout := calculateBoxLayout(
		screen,
		parentBoxLayout,
		textProps,
	)

	return r.CreateScreenElement(
		func(s tcell.Screen) r.BoxLayout {
			if s == nil {
				panic("Text can't render no screen")
			}

			w, h := s.Size()

			if w == 0 || h == 0 {
				panic("Text can't render on a zero size screen")
			}

			renderText(
				s,
				textProps,
				boxLayout,
			)

			return boxLayout
		},
		nil,
	)
}
