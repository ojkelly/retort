package theme

import "github.com/gdamore/tcell"

type Base16 struct {
	// Base00 - Default Background
	Base00 tcell.Color
	// Base01 - Lighter Background (Used for status bars)
	Base01 tcell.Color
	// Base02 - Selection Background
	Base02 tcell.Color
	// Base03 - Comments, Invisibles, Line Highlighting
	Base03 tcell.Color
	// Base04 - Dark Foreground (Used for status bars)
	Base04 tcell.Color
	// Base05 - Default Foreground, Caret, Delimiters, Operators
	Base05 tcell.Color
	// Base06 - Light Foreground (Not often used)
	Base06 tcell.Color
	// Base07 - Light Background (Not often used)
	Base07 tcell.Color
	// Base08 - Variables, XML Tags, Markup Link Text, Markup Lists, Diff Deleted
	Base08 tcell.Color
	// Base09 - Integers, Boolean, Constants, XML Attributes, Markup Link Url
	Base09 tcell.Color
	// Base0A - Classes, Markup Bold, Search Text Background
	Base0A tcell.Color
	// Base0B - Strings, Inherited Class, Markup Code, Diff Inserted
	Base0B tcell.Color
	// Base0C - Support, Regular Expressions, Escape Characters, Markup Quotes
	Base0C tcell.Color
	// Base0D - Functions, Methods, Attribute IDs, Headings
	Base0D tcell.Color
	// Base0E - Keywords, Storage, Selector, Markup Italic, Diff Changed
	Base0E tcell.Color
	// Base0F - Deprecated, Opening/Closing Embedded Language Tags, e.g. <?php ?>
	Base0F tcell.Color
}
