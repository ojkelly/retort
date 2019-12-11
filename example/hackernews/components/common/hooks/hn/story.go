package hn

import (
	"time"

	"github.com/munrocape/hn/hnclient"
	"retort.dev/example/hackernews/components/cache"
	"retort.dev/r"
)

type CurrentStoryState struct {
	Data    hnclient.Story
	Loading bool
	Error   error
}

// UseCurrentStory
func UseCurrentStory(
	storiesContext *r.Context,
	commentsContext *r.Context,
) (
	CurrentStoryState,
	r.SetState,
) {
	s, setState := r.UseState(r.State{
		CurrentStoryState{},
	})
	state := s.GetState(
		CurrentStoryState{},
	).(CurrentStoryState)

	return state, setState
}

func UseStory(id int) (
	story *hnclient.Story,
	loading bool,
	err error,
) {
	c := UseHackerNews()

	// Get our storyCache from the passed in Context
	sc := r.UseContext(cache.StoriesContext)

	storyCache := sc.GetState(
		cache.Stories{},
	).(cache.Stories)

	if storyCache.Update == nil {
		return
	}
	r.UseEffect(func() r.EffectCancel {

		var story *hnclient.Story
		var needToFetch bool

		// Check if story is in the cache
		cachedStory, ok := storyCache.Cache[id]
		if !ok {
			needToFetch = true
		}

		// if it's in the cache, check how fresh it is
		if ok &&
			cachedStory.Story != nil &&
			time.Since(cachedStory.LastUpdated) > CacheTimeout {
			needToFetch = true
		}

		if cachedStory.Loading {
			needToFetch = false
		}

		// If we need to update, go fetch it
		if needToFetch {
			item := cache.StoryItem{
				Story:       nil,
				Loading:     true,
				LastUpdated: time.Now(),
			}
			storyCache.Update(id, item)
			rawStory, err := c.GetStory(id)
			story = &rawStory

			// Update the storyCache with our hydrated story
			loadedItem := cache.StoryItem{
				Story:       story,
				Loading:     false,
				Error:       err,
				LastUpdated: time.Now(),
			}
			storyCache.Update(id, loadedItem)
		}

		return func() {}
	}, r.EffectDependencies{id})

	cachedStory, ok := storyCache.Cache[id]

	if !ok {
		return nil, true, nil
	}

	// TODO: return func to make this story the selected one
	return cachedStory.Story, cachedStory.Loading, nil
}
