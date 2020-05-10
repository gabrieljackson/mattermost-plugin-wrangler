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
	case model.CHANNEL_GROUP:
		if !config.MoveThreadFromGroupMessageChannelEnable {
			return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Wrangler is currently configured to not allow moving posts from group message channels"), false, nil
		}
	}

	postListResponse, appErr := p.API.GetPostThread(postID)
	if appErr != nil {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, fmt.Sprintf("Error: unable to get post with ID %s; ensure this is correct", postID)), true, nil
	}
	wpl := buildWranglerPostList(postListResponse)

	// Validation: let's check a few things before moving any posts.
	if wpl.NumPosts() == 0 {
		return nil, false, fmt.Errorf("Sorting the post list response for post %s resulted in no posts", postID)
	}

	if config.MaxThreadCountMoveSizeInt() != 0 && config.MaxThreadCountMoveSizeInt() < wpl.NumPosts() {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, fmt.Sprintf("Error: the thread is %d posts long, but the move thead command is configured to only move threads of up to %d posts", wpl.NumPosts(), config.MaxThreadCountMoveSizeInt())), true, nil
	}

	if wpl.RootPost().ChannelId != extra.ChannelId {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Error: the 'move thread' command must be run from the channel containing the post"), true, nil
	}

	_, appErr = p.API.GetChannelMember(channelID, extra.UserId)
	if appErr != nil {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, fmt.Sprintf("Error: channel with ID %s doesn't exist or you are not a member", channelID)), true, nil
	}

	targetChannel, appErr := p.API.GetChannel(channelID)
	if appErr != nil {
		return nil, false, fmt.Errorf("unable to get channel with ID %s", channelID)
	}
	if !config.MoveThreadToAnotherTeamEnable && targetChannel.TeamId != originalChannel.TeamId {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Wrangler is currently configured to not allow moving messages to different teams"), false, nil
	}

	if extra.RootId == wpl.RootPost().Id || extra.ParentId == wpl.RootPost().Id {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Error: the 'move thread' command cannot be run from inside the thread being moved; please run directly in the channel containing the thread you wish to move"), true, nil
	}

	// We now know:
	// 1. The postID is valid.
	// 2. The channelID is valid and the user is a member of that channel.
	// 3. The command was run from the original channel with the post, so they
	//    are also a member of that channel.

	targetTeam, appErr := p.API.GetTeam(targetChannel.TeamId)
	if appErr != nil {
		return nil, false, fmt.Errorf("unable to get team with ID %s", targetChannel.TeamId)
	}

	// Cleanup is handled by simply deleting the root post. Any comments/replies
	// are automatically marked as deleted for us.
	cleanupID := wpl.RootPost().Id

	var newRootPost *model.Post
	var originalMessageSummary string

	// Begin creating the new thread.
	p.API.LogInfo("Wrangler is moving a thread",
		"user_id", extra.UserId,
		"original_post_id", cleanupID,
		"original_channel_id", originalChannel.Id,
	)

	if wpl.ContainsFileAttachments() {
		// The thread contains at least one attachment. To properly move the
		// thread, the files will have to be re-uploaded. This is completed
		// before any messages are moved.
		// TODO: check number of files that need to be re-uploaded or file size?
		p.API.LogInfo("Wrangler is re-uploading file attachments",
			"file_count", wpl.FileAttachmentCount,
		)

		for _, post := range wpl.Posts {
			var newFileIDs []string
			var fileBytes []byte
			var oldFileInfo, newFileInfo *model.FileInfo
			for _, fileID := range post.FileIds {
				oldFileInfo, appErr = p.API.GetFileInfo(fileID)
				if appErr != nil {
					return nil, false, errors.Wrap(appErr, "unable to lookup file info to re-upload")
				}
				fileBytes, appErr = p.API.GetFile(fileID)
				if appErr != nil {
					return nil, false, errors.Wrap(appErr, "unable to get file bytes to re-upload")
				}
				newFileInfo, appErr = p.API.UploadFile(fileBytes, targetChannel.Id, oldFileInfo.Name)
				if appErr != nil {
					return nil, false, errors.Wrap(appErr, "unable to re-upload file")
				}

				newFileIDs = append(newFileIDs, newFileInfo.Id)
			}

			post.FileIds = newFileIDs
		}
	}

	for i, post := range wpl.Posts {
		if i == 0 {
			cleanPost(post)
			post.ChannelId = channelID
			newRootPost, appErr = p.API.CreatePost(post)
			if appErr != nil {
				return nil, false, errors.Wrap(appErr, "unable to create new root post")
			}
			originalMessageSummary = cleanAndTrimMessage(post.Message, 500)

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

	p.API.LogInfo("Wrangler thread move complete",
		"user_id", extra.UserId,
		"new_post_id", newRootPost.Id,
		"new_channel_id", channelID,
	)

	newPostLink := makePostLink(*p.API.GetConfig().ServiceSettings.SiteURL, newRootPost.Id)
	if extra.UserId != wpl.RootPost().UserId {
		// The wrangled thread was not started by the user running the command.
		// Send a DM to the user who created the root message to let them know.
		err := p.postMoveThreadBotDM(wpl.RootPost().UserId, newPostLink)
		if err != nil {
			p.API.LogError("Unable to send move-thread DM to user",
				"error", err.Error(),
				"user_id", wpl.RootPost().UserId,
			)
		}
	}

	msg := fmt.Sprintf("A thread has been moved: %s\n", newPostLink)
	msg += fmt.Sprintf(
		"\n| Team | Channel | Messages |\n| -- | -- | -- |\n| %s | %s | %d |\n\n",
		targetTeam.Name, targetChannel.Name, wpl.NumPosts(),
	)
	msg += fmt.Sprintf("Original Thread Root Message:\n%s\n", quoteBlock(originalMessageSummary))

	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_IN_CHANNEL, msg), false, nil
}

func (p *Plugin) postMoveThreadBotDM(userID, newPostLink string) error {
	return p.PostBotDM(userID, fmt.Sprintf(
		"Someone wrangled a thread you started to a new channel for you: %s", newPostLink,
	))
}
