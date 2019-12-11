package cache

import (
	"time"

	"github.com/munrocape/hn/hnclient"
	"retort.dev/r"
)

var StoriesContext = r.CreateContext(r.State{Stories{}})
var CommentContext = r.CreateContext(r.State{Comments{}})

// [ Stories ]------------------------------------------------------------------

type StoryItem struct {
	Story       *hnclient.Story
	Loading     bool
	Error       error
	LastUpdated time.Time
}
type Stories struct {
	Cache  map[int]StoryItem
	Update func(id int, item StoryItem)
}

func (c *Stories) Get(id int) StoryItem {
	story, ok := c.Cache[id]
	if !ok {
		return StoryItem{}
	}
	return story
}

// [ Comments ]-----------------------------------------------------------------

type CommentItem struct {
	Comment     *hnclient.Comment
	Loading     bool
	Error       error
	LastUpdated time.Time
}
type Comments struct {
	Cache  map[int]CommentItem
	Update func(id int, item CommentItem)
}

func Cache(p r.Properties) r.Element {
	children := p.GetProperty(
		r.Children{},
		"Cache requires r.Children",
	).(r.Children)

	s, storySetState := r.UseState(r.State{
		Stories{
			Cache: make(map[int]StoryItem),
		},
	})

	storyState := s.GetState(
		Stories{},
	).(Stories)

	storyState.Update = func(id int, item StoryItem) {
		storySetState(func(s r.State) r.State {
			storyCache := s.GetState(
				Stories{},
			).(Stories)
			storyCache.Cache[id] = item
			return r.State{storyCache}
		})
	}

	StoriesContext.Mount(r.State{storyState})

	// storySetState(func(s r.State) r.State {
	// 	return r.State{
	// 		Stories{
	// 			Cache:  make(map[int]StoryItem),
	// 			Update: updateStoryItem,
	// 		},
	// 	}
	// })

	// commentSetState := CommentContext.Mount()
	// updateCommentItem := func(id int, item CommentItem) {
	// 	commentSetState(func(s r.State) r.State {
	// 		s[id] = item
	// 		return r.State{s}
	// 	})
	// }

	// commentSetState(func(s r.State) r.State {
	// 	return r.State{
	// 		Comments{
	// 			Cache:  make(map[int]CommentItem),
	// 			Update: updateCommentItem,
	// 		},
	// 	}
	// })

	return r.CreateFragment(children)
}
