package r

import (
	"sync"

	"github.com/gdamore/tcell"
)

var c *RetortConfiguration = &RetortConfiguration{}

// quitChan will quit the application when anythin is sent to it
var quitChan chan struct{}

// resizeChan notifies when the screen has resized, and needs to redraw
var resizeChan chan struct{}

// [ Hooks ]--------------------------------------------------------------------

// hookIndex keeps track of the currently used hook, this is a proxy for
// call index
var hookIndex int

var hookFiber *fiber
var hookFiberLock = &sync.Mutex{}

// setStateChan is a channel where SetState actions are sent to be processed
// in the workloop
var setStateChan chan ActionCreator

// [ UseScreen ]----------------------------------------------------------------

var useSimulationScreen bool
var useScreenInstance tcell.Screen
var hasScreenInstance bool
