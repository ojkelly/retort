package r

import (
	"github.com/gdamore/tcell"
	"retort.dev/r/debug"
	"retort.dev/r/internal/quadtree"
)

// TODO: Rewrite to enable event propegation up the chain
// TODO: Enable focussed elements, receiving the event first
type (
	// EventHandler is a Property you can add to a Component that will
	// be called on every *tcell.Event that is created.
	//
	// Use this sparingly as it's very noisy.
	EventHandler = func(e *tcell.Event)

	EventHandlerKey = func(e *tcell.EventKey, meta EventMeta) EventMeta

	// EventHandlerMouse is a Property you can add to a Component to
	// be called when a *tcell.EventMouse is created.
	EventHandlerMouse = func(e *tcell.EventMouse)

	// EventHandlerMouseHover is called when a mouse is over your Component
	EventHandlerMouseHover = func()

	EventMouseScroll = func(up, down, left, right bool)

	// EventMouseClick is called when a mouse clicks on your component.
	// For conveince we pass isPrimary and isSecondary as aliases for
	// Button1 and Button2.
	EventMouseClick = func(
		isPrimary,
		isSecondary bool,
		buttonMask tcell.ButtonMask,
	) EventMouseClickRelease

	// EventMouseClickDrag is not yet implemented, but could be called to allow
	// a component to render a version that is being dragged around
	EventMouseClickDrag = func()

	// EventMouseClickRelease is called when the mouse click has been released.
	// TODO: this can probably be enhanced to enable drag and drop
	EventMouseClickRelease = func()

	EventMeta struct {
		StopPropegation bool
	}
)

// TODO: direct hover and click events
// TODO: keep track of focussed inputs, and direct keyboard input there, when
// focussed
func (r *retort) handleEvents(resizeChan chan struct{}) {
	screen := UseScreen()
	quit := UseQuit()

	meta := EventMeta{}
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
			r.handleKeyEvent(r.wipRoot, ev, meta)
		case *tcell.EventResize:
			resizeChan <- struct{}{}

		case *tcell.EventMouse:
			r.handleMouseEvent(ev)
		case *tcell.EventError:
		case *tcell.EventInterrupt:
		case *tcell.EventTime:
		default:
			if ev != nil {
				debug.Log("Unhandled Event", ev)
			}
		}

	}
}

func (r *retort) handleEvent(e tcell.Event) {

	// Search the quadtree for the matching fiber

	// Get the event handler and call it

}

// handleKeyEvent currently processes outside in
// it would be ideal to track the element/fiber in focus, and expand out from there
func (r *retort) handleKeyEvent(f *fiber, ev *tcell.EventKey, meta EventMeta) {
	// debug.Spew(f)
	if f == nil {
		return
	}

	var eventHandler EventHandlerKey

	eventHandler = f.Properties.GetOptionalProperty(
		eventHandler,
	).(EventHandlerKey)

	if eventHandler != nil {
		meta = eventHandler(ev, meta)
	}

	debug.Spew("handleKeyEvent", ev.Rune(), meta, eventHandler)
	if meta.StopPropegation {
		return
	}

	r.handleKeyEvent(f.child, ev, meta)
	r.handleKeyEvent(f.sibling, ev, meta)
}

// handleMouseEvent determines what type of mouse event needs to be created
// and then routes that event to the correct Component
func (r *retort) handleMouseEvent(ev *tcell.EventMouse) {
	if ev == nil {
		return
	}

	var eventMouseClick EventMouseClick
	var eventHandlerMouseHover EventHandlerMouseHover
	var eventMouseScroll EventMouseScroll
	var smallestArea int

	var isHover bool

	var isClick,
		isPrimaryClick,
		isSecondaryClick bool

	// Vars for EventMouseScroll
	var isScroll,
		scrollDirectionUp,
		scrollDirectionDown,
		scrollDirectionLeft,
		scrollDirectionRight bool

	x, y := ev.Position()

	cursor := quadtree.Bounds{
		X:      x,
		Y:      y,
		Width:  0,
		Height: 0,
	}

	results := r.quadtree.RetrieveIntersections(cursor)

	// Determine the type of mouse event
	switch ev.Buttons() {
	// Scroll Events
	case tcell.WheelUp:
		isScroll = true
		scrollDirectionUp = true
	case tcell.WheelDown:
		isScroll = true
		scrollDirectionDown = true
	case tcell.WheelLeft:
		isScroll = true
		scrollDirectionLeft = true
	case tcell.WheelRight:
		isScroll = true
		scrollDirectionRight = true
		// Hover event?
	case tcell.ButtonNone:
	// ??

	// Click Events
	case tcell.Button1:
		isClick = true
		isPrimaryClick = true
	case tcell.Button2:
		isClick = true
		isSecondaryClick = true
	case tcell.Button3:
		isClick = true
	case tcell.Button4:
		isClick = true
	case tcell.Button5:
		isClick = true
	case tcell.Button6:
		isClick = true
	case tcell.Button7:
		isClick = true
	case tcell.Button8:
		isClick = true
	default:
		// ??

	}

	var eventHandlerProp interface{}
	// Search the matching Components and find the handler
	for _, r := range results {
		if r.Value == nil {
			continue
		}

		// Grab the event handler from this fiber
		matchingFiber := r.Value.(*fiber)

		switch {
		case isClick:
			eventMouseClick = matchingFiber.Properties.GetOptionalProperty(
				eventMouseClick,
			).(EventMouseClick)
		case isHover:
			eventHandlerMouseHover = matchingFiber.Properties.GetOptionalProperty(
				eventHandlerMouseHover,
			).(EventHandlerMouseHover)
		case isScroll:
			eventMouseScroll = matchingFiber.Properties.GetOptionalProperty(
				eventMouseScroll,
			).(EventMouseScroll)
		}

		if eventHandlerProp == nil {
			continue
		}

		match := false
		if cursor.Intersects(r) {
			match = true
		}

		bl := r.Value.(*fiber).BlockLayout

		// find the area of the box
		area := bl.Columns * bl.Rows
		if smallestArea == 0 || smallestArea > area {
			match = true
		}

		if match {
			smallestArea = area
		}
	}

	// Call the event handler from the component, or return if none found
	switch {
	case isClick:
		if eventMouseClick == nil {
			return
		}
		eventMouseClick(isPrimaryClick, isSecondaryClick, ev.Buttons())
	case isHover:
		if eventHandlerMouseHover == nil {
			return
		}
		eventHandlerMouseHover()

	case isScroll:
		if eventMouseScroll == nil {
			return
		}
		eventMouseScroll(
			scrollDirectionUp,
			scrollDirectionDown,
			scrollDirectionLeft,
			scrollDirectionRight,
		)

	}

}

func handleScrollEvent(ev *tcell.EventMouse) {
	//
}
