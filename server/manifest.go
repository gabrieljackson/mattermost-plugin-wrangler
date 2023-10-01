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
  "version": "0.7.0",
  "min_server_version": "5.12.0",
  "server": {
    "executables": {
      "darwin-amd64": "server/dist/plugin-darwin-amd64",
      "darwin-arm64": "server/dist/plugin-darwin-arm64",
      "freebsd-amd64": "server/dist/plugin-freebsd-amd64",
      "linux-amd64": "server/dist/plugin-linux-amd64",
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
        "key": "PermittedWranglerUsers",
        "display_name": "Permitted Wrangler Users",
        "type": "dropdown",
        "help_text": "Choose who is allowed to use the Wrangler plugin. (Other permissions below still apply)",
        "placeholder": "",
        "default": "system-admins",
        "options": [
          {
            "display_name": "System administrators only",
            "value": "system-admins"
          },
          {
            "display_name": "System administrators and users from the 'Allowed Email Domain' list",
            "value": "system-admins-and-email-domain"
          },
          {
            "display_name": "All users",
            "value": "all-users"
          }
        ]
      },
      {
        "key": "AllowedEmailDomain",
        "display_name": "Allowed Email Domain",
        "type": "text",
        "help_text": "(Optional) When set, users must have an email ending in this domain to use Wrangler. Multiple domains can be specified by separating them with commas.",
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
      },
      {
        "key": "ThreadAttachMessage",
        "display_name": "Info-Message: Attached a Message",
        "type": "text",
        "help_text": "The message being sent to the user after attaching his message to a thread. Allowed variables: {executor}, {postLink}",
        "placeholder": "",
        "default": "@{executor} wrangled one of your messages into a thread for you: {postLink}"
      },
      {
        "key": "MoveThreadMessage",
        "display_name": "Info-Message: Moved a Thread",
        "type": "text",
        "help_text": "The message being sent to the user after moving a thread. Allowed variables: {executor}, {postLink}",
        "placeholder": "",
        "default": "@{executor} wrangled a thread you started to a new channel for you: {postLink}"
      },
      {
        "key": "CopyThreadMessage",
        "display_name": "Info-Message: Copied a Thread",
        "type": "text",
        "help_text": "The message being sent to the user after copying a message. Allowed variables: {executor}, {postLink}",
        "placeholder": "",
        "default": "@{executor} wrangled a thread you started to a new channel for you: {postLink}"
      }
    ]
  }
}
`

func init() {
	manifest = model.ManifestFromJson(strings.NewReader(manifestStr))
}
