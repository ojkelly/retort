package r

import (
	"reflect"
	"sync"
)

type (
	Action        = func(s State) State
	ActionCreator struct {
		h *hook
		a Action
	}
	SetState = func(a Action)
)

// UseState provides local state for a Component.
//
// With UseState you can make your components interactive,
// and repsonsive to either user input or anything else that
// changes.
//
// UseState by itself only gives you the ability to change state,
// you generally need to pair this with either an EventHandler
// or UseEffect to provide interactivity.
//
// Don't call setState inside your Component, as this will create
// an infinite rendering loop.
//
// Example
//
// The following example shows how you can use state to change
// the color of a Box border when it's clicked.
//
//  import (
//    "github.com/gdamore/tcell"
//    "retort.dev/component"
//    "retort.dev/r/debug"
//    "retort.dev/r"
//  )
//
//  type MovingBoxState struct {
//    Color tcell.Color
//  }
//
//  func ClickableBox(p r.Properties) r.Element {
//    boxProps := p.GetProperty(
//      component.BoxProps{},
//      "Container requires ContainerProps",
//    ).(component.BoxProps)
//
//    children := p.GetProperty(
//      r.Children{},
//      "Container requires r.Children",
//    ).(r.Children)
//
//    s, setState := r.UseState(r.State{
//      MovingBoxState{Color: boxProps.Border.Foreground},
//    })
//    state := s.GetState(
//      MovingBoxState{},
//    ).(MovingBoxState)
//
//    mouseEventHandler := func(e *tcell.EventMouse) {
//      debug.Log("mouseEventHandler", e, state)
//      color := tcell.ColorGreen
//      if state.Color == tcell.ColorGreen {
//        color = tcell.ColorBlue
//      }
//
//      if state.Color == tcell.ColorBlue {
//        color = tcell.ColorGreen
//      }
//
//      setState(func(s r.State) r.State {
//        debug.Log("mouseEventHandler update state", e, state)
//
//        return r.State{MovingBoxState{
//          Color: color,
//        },
//        }
//      })
//    }
//
//    boxProps.Border.Foreground = state.Color
//
//    return r.CreateElement(
//      component.Box,
//      r.Properties{
//        boxProps,
//        mouseEventHandler,
//      },
//      children,
//    )
//  }
func UseState(initial State) (State, SetState) {
	hookFiberLock.Lock()
	checkStateTypesAreUnique(initial)
	var oldHook *hook

	if hookFiber != nil &&
		hookFiber.alternate != nil &&
		hookFiber.alternate.hooks != nil &&
		len(hookFiber.alternate.hooks) > hookIndex &&
		hookFiber.alternate.hooks[hookIndex] != nil {
		oldHook = hookFiber.alternate.hooks[hookIndex]
	}

	var h *hook
	if oldHook != nil {
		h = oldHook
	} else {
		h = &hook{
			tag:   hookTagState,
			mutex: &sync.Mutex{},
			state: initial,
		}
	}

	var actions []Action
	if oldHook != nil {
		actions = oldHook.queue

		for _, action := range actions {
			h.state = action(h.state)
		}
	}

	if hookFiber != nil {
		hookFiber.hooks = append(hookFiber.hooks, h)
		hookIndex++
	}

	hookFiberLock.Unlock()

	return h.state, h.setState()
}

func (h *hook) setState() SetState {
	return func(a Action) {
		setStateChan <- ActionCreator{
			h: h,
			a: a,
		}
	}
}

func (ac ActionCreator) addToQueue() {
	ac.h.mutex.Lock()
	if ac.h.mutex == nil {
		panic("h is gone")
	}

	ac.h.queue = append(ac.h.queue, ac.a)
	ac.h.mutex.Unlock()

}

func checkStateTypesAreUnique(s State) bool {
	seenPropTypes := make(map[reflect.Type]bool)

	for _, p := range s {
		if seen := seenPropTypes[reflect.TypeOf(p)]; seen {
			return false
		}
		seenPropTypes[reflect.TypeOf(p)] = true
	}
	return true
}

// GetState lets you retrieve the state of your passed in type from a UseState
// hook.
//
// Because we cannot use generics this is the closest we can get. This is like
// Properties where the stateType type is a key to the struct in the slice of
// interfaces.
// As such, you can only have one of a given type in state.
func (state State) GetState(stateType interface{}) interface{} {
	for _, p := range state {
		if reflect.TypeOf(p) == reflect.TypeOf(stateType) {
			return p
		}
	}
	return stateType
}
