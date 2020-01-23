package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/model"
	"github.com/pkg/errors"
)

const moveThreadMessage = `Error: missing arguments

Usage: |/wrangler move thread [MESSAGE_ID] [CHANNEL_ID]|

 * Obtain the message ID via the |Permalink| message dropdown option. It's the last part of the URL.
 * Obtain the channel ID via the Channel |View Info| option or by running |/wrangler list|.
`

func getMoveThreadMessage() string {
	return strings.Replace(moveThreadMessage, "|", "`", -1)
}

func (p *Plugin) runMoveThreadCommand(args []string, extra *model.CommandArgs) (*model.CommandResponse, bool, error) {
	if len(args) < 2 {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, getMoveThreadMessage()), true, nil
	}
	postID := args[0]
	channelID := args[1]

	// Validation: let's check a few things before moving any posts.
	postListResponse, appErr := p.API.GetPostThread(postID)
	if appErr != nil {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, fmt.Sprintf("Error: unable to get post with ID %s; ensure this is correct", postID)), true, nil
	}
	postListResponse.UniqueOrder()
	postListResponse.SortByCreateAt()
	postList := postListResponse.ToSlice()

	if len(postList) == 0 {
		return nil, false, fmt.Errorf("Sorting the post list response for post %s resulted in no posts", postID)
	}

	config := p.getConfiguration()
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

	var finalList []*model.Post
	var cleanupIDs []string
	for i := range postList {
		finalList = append(finalList, postList[len(postList)-i-1])
		cleanupIDs = append(cleanupIDs, postList[i].Id)
	}

	var newRootPost *model.Post

	for i, post := range finalList {
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

	for _, id := range cleanupIDs {
		appErr = p.API.DeletePost(id)
		if appErr != nil {
			return nil, false, errors.Wrap(appErr, "unable to delete post")
		}
	}

	msg := fmt.Sprintf("A thread with %d posts has been moved to %s:%s", len(finalList), targetTeam.Name, targetChannel.Name)

	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_IN_CHANNEL, msg), false, nil
}

func cleanPost(post *model.Post) {
	post.Id = ""
	post.CreateAt = 0
	post.UpdateAt = 0
	post.EditAt = 0
}
