package retort_test

import (
	// import tcell to get access to colors and event types
	"github.com/gdamore/tcell"

	"retort.dev/component"
	"retort.dev/example"
	"retort.dev/r"
)

var exampleVarToMakeGoDocPrintTheWholeFile bool

func Example_app() {
	// Call the main function on retort to start the app,
	// when you call this, retort will take over the screen.
	r.Retort(
		// Root Element
		r.CreateElement(
			example.ClickableBox,
			r.Properties{
				component.BoxProps{
					Width:  100, // Make the root element fill the screen
					Height: 100, // Make the root element fill the screen
					Border: component.Border{
						Style:      component.BorderStyleSingle,
						Foreground: tcell.ColorWhite,
					},
				},
			},
			r.Children{
				// First Child
				r.CreateElement(
					example.ClickableBox,
					r.Properties{
						component.BoxProps{
							Border: component.Border{
								Style:      component.BorderStyleSingle,
								Foreground: tcell.ColorWhite,
							},
						},
					},
					nil, // Pass nil as the third argument if there are no children
				),
				// Second Child
				r.CreateElement(
					example.ClickableBox,
					r.Properties{
						component.BoxProps{
							Border: component.Border{
								Style:      component.BorderStyleSingle,
								Foreground: tcell.ColorWhite,
							},
						},
					},
					nil,
				),
			},
		),
		// Pass in optional configuration
		r.RetortConfiguration{},
	)
}
