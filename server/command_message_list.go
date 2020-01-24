package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/model"
)

func (p *Plugin) runListMessagesCommand(args []string, extra *model.CommandArgs) (*model.CommandResponse, bool, error) {
	var msg string

	channelPosts, appErr := p.API.GetPostsForChannel(extra.ChannelId, 0, 20)
	if appErr != nil {
		return nil, false, appErr
	}

	msg += "The last 20 messages in this channel:\n"
	for _, post := range channelPosts.ToSlice() {
		if post.IsSystemMessage() {
			msg += "[     system message     ] - <skipped>\n"
		} else {
			msg += fmt.Sprintf("%s - %s\n", post.Id, trimMessage(post.Message))
		}
	}

	msg = codeBlock(strings.TrimRight(msg, "\n"))

	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, msg), false, nil
}

func trimMessage(message string) string {
	if len(message) <= 50 {
		return message
	}

	return fmt.Sprintf("%s...", message[:50])
}
