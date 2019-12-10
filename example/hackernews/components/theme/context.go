package theme

import (
	"github.com/gdamore/tcell"
	"retort.dev/component/box"
	"retort.dev/r"
)

var orange = box.Properties{
	Border: box.Border{
		Style:      box.BorderStyleSingle,
		Foreground: tcell.ColorOrange,
	},
}

var white = box.Properties{
	Border: box.Border{
		Style:      box.BorderStyleSingle,
		Foreground: tcell.ColorWhite,
	},
}

var Context = r.CreateContext(r.State{
	box.Properties{
		Border: box.Border{
			Style:      box.BorderStyleSingle,
			Foreground: tcell.ColorOrange,
		},
	},
})

type Color int

const (
	Orange Color = iota
	White
)

type Properties struct {
	Color Color
}

func Theme(p r.Properties) r.Element {
	children := p.GetProperty(
		r.Children{},
		"Theme requires r.Children",
	).(r.Children)

	props := p.GetProperty(
		Properties{
			Color: Orange,
		},
		"Theme requires Properties",
	).(Properties)

	state := box.Properties{}
	switch props.Color {
	case Orange:
		state = orange
	case White:
		state = white
	}

	Context.SetState(r.State{state})

	return r.CreateFragment(children)
}
