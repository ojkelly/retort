package theme

import "github.com/gdamore/tcell"

type Color int

const (
	Orange Color = iota
	White
)

type Colors struct {
	Border     tcell.Color
	Accent     tcell.Color
	Foreground tcell.Color
	Subtle     tcell.Color
}

var orange Colors = Colors{
	Border:     tcell.ColorOrange,
	Accent:     tcell.ColorOrange,
	Foreground: tcell.ColorWhite,
	Subtle:     tcell.ColorGrey,
}

var white Colors = Colors{
	Border:     tcell.ColorGrey,
	Accent:     tcell.ColorWhite,
	Foreground: tcell.ColorWhite,
	Subtle:     tcell.ColorGrey,
}
