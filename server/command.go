package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const helpText = `**Wrangler Plugin - Slash Command Help**

* |/wrangler move thread [MESSAGE_ID] [CHANNEL_ID]| - Move a given message, along with the thread it belongs to, to a given channel
  * This can be on any channel in any team that you have joined
  * Obtain the message ID by running |/wrangler list messages| or via the |Permalink| message dropdown option (it's the last part of the URL)
  * Obtain the channel ID by running |/wrangler list channels| or via the channel |View Info| option
* |/wrangler list channels| - List the IDs of all channels you have joined
* |/wrangler list messages| - List the IDs of the 20 most recent messages in this channel
* |/wrangler info| - Shows plugin information`

func getHelp() string {
	return strings.Replace(helpText, "|", "`", -1)
}

func getCommand() *model.Command {
	return &model.Command{
		Trigger:          "wrangler",
		DisplayName:      "Wrangler",
		Description:      "Manage Mattermost messages!",
		AutoComplete:     false,
		AutoCompleteDesc: "Available commands: move, list, info",
		AutoCompleteHint: "[command]",
	}
}

func getCommandResponse(responseType, text string) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: responseType,
		Text:         text,
		Username:     "wrangler",
		IconURL:      fmt.Sprintf("/plugins/%s/profile.png", manifest.ID),
	}
}

// ExecuteCommand executes a given command and returns a command response.
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	config := p.getConfiguration()

	if config.AllowedEmailDomain != "" {
		user, err := p.API.GetUser(args.UserId)
		if err != nil {
			return nil, err
		}

		if !strings.HasSuffix(user.Email, config.AllowedEmailDomain) {
			return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Permission denied. Please talk to your system administrator to get access."), nil
		}
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
