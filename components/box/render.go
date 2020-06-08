package box

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell"
	runewidth "github.com/mattn/go-runewidth"
	"retort.dev/r"
)

func render(
	s tcell.Screen,
	props Properties,
	layout r.BlockLayout,
) {
	// debug.Spew("render", layout)
	x1 := layout.X
	y1 := layout.Y
	x2 := layout.X + layout.Columns
	y2 := layout.Y + layout.Rows

	borderStyle := tcell.StyleDefault
	borderStyle = borderStyle.Foreground(props.Border.Foreground)
	gl := ' '

	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	if props.Border.Style != BorderStyleNone {
		for col := x1; col <= x2; col++ {
			s.SetContent(col, y1, tcell.RuneHLine, nil, borderStyle)
			s.SetContent(col, y2, tcell.RuneHLine, nil, borderStyle)
		}
		for row := y1 + 1; row < y2; row++ {
			s.SetContent(x1, row, tcell.RuneVLine, nil, borderStyle)
			s.SetContent(x2, row, tcell.RuneVLine, nil, borderStyle)
		}
		if y1 != y2 && x1 != x2 {
			// Only add corners if we need to
			s.SetContent(x1, y1, tcell.RuneULCorner, nil, borderStyle)
			s.SetContent(x2, y1, tcell.RuneURCorner, nil, borderStyle)
			s.SetContent(x1, y2, tcell.RuneLLCorner, nil, borderStyle)
			s.SetContent(x2, y2, tcell.RuneLRCorner, nil, borderStyle)
		}
		for row := y1 + 1; row < y2; row++ {
			for col := x1 + 1; col < x2; col++ {
				s.SetContent(col, row, gl, nil, borderStyle)
			}
		}
	}

	if props.Title.Value != "" {
		renderLabel(
			s,
			props.Title,
			r.BlockLayout{
				X:       layout.X + 2, // Bump it over 1 for the corner, and 1 for style
				Y:       layout.Y,
				Rows:    1,
				Columns: layout.Columns,
			},
			borderStyle,
		)
	}
	if props.Footer.Value != "" {
		renderLabel(
			s,
			props.Footer,
			r.BlockLayout{
				X:       layout.X + 2, // Bump it over 1 for the corner, and 1 for style
				Y:       layout.Y + layout.Rows,
				Rows:    1,
				Columns: layout.Columns,
			},
			borderStyle,
		)
	}
}

func renderLabel(
	s tcell.Screen,
	label Label,
	layout r.BlockLayout,
	style tcell.Style,
) {

	i := 0
	var deferred []rune
	dwidth := 0
	isZeroWidthJoiner := false

	wrapLeft := ""
	wrapRight := ""

	switch label.Wrap {
	case LabelWrapNone:
		wrapLeft = " "
		wrapRight = " "
	case LabelWrapBrace:
		wrapLeft = "{ "
		wrapRight = " }"
	case LabelWrapBracket:
		wrapLeft = "( "
		wrapRight = " )"
	case LabelWrapSquareBracket:
		wrapLeft = "[ "
		wrapRight = " ]"
	case LabelWrapChevron:
		wrapLeft = "< "
		wrapRight = " >"
	}

	wrapLeft = fmt.Sprintf(
		"%s%s%s",
		strings.Repeat(" ", label.Margin.Left),
		wrapLeft,
		strings.Repeat(" ", label.Padding.Left),
	)
	wrapRight = fmt.Sprintf(
		"%s%s%s",
		strings.Repeat(" ", label.Margin.Right),
		wrapRight,
		strings.Repeat(" ", label.Padding.Right),
	)

	value := fmt.Sprintf("%s%s%s", wrapLeft, label.Value, wrapRight)

	// Print each rune to the screen
	for _, r := range value {
		// Check if the rune is a Zero Width Joiner
		if r == '\u200d' {
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				dwidth = 1
			}
			deferred = append(deferred, r)
			isZeroWidthJoiner = true
			continue
		}

		if isZeroWidthJoiner {
			deferred = append(deferred, r)
			isZeroWidthJoiner = false
			continue
		}

		switch runewidth.RuneWidth(r) {
		case 0:
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				dwidth = 1
			}
		case 1:
			if len(deferred) != 0 {
				s.SetContent(layout.X+i, layout.Y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 1
		case 2:
			if len(deferred) != 0 {
				s.SetContent(layout.X+i, layout.Y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 2
		}
		deferred = append(deferred, r)
	}
	if len(deferred) != 0 {
		s.SetContent(layout.X+i, layout.Y, deferred[0], deferred[1:], style)
	}
}
