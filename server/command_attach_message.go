package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

const attachMessageCommand = `Error: missing arguments

/wrangler attach message [MESSAGE_ID_TO_BE_ATTACHED] [ROOT_MESSAGE_ID]
	Attach a given message to a thread in the same channel
	  - Obtain the message IDs by running '/wrangler list messages' or via the 'Permalink' message dropdown option (it's the last part of the URL)
`

func getAttachMessageCommand() string {
	return codeBlock(attachMessageCommand)
}

func (p *Plugin) runAttachMessageCommand(args []string, extra *model.CommandArgs) (*model.CommandResponse, bool, error) {
	if len(args) < 2 {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, getAttachMessageCommand()), true, nil
	}
	postToBeAttachedID := args[0]
	postToAttachToID := args[1]

	postToBeAttached, appErr := p.API.GetPost(postToBeAttachedID)
	if appErr != nil {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, fmt.Sprintf("Error: unable to get message with ID %s; ensure this is correct", postToBeAttachedID)), true, nil
	}
	postToAttachTo, appErr := p.API.GetPost(postToAttachToID)
	if appErr != nil {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, fmt.Sprintf("Error: unable to get message with ID %s; ensure this is correct", postToAttachToID)), true, nil
	}

	if postToBeAttached.ChannelId != extra.ChannelId {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Error: the attach command must be run from the channel containing the messages"), true, nil
	}
	if postToAttachTo.ChannelId != extra.ChannelId {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Error: unable to attach message to a thread in another channel"), true, nil
	}
	if len(postToBeAttached.RootId) != 0 || len(postToBeAttached.ParentId) != 0 {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Error: the message to be attached is already part of a thread"), true, nil
	}

	// We now know:
	// 1. The post IDs are valid.
	// 2. The post to be attached is not part of a thread already.
	// 3. The posts are in the same channel.
	// 4. The command was run from the original channel with the posts, so they
	//    are also a member of that channel.

	newRootID := postToAttachTo.Id
	if len(postToAttachTo.RootId) != 0 {
		newRootID = postToAttachTo.RootId
	}
	cleanupID := postToBeAttached.Id

	// Begin attaching message to the thread.
	p.API.LogInfo("Wrangler is attaching a message",
		"user_id", extra.UserId,
		"post_to_be_attached", postToBeAttachedID,
		"new_root_id", newRootID,
	)

	cleanPost(postToBeAttached)
	postToBeAttached.RootId = newRootID
	postToBeAttached.ParentId = newRootID

	_, appErr = p.API.CreatePost(postToBeAttached)
	if appErr != nil {
		return nil, false, errors.Wrap(appErr, "unable to create new post")
	}

	appErr = p.API.DeletePost(cleanupID)
	if appErr != nil {
		return nil, false, errors.Wrap(appErr, "unable to delete post")
	}

	p.API.LogInfo("Wrangler has attached a message",
		"user_id", extra.UserId,
		"post_to_be_attached", postToBeAttachedID,
		"new_root_id", newRootID,
	)

	msg := fmt.Sprintf("Message successfully attached to thread")

	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, msg), false, nil
}
