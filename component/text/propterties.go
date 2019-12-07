package text

import "github.com/gdamore/tcell"

type Properties struct {
	Value    string
	Overflow Overflow

	Background tcell.Color
	Foreground tcell.Color
}

// Overflow controls if the text is allowed to spill outside it's contained
// While WordBreak controls what to do in the case of OverflowWrap
type Overflow int

const (
	OverflowWrap Overflow = iota
	OverflowElipsis
	OverflowHidden
)

type Align int

const (
	AlignLeft Align = iota
	AlignCenter
	AlignRight
)

// WordBreak controls what happens to text greater than its width when
// OverflowWrap is selected.
type WordBreak int

const (
	// Normal Use the default line break rule.
	Normal WordBreak = iota
	// BreakAll To prevent overflow, word breaks should be inserted between any
	// two characters (excluding Chinese/Japanese/Korean text).
	BreakAll
	// KeepAll Word breaks should not be used for Chinese/Japanese/Korean (CJK)
	// text. Non-CJK text behavior is the same as for normal.
	KeepAll
)
