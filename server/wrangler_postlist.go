package main

import (
	"github.com/mattermost/mattermost-server/v5/model"
)

// WranglerPostList provides a list of posts along with metadata about those
// posts.
type WranglerPostList struct {
	Posts                []*model.Post
	ThreadUserIDs        []string
	EarlistPostTimestamp int64
	LatestPostTimestamp  int64
	ContainsAttachments  bool
}

// NumPosts returns the number of posts in a post list.
func (wpl *WranglerPostList) NumPosts() int {
	return len(wpl.Posts)
}

// RootPost returns the root post in a post list.
func (wpl *WranglerPostList) RootPost() *model.Post {
	if wpl.NumPosts() < 1 {
		return nil
	}

	return wpl.Posts[0]
}

func buildWranglerPostList(postList *model.PostList) *WranglerPostList {
	wpl := &WranglerPostList{}

	postList.UniqueOrder()
	postList.SortByCreateAt()
	posts := postList.ToSlice()

	if len(posts) == 0 {
		// Something was sorted wrong or an empty PostList was provided.
		return wpl
	}

	// A separate ID key map to ensure no duplicates.
	idKeys := make(map[string]bool)

	for i := range posts {
		p := posts[len(posts)-i-1]

		// Add UserID to metadata if it's new.
		if _, ok := idKeys[p.UserId]; !ok {
			idKeys[p.UserId] = true
			wpl.ThreadUserIDs = append(wpl.ThreadUserIDs, p.UserId)
		}

		// Mark postlist as containing attachments if post has attachment(s).
		if !wpl.ContainsAttachments && len(p.Attachments()) != 0 {
			wpl.ContainsAttachments = true
		}

		wpl.Posts = append(wpl.Posts, p)
	}

	// Set metadata for earliest and latest posts
	wpl.EarlistPostTimestamp = wpl.RootPost().CreateAt
	wpl.LatestPostTimestamp = wpl.Posts[wpl.NumPosts()-1].CreateAt

	return wpl
}
