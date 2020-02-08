package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

const moveThreadMessage = `Error: missing arguments

/wrangler move thread [MESSAGE_ID] [CHANNEL_ID]
  Move a given message, along with the thread it belongs to, to a given channel
    - This can be on any channel in any team that you have joined
    - Obtain the message ID by running '/wrangler list messages' or via the 'Permalink' message dropdown option (it's the last part of the URL)
    - Obtain the channel ID by running '/wrangler list channels' or via the channel 'View Info' option
`

func getMoveThreadMessage() string {
	return codeBlock(moveThreadMessage)
}

func (p *Plugin) runMoveThreadCommand(args []string, extra *model.CommandArgs) (*model.CommandResponse, bool, error) {
	if len(args) < 2 {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, getMoveThreadMessage()), true, nil
	}
	postID := args[0]
	channelID := args[1]

	config := p.getConfiguration()

	originalChannel, appErr := p.API.GetChannel(extra.ChannelId)
	if appErr != nil {
		return nil, false, fmt.Errorf("unable to get channel with ID %s", extra.ChannelId)
	}
	switch originalChannel.Type {
	case model.CHANNEL_PRIVATE:
		if !config.MoveThreadFromPrivateChannelEnable {
			return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Wrangler is currently configured to not allow moving posts from private channels"), false, nil
		}
	case model.CHANNEL_DIRECT:
		if !config.MoveThreadFromDirectMessageChannelEnable {
			return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Wrangler is currently configured to not allow moving posts from direct message channels"), false, nil
		}
	}

	postListResponse, appErr := p.API.GetPostThread(postID)
	if appErr != nil {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, fmt.Sprintf("Error: unable to get post with ID %s; ensure this is correct", postID)), true, nil
	}
	postList := sortedPostsFromPostList(postListResponse)

	// Validation: let's check a few things before moving any posts.
	if len(postList) == 0 {
		return nil, false, fmt.Errorf("Sorting the post list response for post %s resulted in no posts", postID)
	}

	if config.MaxThreadCountMoveSizeInt() != 0 && config.MaxThreadCountMoveSizeInt() < len(postList) {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, fmt.Sprintf("Error: the thread is %d posts long, but the move thead command is configured to only move threads of up to %d posts", len(postList), config.MaxThreadCountMoveSizeInt())), true, nil
	}

	if postList[0].ChannelId != extra.ChannelId {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Error: the move command must be run from the channel containing the post"), true, nil
	}

	_, appErr = p.API.GetChannelMember(channelID, extra.UserId)
	if appErr != nil {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, fmt.Sprintf("Error: channel with ID %s doesn't exist or you are not a member", channelID)), true, nil
	}

	// We now know:
	// 1. The postID is valid.
	// 2. The channelID is valid and the user is a member of that channel.
	// 3. The command was run from the original channel with the post, so they
	//    are also a member of that channel.

	targetChannel, appErr := p.API.GetChannel(channelID)
	if appErr != nil {
		return nil, false, fmt.Errorf("unable to get channel with ID %s", channelID)
	}

	targetTeam, appErr := p.API.GetTeam(targetChannel.TeamId)
	if appErr != nil {
		return nil, false, fmt.Errorf("unable to get team with ID %s", targetChannel.TeamId)
	}

	// Cleanup is handled by simply deleting the root post. Any comments/replies
	// are automatically marked as deleted for us.
	cleanupID := postList[0].Id

	var newRootPost *model.Post

	for i, post := range postList {
		if i == 0 {
			cleanPost(post)
			post.ChannelId = channelID
			newRootPost, appErr = p.API.CreatePost(post)
			if appErr != nil {
				return nil, false, errors.Wrap(appErr, "unable to create new root post")
			}

			continue
		}

		cleanPost(post)
		post.ChannelId = channelID
		post.RootId = newRootPost.Id
		post.ParentId = newRootPost.Id
		_, appErr = p.API.CreatePost(post)
		if appErr != nil {
			return nil, false, errors.Wrap(appErr, "unable to create new post")
		}
	}

	_, appErr = p.API.CreatePost(&model.Post{
		UserId:    p.BotUserID,
		RootId:    newRootPost.Id,
		ParentId:  newRootPost.Id,
		ChannelId: channelID,
		Message:   "This thread was moved from another channel",
	})
	if appErr != nil {
		return nil, false, errors.Wrap(appErr, "unable to create new bot post")
	}

	appErr = p.API.DeletePost(cleanupID)
	if appErr != nil {
		return nil, false, errors.Wrap(appErr, "unable to delete post")
	}

	msg := fmt.Sprintf("A thread with %d posts has been moved [ team=%s, channel=%s ]", len(postList), targetTeam.Name, targetChannel.Name)

	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_IN_CHANNEL, msg), false, nil
}

func sortedPostsFromPostList(postList *model.PostList) []*model.Post {
	postList.UniqueOrder()
	postList.SortByCreateAt()
	posts := postList.ToSlice()

	var reversedPosts []*model.Post
	for i := range posts {
		reversedPosts = append(reversedPosts, posts[len(posts)-i-1])
	}

	return reversedPosts
}

func cleanPost(post *model.Post) {
	post.Id = ""
	post.CreateAt = 0
	post.UpdateAt = 0
	post.EditAt = 0
}
