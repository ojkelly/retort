package r

import (
	"github.com/gdamore/tcell"
	"retort.dev/r/debug"
	"retort.dev/r/internal/quadtree"
)

type EdgeSizes struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

// BlockLayout is used by ScreenElements to determine the exact location to
// calculate/render from.
// It represents the concrete positioning information specific
// to the size of the terminal screen.
//
// You never set this directly, it's calculated via a component like
// Box. Which allows for more expressive objects, with padding, and margin.
//
// It is recalculated when the screen is resized.
//
// This layout information is also used to calculate which elements mouse events
// effect.
//
// You shouldn't use this except for a call to r.CreateScreenElement
type BlockLayout struct {
	X, Y, Rows, Columns int

	Padding, Border, Margin EdgeSizes

	// Grow, like flex-grow
	// TODO: better docs
	Grow int

	// ZIndex is the layer this Box is printed on.
	// Specifically, it determines the order of painting on the screen, with
	// higher numbers being painted later, and appearing on top.
	// This is also used to direct some events, where the highest zindex is used.
	ZIndex int

	// Order is set to control the display order of a group of children
	Order int

	// Fixed if the Rows, Columns are an explicit fixed size, else they're fluid
	FixedRows, FixedColumns bool

	// Valid is set when the BlockLayout has been initialised somewhere
	// if it's false, it means we've got a default
	Valid bool

	// Hide is set when this BlockLayout should not be rendered
	Hide bool
}

type BlockLayouts = []BlockLayout
type BlockLayoutWithProperties struct {
	BlockLayout BlockLayout
	Properties  Properties
}

type CalculateLayoutStage int

const (
	// Initial Pass
	// Calculate implicit or explicit absolute bounds
	CalculateLayoutStageInitial CalculateLayoutStage = iota

	// After this Blocks children have calculated their layouts
	// we recalculate this blocks layou
	CalculateLayoutStageWithChildren

	// Final Pass
	CalculateLayoutStageFinal
)

// CalculateLayout
//
// childrenBlockLayouts will be empty until at
// least CalculateLayoutStageWithChildren
//
// innerBlockLayout is the draw area for children blocks, and will
// be smaller due to padding or border effects
type CalculateLayout func(
	s tcell.Screen,
	stage CalculateLayoutStage,
	parentBlockLayout BlockLayout,
	children []BlockLayoutWithProperties,
) (
	blockLayout BlockLayout,
	innerBlockLayout BlockLayout,
	childrenBlockLayouts []BlockLayoutWithProperties,
)

// reconcileQuadTree updates the quadtree with our new layout, and provides
// the default box layout (from the parent) if none is available on the element
func (r *retort) reconcileQuadTree(f *fiber) {
	if f == nil {
		return
	}

	skip := false

	BlockLayout := f.Properties.GetOptionalProperty(
		BlockLayout{},
	).(BlockLayout)

	if BlockLayout.X == 0 &&
		BlockLayout.Y == 0 &&
		BlockLayout.Rows == 0 &&
		BlockLayout.Columns == 0 {
		skip = true
	}

	if !skip {
		r.quadtree.Insert(quadtree.Bounds{
			X:      BlockLayout.X,
			Y:      BlockLayout.Y,
			Width:  BlockLayout.Columns,
			Height: BlockLayout.Rows,

			// Store a pointer to our fiber for retrieval
			// We will need to cast this on the way out
			Value: f,
		})
	}

	r.reconcileQuadTree(f.child)
	r.reconcileQuadTree(f.sibling)
}

func (r *retort) calculateLayout(f *fiber) {
	if f == nil {
		return
	}

	if f.BlockLayout.Hide {
		debug.Log("f.BlockLayout.Hide")
		return
	}

	screen := UseScreen()

	if f.calculateLayout != nil {
		calcLayout := *f.calculateLayout

		blockLayout, innerBlockLayout, _ := calcLayout(
			screen,
			CalculateLayoutStageInitial,
			f.InnerBlockLayout,
			nil,
		)

		if blockLayout.Hide {
			debug.Log("blockLayout.Hide")
			return
		}

		f.BlockLayout = blockLayout

		f.InnerBlockLayout = innerBlockLayout

		r.calculateLayout(f.child)
		r.calculateLayout(f.sibling)

		children := f.ImmeditateChildren()

		if len(children) > 0 {
			childrenBlockLayoutsWithProperties := []BlockLayoutWithProperties{}

			for _, c := range children {
				if c != nil {
					cbl := BlockLayoutWithProperties{
						BlockLayout: c.BlockLayout,
						Properties:  c.Properties,
					}
					childrenBlockLayoutsWithProperties = append(childrenBlockLayoutsWithProperties, cbl)
				}
			}

			bl, _, childrenBlockLayouts := calcLayout(
				screen,
				CalculateLayoutStageWithChildren,
				f.InnerBlockLayout,
				childrenBlockLayoutsWithProperties,
			)
			if bl.Hide {
				debug.Log("bl.Hide")
				return
			}

			// Put the updated blockLayouts back onto the children
			for i, c := range children {
				if c == nil {
					continue
				}

				c.BlockLayout = childrenBlockLayouts[i].BlockLayout
				c.InnerBlockLayout = childrenBlockLayouts[i].BlockLayout
			}
		}

		f.Properties = ReplaceProps(f.Properties, f.BlockLayout)
	} else {
		// Pass on the BlockLayouts down the tree
		children := f.ImmeditateChildren()

		if len(children) > 0 {
			// Put the updated blockLayouts back onto the children
			for _, c := range children {
				if c == nil {
					continue
				}

				c.BlockLayout = f.InnerBlockLayout
				c.InnerBlockLayout = f.InnerBlockLayout
			}
		}
	}

	r.calculateLayout(f.child)
	r.calculateLayout(f.sibling)
}
