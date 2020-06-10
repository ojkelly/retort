package text

import (
	"strings"

	"github.com/gdamore/tcell"
	runewidth "github.com/mattn/go-runewidth"
	"retort.dev/r"
	"retort.dev/r/debug"
	"retort.dev/r/intmath"
)

func render(
	s tcell.Screen,
	props Properties,
	layout r.BlockLayout,
	offsetX, offsetY int,
) {

	style := tcell.StyleDefault
	style = style.Foreground(props.Foreground)

	var lines []string
	for _, text := range strings.Split(props.Value, "\n\n") {
		lines = append(lines, breakText(text, props, layout)...)
		lines = append(lines, "")
	}

	scrollLimit := int(float64(len(lines)) / 1.2)
	offset := 0
	if offsetX < len(lines) {
		offset = offsetX
	}
	if offsetX > scrollLimit {
		offset = scrollLimit
	}

	linesToRender := lines[offset:]

	for i, line := range linesToRender {
		if i > layout.Rows {
			return
		}

		renderLine(s, style, layout.X, layout.Y+i, line)
	}
}

// breakText into rows to text that can be printed.
// This function handles all logic related to word breaking.
func breakText(
	text string,
	props Properties,
	layout r.BlockLayout,
) (lines []string) {
	width := layout.Columns

	// Break up words by whitespace characters
	words := strings.Fields(text)

	// if there's no words bail here
	if len(words) == 0 {
		return
	}

	line := ""
	colsRemaining := width

	for _, word := range words {
		if colsRemaining == 0 {
			// Save this line
			lines = append(lines, line)

			// And make a new one
			line = word
			colsRemaining = width
			continue
		}

		if len(word)+2 > colsRemaining {
			// Can we break the word?
			if props.WordBreak == BreakAll {
				lengthToSplit := intmath.Min(len(word), colsRemaining-1)
				// TODO: this isn't great, and could be greatly improved
				wordPart := word[:lengthToSplit] + "-"
				line = line + wordPart
				word = word[lengthToSplit:]
			}

			// Save this line
			lines = append(lines, line)

			// And make a new one
			line = word
			colsRemaining = width
			continue
		}

		line = line + word + " "
		colsRemaining = colsRemaining - len(word) - 1
		if colsRemaining < 0 {
			colsRemaining = 0
		}
	}

	// TODO: there's probably a better way
	// save last line
	lines = append(lines, line)

	return
}

func renderLine(s tcell.Screen, style tcell.Style, x, y int, str string) {
	debug.Spew("renderLine", str)
	i := 0
	var deferred []rune
	dwidth := 0
	zwj := false
	for _, r := range str {
		if r == '\u200d' {
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				dwidth = 1
			}
			deferred = append(deferred, r)
			zwj = true
			continue
		}
		if zwj {
			deferred = append(deferred, r)
			zwj = false
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
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 1
		case 2:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 2
		}
		deferred = append(deferred, r)
	}
	if len(deferred) != 0 {
		s.SetContent(x+i, y, deferred[0], deferred[1:], style)
		i += dwidth
	}
}
