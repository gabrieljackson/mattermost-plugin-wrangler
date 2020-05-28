// This file is automatically generated. Do not modify it manually.

package main

import (
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

var manifest *model.Manifest

const manifestStr = `
{
  "id": "com.mattermost.wrangler",
  "name": "Wrangler",
  "description": "Manage messages across teams and channels",
  "version": "0.4.1",
  "min_server_version": "5.12.0",
  "server": {
    "executables": {
      "linux-amd64": "server/dist/plugin-linux-amd64",
      "darwin-amd64": "server/dist/plugin-darwin-amd64",
      "windows-amd64": "server/dist/plugin-windows-amd64.exe"
    },
    "executable": ""
  },
  "webapp": {
    "bundle_path": "webapp/dist/main.js"
  },
  "settings_schema": {
    "header": "",
    "footer": "",
    "settings": [
      {
        "key": "AllowedEmailDomain",
        "display_name": "Allowed Email Domain",
        "type": "text",
        "help_text": "(Optional) When set, users must have an email ending in this domain to use the wrangler slash command.",
        "placeholder": "",
        "default": null
      },
      {
        "key": "EnableWebUI",
        "display_name": "Enable Wrangler webapp functionality [BETA]",
        "type": "bool",
        "help_text": "Enable the work-in-progress Wrangler webapp functionality.",
        "placeholder": "",
        "default": false
      },
      {
        "key": "CommandAutoCompleteEnable",
        "display_name": "Enable Wrangler Command AutoComplete",
        "type": "bool",
        "help_text": "Control whether command autocomplete is enabled or not. If enabled and Allowed Email Domain is set, then some users will be able to see the Wrangler commands, but will be unable to run them.",
        "placeholder": "",
        "default": false
      },
      {
        "key": "MoveThreadMaxCount",
        "display_name": "Max Thread Count Move Size",
        "type": "text",
        "help_text": "The maximum number of messages in a thread that the plugin is allowed to move. Leave empty for unlimited messages.",
        "placeholder": "",
        "default": null
      },
      {
        "key": "MoveThreadToAnotherTeamEnable",
        "display_name": "Enable Moving Threads To Different Teams",
        "type": "bool",
        "help_text": "Control whether Wrangler is permitted to move message threads from one team to another or not.",
        "placeholder": "",
        "default": false
      },
      {
        "key": "MoveThreadFromPrivateChannelEnable",
        "display_name": "Enable Moving Threads From Private Channels",
        "type": "bool",
        "help_text": "Control whether Wrangler is permitted to move message threads from private channels or not.",
        "placeholder": "",
        "default": false
      },
      {
        "key": "MoveThreadFromDirectMessageChannelEnable",
        "display_name": "Enable Moving Threads From Direct Message Channels",
        "type": "bool",
        "help_text": "Control whether Wrangler is permitted to move message threads from direct message channels or not.",
        "placeholder": "",
        "default": false
      },
      {
        "key": "MoveThreadFromGroupMessageChannelEnable",
        "display_name": "Enable Moving Threads From Group Message Channels",
        "type": "bool",
        "help_text": "Control whether Wrangler is permitted to move message threads from group message channels or not.",
        "placeholder": "",
        "default": false
      }
    ]
  }
}
`

func init() {
	manifest = model.ManifestFromJson(strings.NewReader(manifestStr))
}
