package r

import (
	"fmt"
	"os"
	"reflect"
	"sync"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"retort.dev/debug"
)

// func UseContext() context.Context {
// 	return r.ctx
// }

// [ Hooks ]----------------------------------------------------------------------------------------

var hookIndex int

type hook struct {
	tag   hookTag
	mutex *sync.Mutex

	// State
	state State
	queue []Action

	// Effect
	deps   EffectDependencies
	effect Effect
	cancel EffectCancel
}

type hookTag int

const (
	hookTagNothing hookTag = iota
	hookTagState
	hookTagEffect
	hookTagReducer
)

func (h *hook) Clone() *hook {
	return &hook{
		tag:   h.tag,
		mutex: h.mutex,

		state: h.state,
		queue: h.queue,

		deps:   h.deps,
		effect: h.effect,
		cancel: h.cancel,
	}
}

// [ UseState ]-------------------------------------------------------------------------------------

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
//    "retort.dev/debug"
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
//      debug.Spew("mouseEventHandler", e, state)
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
//        debug.Spew("mouseEventHandler update state", e, state)
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
		debug.Spew("using old hook", oldHook)
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
		debug.Spew("updated state", h.state, actions)
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

// GetState lets you retrieve the state of your passed in type from a UseState hook.
//
// Because we cannot use generics this is the closest we can get. This is like Properties
// where the stateType type is a key to the struct in the slice of interfaces.
// As such, you can only have one of a given type in state.
func (state State) GetState(stateType interface{}) interface{} {
	for _, p := range state {
		if reflect.TypeOf(p) == reflect.TypeOf(stateType) {
			return p
		}
	}
	return stateType
}

// [ UseEffect ]------------------------------------------------------------------------------------

type (
	// Effect is the function type you pass to UseEffect.
	//
	// It must return an EffectCancel, even if there is nothing to clean up.
	//
	// In the Effect you can have a routine to do something (such as fetching data),
	// and then call SetState from a UseState hook, to update your Component.
	Effect func() EffectCancel
	// EffectCancel is a function that must be returned by your Effect, and is called
	// when the effect is cleaned up or canceled. This allows you to finish anything
	// you were doing such as closing channels, connections or files.
	EffectCancel func()
	// EffectDependencies lets you pass in what your Effect depends upon.
	// If they change, your Effect will be re-run.
	EffectDependencies []interface{}
)

// UseEffect is a retort hook that can be called in your Component to run side effects.
//
// Where UseState allows your components to re-render when their State changes, UseEffect
// allows you to change that state when you need to.
//
// Data fetching is a good example of when you would want something like UseEffect.
//
// Example
//
// The example below is a reasonably simple one that changes the color of the border of a
// box ever 2 seconds. The point here is to show how you can run a goroutine in the UseEffect
// callback, and clean up the channel in the EffectCancel return function.
//  import (
//   "time"
//
//   "github.com/gdamore/tcell"
//   "retort.dev/component"
//   "retort.dev/r"
//  )
//
//  type EffectExampleBoxState struct {
//    Color tcell.Color
//  }
//
//  func EffectExampleBox(p r.Properties) r.Element {
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
//      EffectExampleBoxState{Color: boxProps.Border.Foreground},
//    })
//    state := s.GetState(
//      EffectExampleBoxState{},
//    ).(EffectExampleBoxState)
//
//    r.UseEffect(func() r.EffectCancel {
//      ticker := time.NewTicker(2 * time.Second)
//
//      done := make(chan bool)
//
//      go func() {
//        for {
//          select {
//          case <-done:
//            return
//          case <-ticker.C:
//            setState(func(s r.State) r.State {
//              ms := s.GetState(
//                EffectExampleBoxState{},
//              ).(EffectExampleBoxState)
//
//              color := tcell.ColorGreen
//              if ms.Color == tcell.ColorGreen {
//                color = tcell.ColorBlue
//              }
//
//              if ms.Color == tcell.ColorBlue {
//                color = tcell.ColorGreen
//              }
//
//              return r.State{EffectExampleBoxState{
//                Color: color,
//              },
//              }
//            })
//          }
//        }
//      }()
//      return func() {
//        <-done
//      }
//    }, r.EffectDependencies{})
//
//    // var mouseEventHandler r.MouseEventHandler
//    mouseEventHandler := func(e *tcell.EventMouse) {
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
//        return r.State{EffectExampleBoxState{
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
func UseEffect(effect Effect, deps EffectDependencies) {
	hookFiberLock.Lock()

	var oldHook *hook
	if hookFiber != nil &&
		hookFiber.alternate != nil &&
		hookFiber.alternate.hooks != nil &&
		len(hookFiber.alternate.hooks) > hookIndex &&
		hookFiber.alternate.hooks[hookIndex] != nil {
		oldHook = hookFiber.alternate.hooks[hookIndex]
	}

	hasChanged := true

	if oldHook != nil {
		hasChanged = hasDepsChanged(oldHook.deps, deps)
	}

	h := &hook{
		tag:    hookTagEffect,
		effect: nil,
		cancel: nil,
	}

	if hasChanged {
		h.effect = effect
		h.deps = deps
	}

	if hookFiber != nil {
		hookFiber.hooks = append(hookFiber.hooks, h)
		hookIndex++
	}
	hookFiberLock.Unlock()
}

func hasDepsChanged(
	prevDeps,
	nextDeps EffectDependencies,
) (changed bool) {

	// TODO cleanup
	// if len(prevDeps) != len(nextDeps) {
	// 	changed = true
	// }
	// if len(prevDeps) == 0 {
	// 	changed = true
	// }
	// if len(nextDeps) == 0 {
	// 	changed = true
	// }

	// Check the slices have the same contents, in the same order
	for i, pd := range prevDeps {
		if nextDeps[i] != pd {
			changed = true
		}
	}

	return
}

func (r *retort) processEffects(f *fiber) {
	runEffects := false
	cancelEffects := false

	if f == nil {
		return
	}

	switch f.effect {
	case fiberEffectNothing:
		runEffects = true
	case fiberEffectPlacement:
		runEffects = true
	case fiberEffectUpdate:
		runEffects = true
	case fiberEffectDelete:
		cancelEffects = true
	}

	if f.hooks != nil {
		if cancelEffects {
			for _, h := range f.hooks {
				if h.tag != hookTagEffect ||
					h.effect == nil {
					continue
				}
				h.cancel()
			}
		}

		if runEffects {
			for _, h := range f.hooks {
				if h.tag != hookTagEffect ||
					h.effect == nil {
					continue
				}
				h.cancel = h.effect()
			}
		}
	}

	r.processEffects(f.child)
	r.processEffects(f.sibling)
}

// [ UseQuit ]--------------------------------------------------------------------------------------

// UseQuit returns a single function that when invoked
// will exit the application
func UseQuit() func() {
	return func() {
		close(quitChan)
	}
}

// [ UseScreen ]------------------------------------------------------------------------------------

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
