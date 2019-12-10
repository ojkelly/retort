package r

import (
	"fmt"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/gdamore/tcell/encoding"
	"retort.dev/debug"
	"retort.dev/r/internal/quadtree"
)

type retort struct {
	root *fiber

	nextUnitOfWork *fiber
	currentRoot    *fiber
	wipRoot        *fiber
	wipFiber       *fiber
	deletions      []*fiber

	hasChangesToRender bool
	hasNewState        bool

	rootBoxLayout BoxLayout
	quadtree      quadtree.Quadtree
	config        RetortConfiguration
}

// RetortConfiguration allows you to enable features your app
// may want to use
type RetortConfiguration struct {
	// UseSimulationScreen to output to a simulated screen
	// this is useful for automated testing
	UseSimulationScreen bool

	// UseDebugger to show a debug overlay with output from
	// the retort.dev/debug#Log function
	UseDebugger bool

	// DisableMouse to prevent Mouse Events from being created
	DisableMouse bool
}

var setStateChan chan ActionCreator

var quitChan chan struct{}

var c *RetortConfiguration = &RetortConfiguration{}

// Retort is called with your root Component and any optional
// configuration to begin running retort.
//
//  func Example_app() {
//    // Call the main function on retort to start the app,
//    // when you call this, retort will take over the screen.
//    r.Retort(
//      // Root Element
//      r.CreateElement(
//        example.ClickableBox,
//        r.Properties{
//          component.BoxProps{
//            Width:  100, // Make the root element fill the screen
//            Height: 100, // Make the root element fill the screen
//            Border: component.Border{
//              Style:      component.BorderStyleSingle,
//              Foreground: tcell.ColorWhite,
//            },
//          },
//        },
//        r.Children{
//          // First Child
//          r.CreateElement(
//            example.ClickableBox,
//            r.Properties{
//              component.BoxProps{
//                Border: component.Border{
//                  Style:      component.BorderStyleSingle,
//                  Foreground: tcell.ColorWhite,
//                },
//              },
//            },
//            nil, // Pass nil as the third argument if there are no children
//          ),
//          // Second Child
//          r.CreateElement(
//            example.ClickableBox,
//            r.Properties{
//              component.BoxProps{
//                Border: component.Border{
//                  Style:      component.BorderStyleSingle,
//                  Foreground: tcell.ColorWhite,
//                },
//              },
//            },
//            nil,
//          ),
//        },
//      ),
//      // Pass in optional configuration
//      r.RetortConfiguration{}
//    )
//  }
func Retort(root Element, config RetortConfiguration) {
	r := &retort{
		root:   root,
		config: config,
		quadtree: quadtree.Quadtree{
			MaxObjects: 2000,
			MaxLevels:  1000,
			Level:      0,
		},
	}

	c = &config

	quitChan = make(chan struct{})

	setStateChan = make(chan ActionCreator, 2000)

	screen := UseScreen()
	defer screen.Fini()

	encoding.Register()

	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	r.root = root

	w, h := screen.Size()

	r.parseRetortConfiguration()

	r.rootBoxLayout = BoxLayout{
		X:       0,
		Y:       0,
		Columns: w + 1, // +1 to account for zero-indexing
		Rows:    h + 1, // +1 to account for zero-indexing
	}
	r.quadtree.Bounds.Width = w
	r.quadtree.Bounds.Height = h

	root.Properties = append(root.Properties, r.rootBoxLayout)

	r.wipRoot = &fiber{
		componentType: nothingComponent,
		Properties:    Properties{Children{root}},
		alternate:     r.currentRoot,
	}
	r.nextUnitOfWork = r.wipRoot
	r.currentRoot = r.wipRoot.Clone()
	r.hasChangesToRender = true

	var frame int
	var deadline time.Time

	// event handling
	go r.handleEvents()

	// work loop
	go func() {
		deadline = time.Now().Add(14 * time.Millisecond)
		workTick := time.NewTicker(1 * time.Nanosecond)
		frameTick := time.NewTicker(16 * time.Millisecond)
		shouldYield := false

		var droppedFrames int

	workloop:
		for {
			select {
			case <-quitChan:
				workTick.Stop()
				break workloop
			// workloop
			case action := <-setStateChan:
				action.addToQueue()
				r.hasNewState = true
			case <-frameTick.C:
				if r.hasNewState {
					r.updateTree()
				}
				if r.hasChangesToRender {
					// If there's still setStates to add to the queue, give them a chance
					// to be added
					if len(setStateChan) > 0 && droppedFrames == 0 {
						droppedFrames++
						continue
					}
					workTick = time.NewTicker(1 * time.Nanosecond)
					droppedFrames = 0
				}
				deadline = time.Now().Add(14 * time.Millisecond)

			case <-workTick.C:
				if !r.hasChangesToRender {
					// While we have work to do, this case is run very frequently
					// But when we have no work to do it can consume considerable CPU time
					// So we only start this ticker when we actuall have work to do,
					// and we stop it the rest of the time.
					// We use a frame tick to ensure at least once every 16ms (60fps)
					// we are checking if we need to do work
					workTick.Stop()
				}
				if r.nextUnitOfWork != nil && !shouldYield {
					start := time.Now()
					r.nextUnitOfWork = r.performWork(r.nextUnitOfWork)
					debug.Spew("performWork: ", time.Since(start))

					// yield with time to render
					if time.Since(deadline) > 100*time.Nanosecond {
						shouldYield = true
					}
				}

				if r.nextUnitOfWork == nil && r.wipRoot != nil {
					start := time.Now()
					r.commitRoot()
					debug.Spew("commitRoot: ", time.Since(start))
					shouldYield = false
				}

				if time.Since(deadline) > 0 {
					shouldYield = false
					frame++
				}
			}
		}
	}()

	// Wait until quit
	<-quitChan
	screen.Clear()
	screen.Fini()
}

func (r *retort) parseRetortConfiguration() {
	screen := UseScreen()

	if !r.config.DisableMouse {
		screen.EnableMouse()
	}
}

// ForceRender can be called at any point to ask
// retort to start a whole new update
func (r *retort) ForceRender() {
	r.updateTree()
}

// [ Working ]------------------------------------------------------------------

func (r *retort) updateTree() {
	r.wipRoot = r.currentRoot
	r.wipRoot.alternate = r.currentRoot.Clone()
	r.wipRoot.dirty = true

	r.nextUnitOfWork = r.wipRoot
	r.wipFiber = nil
	r.deletions = nil
	r.hasChangesToRender = true
	r.hasNewState = false
}

func (r *retort) performWork(f *fiber) *fiber {
	r.updateComponent(f)

	if f.child != nil {
		return f.child
	}

	nextFiber := f

	for nextFiber != nil {
		if nextFiber.sibling != nil {
			return nextFiber.sibling
		}
		nextFiber = nextFiber.parent
	}

	return nil
}

var hookFiber *fiber
var hookFiberLock = &sync.Mutex{}

// [ Components ]---------------------------------------------------------------

func (r *retort) updateComponent(f *fiber) {
	hookFiberLock.Lock()
	hookFiber = f
	hookFiberLock.Unlock()

	switch f.componentType {
	case nothingComponent:
		r.updateNothingComponent(f)
	case elementComponent:
		r.updateElementComponent(f)
	case fragmentComponent:
		r.updateFragmentComponent(f)
	case screenComponent:
		r.updateScreenComponent(f)
	}

	// debug.Spew("updateComponent", f)
	hookFiberLock.Lock()
	hookFiber = nil
	hookFiberLock.Unlock()
}

func (r *retort) updateElementComponent(f *fiber) {
	if f == nil || f.componentType != elementComponent {
		return
	}

	if f.component == nil || f.Properties == nil {
		return
	}
	r.wipFiber = f

	hookFiberLock.Lock()
	hookIndex = 0
	hookFiberLock.Unlock()

	r.wipFiber.hooks = nil

	children := f.component(f.Properties)
	// debug.Spew("updateElementComponent children", children)
	r.reconcileChildren(f, []*fiber{children})
}

func (r *retort) updateFragmentComponent(f *fiber) {
	if f == nil || f.componentType != fragmentComponent || f.Properties == nil {
		return
	}

	r.wipFiber = f

	children := f.Properties.GetProperty(
		Children{},
		"Fragment requires r.Children",
	).(Children)

	r.reconcileChildren(f, children)
}

func (r *retort) updateNothingComponent(f *fiber) {
	if f == nil || f.componentType != nothingComponent {
		return
	}

	r.wipFiber = f

	children := f.Properties.GetOptionalProperty(
		Children{},
	).(Children)
	r.reconcileChildren(f, children)
}

func (r *retort) updateScreenComponent(f *fiber) {
	if f == nil || f.componentType != screenComponent {
		return
	}

	r.wipFiber = f

	children := f.Properties.GetOptionalProperty(
		Children{},
	).(Children)

	r.reconcileChildren(f, children)
	// debug.Spew("updateScreenComponent", f)

}

// [ Children ]-----------------------------------------------------------------

func (r *retort) reconcileChildren(f *fiber, elements []*fiber) {
	index := 0

	f.dirty = false

	var oldFiber *fiber
	var boxLayout BoxLayout
	if r.wipFiber != nil && r.wipFiber.alternate != nil {
		oldFiber = r.wipFiber.alternate.child
	}

	boxLayout = f.Properties.GetOptionalProperty(BoxLayout{}).(BoxLayout)

	var prevSibling *fiber

	// Add newly generated child elements, as children to this fiber
	for index < len(elements) || oldFiber != nil {
		var element *fiber
		if len(elements) != 0 {
			element = elements[index]
		}

		var newFiber *fiber

		sameType := false

		if oldFiber != nil && element != nil &&
			reflect.TypeOf(element.component) == reflect.TypeOf(oldFiber.component) {
			sameType = true
		}

		if sameType { // Update
			f.dirty = true
			newFiber = &fiber{
				dirty:          true,
				componentType:  oldFiber.componentType,
				component:      oldFiber.component,
				Properties:     AddPropsIfNone(element.Properties, boxLayout),
				parent:         f,
				alternate:      oldFiber,
				effect:         fiberEffectUpdate,
				renderToScreen: element.renderToScreen,
			}
		}

		if element != nil && !sameType { // New Placement
			f.dirty = true
			newFiber = &fiber{
				dirty:          true,
				componentType:  element.componentType,
				component:      element.component,
				Properties:     AddPropsIfNone(element.Properties, boxLayout),
				parent:         f,
				alternate:      nil,
				effect:         fiberEffectPlacement,
				renderToScreen: element.renderToScreen,
			}
		}

		if oldFiber != nil && !sameType { // Delete
			oldFiber.effect = fiberEffectDelete
			r.deletions = append(r.deletions, oldFiber)
		}

		if oldFiber != nil { // nothing to update
			oldFiber = oldFiber.sibling
		}

		if index == 0 {
			f.dirty = true
			f.child = newFiber
		} else if element != nil {
			f.dirty = true
			prevSibling.sibling = newFiber
		}

		prevSibling = newFiber
		index++
	}
}
