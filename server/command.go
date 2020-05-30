package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const helpText = `Wrangler Plugin - Slash Command Help

%s

%s

/wrangler attach message [MESSAGE_ID_TO_BE_ATTACHED] [ROOT_MESSAGE_ID]
  Attach a given message to a thread in the same channel
    - Obtain the message IDs by running '/wrangler list messages' or via the 'Permalink' message dropdown option (it's the last part of the URL)

/wrangler list channels [flags]
  List the IDs of all channels you have joined
	Flags:
%s
/wrangler list messages [flags]
  List the IDs of recent messages in this channel
    Flags:
%s
/wrangler info
  Shows plugin information`

func getHelp() string {
	return codeBlock(fmt.Sprintf(
		helpText,
		moveThreadUsage,
		copyThreadUsage,
		getListChannelsFlagSet().FlagUsages(),
		getListMessagesFlagSet().FlagUsages(),
	))
}

func getCommand(autocomplete bool) *model.Command {
	return &model.Command{
		Trigger:          "wrangler",
		DisplayName:      "Wrangler",
		Description:      "Manage Mattermost messages!",
		AutoComplete:     autocomplete,
		AutoCompleteDesc: "Available commands: move thread, attach message, list messages, list channels, info",
		AutoCompleteHint: "[command]",
	}
}

func getCommandResponse(responseType, text string) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: responseType,
		Text:         text,
		Username:     "wrangler",
		IconURL:      fmt.Sprintf("/plugins/%s/profile.png", manifest.Id),
	}
}

// ExecuteCommand executes a given command and returns a command response.
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	if !p.authorizedPluginUser(args.UserId) {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Permission denied. Please talk to your system administrator to get access."), nil
	}

	stringArgs := strings.Split(args.Command, " ")

	if len(stringArgs) < 2 {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, getHelp()), nil
	}

	command := stringArgs[1]

	var handler func([]string, *model.CommandArgs) (*model.CommandResponse, bool, error)

	switch command {
	case "move":
		if len(stringArgs) < 3 {
			break
		}

		switch stringArgs[2] {
		case "thread":
			handler = p.runMoveThreadCommand
			stringArgs = stringArgs[3:]
		}
	case "copy":
		if len(stringArgs) < 3 {
			break
		}

		switch stringArgs[2] {
		case "thread":
			handler = p.runCopyThreadCommand
			stringArgs = stringArgs[3:]
		}
	case "attach":
		if len(stringArgs) < 3 {
			break
		}

		switch stringArgs[2] {
		case "message":
			handler = p.runAttachMessageCommand
			stringArgs = stringArgs[3:]
		}
	case "list":
		if len(stringArgs) < 3 {
			break
		}

		switch stringArgs[2] {
		case "channels":
			handler = p.runListChannelsCommand
			stringArgs = stringArgs[3:]
		case "messages":
			handler = p.runListMessagesCommand
			stringArgs = stringArgs[3:]
		}
	case "info":
		handler = p.runInfoCommand
		stringArgs = stringArgs[2:]
	}

	if handler == nil {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, getHelp()), nil
	}

	resp, userError, err := handler(stringArgs, args)

	if err != nil {
		p.API.LogError(err.Error())
		if userError {
			return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, fmt.Sprintf("__Error: %s__\n\nRun `/wrangler help` for usage instructions.", err.Error())), nil
		}

		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "An unknown error occurred. Please talk to your administrator for help."), nil
	}

	return resp, nil
}

func (p *Plugin) runInfoCommand(args []string, extra *model.CommandArgs) (*model.CommandResponse, bool, error) {
	resp := fmt.Sprintf("Wrangler plugin version: %s, "+
		"[%s](https://github.com/gabrieljackson/mattermost-plugin-wrangler/commit/%s), built %s\n\n",
		manifest.Version, BuildHashShort, BuildHash, BuildDate)

	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, resp), false, nil
}

func (p *Plugin) authorizedPluginUser(userID string) bool {
	config := p.getConfiguration()

	if len(config.AllowedEmailDomain) != 0 {
		user, err := p.API.GetUser(userID)
		if err != nil {
			return false
		}

		if !strings.HasSuffix(user.Email, config.AllowedEmailDomain) {
			return false
		}
	}

	return true
}
