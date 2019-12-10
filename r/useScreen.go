package r

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
)

var useSimulationScreen bool
var useScreenInstance tcell.Screen
var hasScreenInstance bool

// UseScreen returns a tcell.Screen allowing you to read and
// interact with the Screen directly.
//
// Even though this means you can modify the Screen from
// anywhere, just as you should avoid DOM manipulation directly
// in React, you should avoid manipulating the Screen with
// this hook.
//
// Use this hook to read information from the screen only.
//
// If you need to write to the Screen, use a ScreenElement.
// This ensures when your Component has changes, retort will
// call your RenderToScreen function. Doing this any other way
// will gaurentee at some point things will get out of sync.
func UseScreen() tcell.Screen {
	if hasScreenInstance {
		return useScreenInstance
	}

	var s tcell.Screen
	var err error

	if c.UseSimulationScreen {
		s = tcell.NewSimulationScreen("UTF-8")
	} else {
		s, err = tcell.NewScreen()
	}
	useScreenInstance = s
	encoding.Register()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	hasScreenInstance = true
	return useScreenInstance
}
