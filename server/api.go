package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mattermost/mattermost-server/plugin"
)

// ServeHTTP handles HTTP requests to the plugin.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	config := p.getConfiguration()

	if err := config.IsValid(); err != nil {
		http.Error(w, "This plugin is not configured.", http.StatusNotImplemented)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch path := r.URL.Path; path {
	case "/profile.png":
		p.handleProfileImage(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (p *Plugin) handleProfileImage(w http.ResponseWriter, r *http.Request) {
	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.API.LogError("Unable to get bundle path, err=" + err.Error())
		return
	}

	img, err := os.Open(filepath.Join(bundlePath, "assets", "profile.png"))
	if err != nil {
		http.NotFound(w, r)
		p.API.LogError("Unable to read profile image, err=" + err.Error())
		return
	}
	defer img.Close()

	w.Header().Set("Content-Type", "image/png")
	io.Copy(w, img)
}

// TeamAndChannelsResponse is an API response for a list of teams and channels.
type TeamAndChannelsResponse struct {
	Response map[string]Team `json:"response"`
}

// Team is a list of channels in a team.
type Team struct {
	Name     string    `json:"name"`
	Channels []Channel `json:"channels"`
}

// Channel is a list of channels in a team.
type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// TODO: this is a stub for future webapp support
func (p *Plugin) handleTeamsAndChannelsRequest(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	teams, appErr := p.API.GetTeamsForUser(userID)
	if appErr != nil {
		p.API.LogError("Unable to get teams for user err=" + appErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	teamMap := make(map[string]Team)
	for _, team := range teams {
		channels, appErr := p.API.GetChannelsForTeamForUser(team.Id, userID, false)
		if appErr != nil {
			p.API.LogError("Unable to get channels for user err=" + appErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var channelResponse []Channel
		for _, channel := range channels {
			if channel.IsGroupOrDirect() {
				continue
			}
			channelResponse = append(
				channelResponse,
				Channel{ID: channel.Id, Name: channel.Name},
			)
		}

		teamMap[team.Id] = Team{Name: team.Name, Channels: channelResponse}
	}

	response := TeamAndChannelsResponse{teamMap}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		p.API.LogError("Unable marhsal channel list to json err=" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	p.API.LogError(string(responseJSON))

	w.Write(responseJSON)
}
