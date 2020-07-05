package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

const moveThreadUsage = `/wrangler move thread [MESSAGE_ID] [CHANNEL_ID]
  Move a given message, along with the thread it belongs to, to a given channel
    - This can be on any channel in any team that you have joined
    - Obtain the message ID by running '/wrangler list messages' or via the 'Permalink' message dropdown option (it's the last part of the URL)
    - Obtain the channel ID by running '/wrangler list channels' or via the channel 'View Info' option`

func getMoveThreadMessage() string {
	return codeBlock(fmt.Sprintf("`Error: missing arguments\n\n%s", moveThreadUsage))
}

func (p *Plugin) runMoveThreadCommand(args []string, extra *model.CommandArgs) (*model.CommandResponse, bool, error) {
	if len(args) < 2 {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, getMoveThreadMessage()), true, nil
	}
	postID := args[0]
	channelID := args[1]

	postListResponse, appErr := p.API.GetPostThread(postID)
	if appErr != nil {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, fmt.Sprintf("Error: unable to get post with ID %s; ensure this is correct", postID)), true, nil
	}
	wpl := buildWranglerPostList(postListResponse)

	originalChannel, appErr := p.API.GetChannel(extra.ChannelId)
	if appErr != nil {
		return nil, false, fmt.Errorf("unable to get channel with ID %s", extra.ChannelId)
	}
	_, appErr = p.API.GetChannelMember(channelID, extra.UserId)
	if appErr != nil {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, fmt.Sprintf("Error: channel with ID %s doesn't exist or you are not a member", channelID)), true, nil
	}
	targetChannel, appErr := p.API.GetChannel(channelID)
	if appErr != nil {
		return nil, false, fmt.Errorf("unable to get channel with ID %s", channelID)
	}

	response, userErr, err := p.validateMoveOrCopy(wpl, originalChannel, targetChannel, extra)
	if response != nil || err != nil {
		return response, userErr, err
	}

	targetTeam, appErr := p.API.GetTeam(targetChannel.TeamId)
	if appErr != nil {
		return nil, false, fmt.Errorf("unable to get team with ID %s", targetChannel.TeamId)
	}

	// Begin creating the new thread.
	p.API.LogInfo("Wrangler is moving a thread",
		"user_id", extra.UserId,
		"original_post_id", wpl.RootPost().Id,
		"original_channel_id", originalChannel.Id,
	)

	// To simulate the move, we first copy the original messages(s) to the
	// new channel and later delete the original messages(s).
	newRootPost, err := p.copyWranglerPostlist(wpl, targetChannel)
	if err != nil {
		return nil, false, err
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

	// Cleanup is handled by simply deleting the root post. Any comments/replies
	// are automatically marked as deleted for us.
	appErr = p.API.DeletePost(wpl.RootPost().Id)
	if appErr != nil {
		return nil, false, errors.Wrap(appErr, "unable to delete post")
	}

	p.API.LogInfo("Wrangler thread move complete",
		"user_id", extra.UserId,
		"new_post_id", newRootPost.Id,
		"new_channel_id", channelID,
	)

	newPostLink := makePostLink(*p.API.GetConfig().ServiceSettings.SiteURL, targetTeam.Name, newRootPost.Id)
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

	originalMessageSummary := cleanAndTrimMessage(wpl.RootPost().Message, 500)

	msg := fmt.Sprintf("A thread has been moved: %s\n", newPostLink)
	msg += fmt.Sprintf(
		"\n| Team | Channel | Messages |\n| -- | -- | -- |\n| %s | %s | %d |\n\n",
		targetTeam.DisplayName, targetChannel.DisplayName, wpl.NumPosts(),
	)
	msg += fmt.Sprintf("Original Thread Root Message:\n%s\n", quoteBlock(originalMessageSummary))

	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_IN_CHANNEL, msg), false, nil
}

func (p *Plugin) postMoveThreadBotDM(userID, newPostLink string) error {
	return p.PostBotDM(userID, fmt.Sprintf(
		"Someone wrangled a thread you started to a new channel for you: %s", newPostLink,
	))
}
