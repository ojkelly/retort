package hn

import (
	"time"

	"retort.dev/example/hackernews/components/cache"
	"retort.dev/r"
)

var CacheTimeout time.Duration = time.Minute

type TopStoriesState struct {
	Data    []int
	Loading bool
	Error   error
}

var checked bool

// UseTopStories
func UseTopStories() TopStoriesState {
	c := UseHackerNews()

	sc := r.UseContext(cache.StoriesContext)
	storyCache := sc.GetState(
		cache.Stories{},
	).(cache.Stories)

	s, setState := r.UseState(r.State{
		TopStoriesState{Loading: true},
	})
	state := s.GetState(
		TopStoriesState{Loading: true},
	).(TopStoriesState)

	// debug.Spew(state)
	// Update list of Top Stories
	r.UseEffect(func() r.EffectCancel {
		if storyCache.Update == nil {
			return func() {}
		}

		if checked {
			return func() {}
		}

		checked = true

		topStories, err := c.GetTopStories(20)

		if err != nil {
			// debug.Spew("topStories err", err)
			setState(func(s r.State) r.State {
				return r.State{TopStoriesState{
					Loading: false,
					Error:   err,
				}}
			})
		} else {
			// debug.Spew("topStories", topStories)
			setState(func(s r.State) r.State {
				return r.State{TopStoriesState{
					Data:    topStories,
					Loading: false,
					Error:   nil,
				}}
			})
		}

		for _, id := range topStories {
			item := cache.StoryItem{
				Story:       nil,
				Loading:     false,
				LastUpdated: time.Now(),
			}
			storyCache.Update(id, item)
		}

		return func() {}
	}, r.EffectDependencies{storyCache.Update})

	// storiesContext := r.UseContext(cache.StoriesContext)

	// // Hydrate stories into cache
	// r.UseEffect(func() r.EffectCancel {
	// 	HydrateStories(state.Data, storiesContext)

	// 	return func() {}
	// }, r.EffectDependencies{state.Data})

	return state
}
