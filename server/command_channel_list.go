package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/model"
)

func (p *Plugin) runListChannelsCommand(args []string, extra *model.CommandArgs) (*model.CommandResponse, bool, error) {
	teams, appErr := p.API.GetTeamsForUser(extra.UserId)
	if appErr != nil {
		return nil, false, appErr
	}

	var msg string
	for _, team := range teams {
		msg += fmt.Sprintf("%s\n", team.Name)

		channels, appErr := p.API.GetChannelsForTeamForUser(team.Id, extra.UserId, false)
		if appErr != nil {
			return nil, false, appErr
		}

		for _, channel := range channels {
			if channel.IsGroupOrDirect() {
				continue
			}
			msg += fmt.Sprintf("%s - %s\n", channel.Id, channel.Name)
		}
		msg += "\n"
	}

	msg = codeBlock(strings.TrimRight(msg, "\n"))

	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, msg), false, nil
}
