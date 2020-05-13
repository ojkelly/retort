package r

type fiberEffect int

const (
	fiberEffectNothing fiberEffect = iota
	fiberEffectUpdate
	fiberEffectPlacement
	fiberEffectDelete
)

type componentType int

const (
	nothingComponent componentType = iota
	elementComponent
	fragmentComponent
	screenComponent
)

type fiber struct {
	// dirty when there are changes to commit
	dirty bool

	componentType componentType
	component     Component
	Properties    Properties
	parent        *fiber
	sibling       *fiber
	child         *fiber
	alternate     *fiber
	effect        fiberEffect
	hooks         []*hook

	// Layout Information
	renderToScreen *RenderToScreen

	calculateLayout *CalculateLayout

	// this BlockLayout is used internally, mainly to route events
	// Its different to the BlockLayout that may be passed around in props
	BlockLayout BlockLayout

	// InnerBlockLayout is passed to children of this fiber
	InnerBlockLayout BlockLayout

	// focus bool
}

func (f *fiber) Parent() *fiber {
	return f.parent
}

func (f *fiber) Sibling() *fiber {
	return f.sibling
}

func (f *fiber) Child() *fiber {
	return f.child
}

// Clone safely makes a copy of a hook for use with fiber updates
func (f *fiber) Clone() (newFiber *fiber) {
	// Parent, sibling, and alternate are not cloned
	// as doing so will recurse forever
	newFiber = &fiber{
		componentType:    f.componentType,
		component:        f.component,
		Properties:       f.Properties,
		effect:           f.effect,
		alternate:        f.alternate,
		hooks:            f.hooks,
		BlockLayout:      f.BlockLayout,
		InnerBlockLayout: f.InnerBlockLayout,
	}

	if f.renderToScreen != nil {
		render := *f.renderToScreen
		newFiber.renderToScreen = &render
	}

	if f.calculateLayout != nil {
		calcLayout := *f.calculateLayout
		newFiber.calculateLayout = &calcLayout
	}

	if f.child != nil {
		newFiber.child = f.child.Clone()
	}
	if f.sibling != nil {
		newFiber.sibling = f.sibling.Clone()
	}

	return
}

// cloneElements safely makes a copy of elements of a fiber for use with updates
func cloneElements(fibers []*fiber) (cloned []*fiber) {
	cloned = []*fiber{}

	for _, f := range fibers {
		cloned = append(cloned, f.Clone())
	}

	return
}
