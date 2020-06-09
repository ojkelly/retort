package r

import (
	"sort"

	"github.com/gdamore/tcell"
	"retort.dev/debug"
)

// RenderToScreen is the callback passed to create a Screen Element
type RenderToScreen func(
	s tcell.Screen,
	blockLayout BlockLayout,
)

type DisplayCommand struct {
	RenderToScreen *RenderToScreen
	BlockLayout    BlockLayout
}

// DisplayList
// https://en.wikipedia.org/wiki/Display_list
type DisplayList []DisplayCommand
type DisplayListSortZIndex []DisplayCommand

func (dl DisplayListSortZIndex) Len() int { return len(dl) }
func (dl DisplayListSortZIndex) Less(i, j int) bool {
	return dl[i].BlockLayout.ZIndex < dl[j].BlockLayout.ZIndex
}
func (dl DisplayListSortZIndex) Swap(i, j int) { dl[i], dl[j] = dl[j], dl[i] }

// Sort a DisplayList for rendering to screen, with respect to ZIndexes
func (dl DisplayList) Sort() {
	sort.Sort(DisplayListSortZIndex(dl))
}

// commitRoot processes a tree root, and commits the results
// It's used to process updates for a fiber render, and is called when the
// main workloop has run out of tasks
func (r *retort) commitRoot() {
	screen := UseScreen()
	displayList := DisplayList{}

	// for _, deletion := range r.deletions {
	// 	displayList = append(displayList, r.processDisplayCommands(deletion)...)
	// }

	// w, h := screen.Size()
	debug.Log("Render: Calculate Layout")

	r.calculateLayout(r.wipRoot)

	// debug.Spew(r.wipRoot)

	debug.Log("Render: Create Display List")
	// Draw
	// TODO: conver this to a 2 step, first create a DisplayList (a list of commands for what to draw)
	// then optmise the list, by sorting by z-index, and removing commands that are occuluded
	// then run the commands sequentially
	displayList = append(displayList, r.processDisplayCommands(r.wipRoot)...)

	displayList.Sort()

	debug.Spew("DisplayList", displayList)

	debug.Log("Render: Paint")
	r.paint(displayList)

	screen.Show()

	// Update effects
	r.processEffects(r.wipRoot)

	// Update our quadtree for collisions
	r.quadtree.Clear()
	r.reconcileQuadTree(r.wipRoot)

	// Clean up and prepare for the next render pass
	r.currentRoot = r.wipRoot
	r.wipRoot = nil
	r.hasChangesToRender = false

	hookFiber = nil
	debug.Log("committed root")
}

// commitWork walks the tree and commits any fiber updates
func (r *retort) processDisplayCommands(f *fiber) (displayList DisplayList) {
	if f == nil {
		return
	}

	// debug.Log(fmt.Sprintf("processDisplayCommands address: %p", f))
	// debug.Spew(f)

	// TODO: collect all the renderToScreen paired with their zIndex
	// render all from lowest to highest index
	switch f.effect {
	case fiberEffectNothing:
	case fiberEffectPlacement:

		if f.renderToScreen == nil {
			break
		}

		// debug.Spew(fmt.Sprintf("f address %p", f), "render b", f.BlockLayout)
		// debug.Spew(f)

		displayCommand := DisplayCommand{
			RenderToScreen: f.renderToScreen,
			BlockLayout:    f.BlockLayout,
		}

		displayList = append(displayList, displayCommand)
	case fiberEffectUpdate:
		// cancelEffects(f)

		if f.renderToScreen == nil {
			break
		}
		// debug.Log(fmt.Sprintf("f address %p", f), "render update b", f.BlockLayout)

		displayCommand := DisplayCommand{
			RenderToScreen: f.renderToScreen,
			BlockLayout:    f.BlockLayout,
		}

		displayList = append(displayList, displayCommand)

	case fiberEffectDelete:
	}

	if f.child != nil {
		displayList = append(displayList, r.processDisplayCommands(f.child)...)
	}

	if f.sibling != nil {
		displayList = append(displayList, r.processDisplayCommands(f.sibling)...)
	}

	f.dirty = false

	return
}

func (r *retort) commitDeletion(f *fiber) {
	// if (fiber.dom) {
	//   domParent.removeChild(fiber.dom);
	// } else {
	//   commitDeletion(fiber.child, domParent);
	// }
}

// paint the DisplayList to screen
func (r *retort) paint(displayList DisplayList) {
	screen := UseScreen()

	for _, command := range displayList {
		render := *command.RenderToScreen

		render(
			screen,
			command.BlockLayout,
		)
	}

	screen.Show()
}
