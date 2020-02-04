package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/spf13/pflag"
)

const (
	flagCount            = "count"
	minListMessagesCount = 1
	maxListMessagesCount = 100

	flagTrimLength            = "trim-length"
	minListMessagesTrimLength = 10
	maxListMessagesTrimLength = 500
)

type listOptions struct {
	count      int
	trimLength int
}

func getListMessagesFlagSet() *pflag.FlagSet {
	listMessagesFlagSet := pflag.NewFlagSet("list messages", pflag.ContinueOnError)
	listMessagesFlagSet.Int(flagCount, 20, fmt.Sprintf("Number of messages to return. Must be between %d and %d", minListMessagesCount, maxListMessagesCount))
	listMessagesFlagSet.Int(flagTrimLength, 50, fmt.Sprintf("The max character count of messages listed before they are trimmed. Must be between %d and %d", minListMessagesTrimLength, maxListMessagesTrimLength))

	return listMessagesFlagSet
}

func parseListMessagesArgs(args []string) (listOptions, error) {
	var options listOptions

	listMessagesFlagSet := getListMessagesFlagSet()
	err := listMessagesFlagSet.Parse(args)
	if err != nil {
		return options, err
	}

	options.count, err = listMessagesFlagSet.GetInt(flagCount)
	if err != nil {
		return options, err
	}
	if options.count < minListMessagesCount || options.count > maxListMessagesCount {
		return options, fmt.Errorf("%s (%d) must be between %d and %d", flagCount, options.count, minListMessagesCount, maxListMessagesCount)
	}

	options.trimLength, err = listMessagesFlagSet.GetInt(flagTrimLength)
	if err != nil {
		return options, err
	}
	if options.trimLength < minListMessagesTrimLength || options.trimLength > maxListMessagesTrimLength {
		return options, fmt.Errorf("%s (%d) must be between %d and %d", flagTrimLength, options.trimLength, minListMessagesTrimLength, maxListMessagesTrimLength)
	}

	return options, err
}

func (p *Plugin) runListMessagesCommand(args []string, extra *model.CommandArgs) (*model.CommandResponse, bool, error) {
	options, err := parseListMessagesArgs(args)
	if err != nil {
		return nil, true, err
	}

	channelPosts, appErr := p.API.GetPostsForChannel(extra.ChannelId, 0, options.count)
	if appErr != nil {
		return nil, false, appErr
	}

	msg := fmt.Sprintf("The last %d messages in this channel:\n", options.count)
	for _, post := range channelPosts.ToSlice() {
		if post.IsSystemMessage() {
			msg += "[     system message     ] - <skipped>\n"
		} else {
			msg += fmt.Sprintf("%s - %s\n", post.Id, trimMessage(post.Message, options.trimLength))
		}
	}

	msg = codeBlock(strings.TrimRight(msg, "\n"))

	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, msg), false, nil
}

func trimMessage(message string, trimLength int) string {
	if len(message) <= trimLength {
		return message
	}

	return fmt.Sprintf("%s...", message[:trimLength])
}
