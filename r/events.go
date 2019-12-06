package r

import "github.com/gdamore/tcell"

import "retort.dev/r/internal/quadtree"

type (
	// EventHandler is a Property you can add to a Component that will
	// be called on every *tcell.Event that is created.
	//
	// Use this sparingly as it's very noisy.
	EventHandler = func(e *tcell.Event)
	// MouseEventHandler is a Property you can add to a Component to
	// be called when a *tcell.EventMouse is created.
	MouseEventHandler = func(e *tcell.EventMouse)
)

// TODO: direct hover and click events
// TODO: keep track of focussed inputs, and direct keyboard input there, when
// focussed
func (r *retort) handleEvents() {
	screen := UseScreen()
	quit := UseQuit()

	for {
		// Grab events from tcell
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			// Keyboard event
			switch ev.Key() {
			case tcell.KeyEscape:
				quit()
			case tcell.KeyCtrlQ:
				quit()
			}
		case *tcell.EventResize:
			w, h := screen.Size()

			r.quadtree.Bounds.Width = w
			r.quadtree.Bounds.Height = h

			screen.Sync()
		case *tcell.EventMouse:
			r.handleMouseEvent(ev)
		case *tcell.EventError:
		case *tcell.EventInterrupt:
		case *tcell.EventTime:
		default:
			if ev != nil {
				// debug.Spew("Unhandled Event", ev)
			}
		}

	}
}

func (r *retort) handleEvent(e tcell.Event) {

	// Search the quadtree for the matching fiber

	// Get the event handler and call it

}

func (r *retort) handleMouseEvent(ev *tcell.EventMouse) {
	if ev == nil {
		return
	}

	x, y := ev.Position()
	button := ev.Buttons()

	// Only buttons, not wheel events
	button &= tcell.ButtonMask(0xff)

	if button == tcell.ButtonNone {
		return
	}

	switch ev.Buttons() {
	case tcell.Button1:
	case tcell.Button2:
	case tcell.Button3:
	case tcell.Button4:
	case tcell.Button5:
	case tcell.Button6:
	case tcell.Button7:
	case tcell.Button8:
	}

	cursor := quadtree.Bounds{
		X:      x,
		Y:      y,
		Width:  0,
		Height: 0,
	}

	results := r.quadtree.RetrieveIntersections(cursor)

	var eventHandler MouseEventHandler
	var f *fiber
	var smallestArea int

	for _, r := range results {
		if r.Value == nil {
			continue
		}
		matchingFiber := r.Value.(*fiber)

		eventHandlerProp := matchingFiber.Properties.GetOptionalProperty(
			eventHandler,
		).(MouseEventHandler)

		if eventHandlerProp == nil {
			continue
		}

		if f == nil {
			f = matchingFiber
		}

		match := false
		if cursor.Intersects(r) {
			match = true
		}

		bl := r.Value.(*fiber).boxLayout

		// find the area of the box
		area := bl.Columns * bl.Rows
		if smallestArea == 0 || smallestArea > area {
			match = true
		}

		if match {
			f = r.Value.(*fiber)
			eventHandler = eventHandlerProp
			smallestArea = area
		}
	}

	if f == nil || eventHandler == nil {
		return
	}

	eventHandler(ev)
}
