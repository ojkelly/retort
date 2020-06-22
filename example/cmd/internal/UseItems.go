package internal

import (
	"github.com/tjarratt/babble"
	"retort.dev/r"
)

type Item struct {
	Title string
	Date  string
	Body  string
}

type ItemsState struct {
	Items        []Item
	SelectedItem int
}

func UseItems(itemCount int) (itemsState ItemsState, setItemState func(selected int)) {
	babbler := babble.NewBabbler()

	initialItems := []Item{}

	for i := 1; i <= itemCount; i++ {
		item := Item{
			Title: babbler.Babble(),
			Date:  babbler.Babble(),
			Body:  loremIpsum,
		}
		initialItems = append(initialItems, item)
	}

	initialState := r.State{
		ItemsState{
			Items:        initialItems,
			SelectedItem: 0,
		},
	}

	state, setState := r.UseState(initialState)

	setItemState = func(selected int) {
		setState(func(s r.State) r.State {
			iState := s.GetState(
				ItemsState{},
			).(ItemsState)

			return r.State{
				ItemsState{
					Items:        iState.Items,
					SelectedItem: selected,
				},
			}
		},
		)

	}

	itemsState = state.GetState(
		ItemsState{},
	).(ItemsState)

	return itemsState, setItemState
}
