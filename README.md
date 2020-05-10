<img src="https://github.com/gabrieljackson/mattermost-plugin-wrangler/blob/master/assets/profile.png?raw=true" width="75" height="75" alt="wrangler">

# Mattermost Wrangler Plugin

[![CircleCI](https://circleci.com/gh/gabrieljackson/mattermost-plugin-wrangler.svg?style=shield)](https://circleci.com/gh/gabrieljackson/mattermost-plugin-wrangler) [![Go Report Card](https://goreportcard.com/badge/github.com/gabrieljackson/mattermost-plugin-wrangler)](https://goreportcard.com/report/github.com/gabrieljackson/mattermost-plugin-wrangler)

A [Mattermost](https://mattermost.com) plugin for managing messages.

## About

The wrangler plugin was created to provide advanced message management options. Currently, it supports moving messages, and any thread they may belong to, across channels and teams. The primary use-case for this is to expose a conversation to a more-appropriate channel for even greater participation.

In the future, the wrangler plugin will be developed to support additional tasks such as moving messages to new threads. Additionally, web UI integration will be supported for easier use and non-message tasks could even be implemented.

## Install

1. Go the releases page and download the latest release.
2. On your Mattermost, go to System Console -> Plugin Management and upload it.
3. Configure plugin settings as desired.
4. In order for the plugin to properly recreate messages, ensure the following system console settings are set to true:
    1. [Enable integrations to override usernames](https://docs.mattermost.com/administration/config-settings.html#enable-integrations-to-override-usernames)
    2. [Enable integrations to override profile picture icons](https://docs.mattermost.com/administration/config-settings.html#enable-integrations-to-override-profile-picture-icons)
5. Start using the plugin!

## Commands

Type `/wrangler` for a list of all Wrangler commands.

```
Wrangler Plugin - Slash Command Help

/wrangler move thread [MESSAGE_ID] [CHANNEL_ID]
  Move a given message, along with the thread it belongs to, to a given channel
    - This can be on any channel in any team that you have joined
    - Obtain the message ID by running '/wrangler list messages' or via the 'Permalink' message dropdown option (it's the last part of the URL)
    - Obtain the channel ID by running '/wrangler list channels' or via the channel 'View Info' option

/wrangler attach message [MESSAGE_ID_TO_BE_ATTACHED] [ROOT_MESSAGE_ID]
  Attach a given message to a thread in the same channel
    - Obtain the message IDs by running '/wrangler list messages' or via the 'Permalink' message dropdown option (it's the last part of the URL)

/wrangler list channels [flags]
  List the IDs of all channels you have joined
    Flags:
      --channel-filter string   A filter value that channel names must contain to be shown on the list
      --team-filter string      A filter value that team names must contain to be shown on the list

/wrangler list messages [flags]
  List the IDs of recent messages in this channel
    Flags:
      --count int         Number of messages to return. Must be between 1 and 100 (default 20)
      --trim-length int   The max character count of messages listed before they are trimmed. Must be between 10 and 500 (default 50)

/wrangler info
  Shows plugin information
```

#### /wrangler move thread

A powerful command that can "move" a message along with its parent thread to a new channel.

Note that the command works by creating new messages in the target channel, but preserves most of the original message metadata. Ordering is kept intact, but the messages contain new timestamps so that channel message history is not altered.

##### Example

A thread that was started in `channel1` is moved to `channel2`.

![channel1](https://user-images.githubusercontent.com/3694686/73672948-d1066380-467b-11ea-9772-097f9fdfcdf0.png)

The thread after being "moved" to `channel2`.

![channel2](https://user-images.githubusercontent.com/3694686/73672959-d499ea80-467b-11ea-97dc-4a2e33c8829e.png)

#### /wrangler attach message

Attaches a message that is not currently in a thread to an existing message or thread in the same channel.

This is useful for bringing normal messages about a topic into threads that they relate to.

#### /wrangler list channels

Lists channel IDs that you belong to across all teams.

#### /wrangler list channels

Lists recent message IDs from the current channel.

#### /wrangler info

Shows version and commit information for the currently-running plugin build.

## Configuration Options

The following plugin configuration is available:

 - Allowed Email Domain: an optional setting to limit plugin usage to specific users
 - Enable Wrangler Command AutoComplete: Control whether command autocomplete is enabled or not. If enabled and Allowed Email Domain is set, then some users will be able to see the Wrangler commands, but will be unable to run them.
 - Max Thread Count Move Size: an optional setting to limit the size of threads that can be moved
 - Enable Moving Threads To Different Teams: Control whether Wrangler is permitted to move message threads from one team to another or not.
 - Enable Moving Threads From Private Channels: Control whether Wrangler is permitted to move message threads from private channels or not.
 - Enable Moving Threads From Direct Message Channels: Control whether Wrangler is permitted to move message threads from direct message channels or not.
 - Enable Moving Threads From Group Message Channels: Control whether Wrangler is permitted to move message threads from group message channels or not.

## FAQ

Q: If I move a thread to a new channel, can I expect that every message will be identical to the original?
A: Currently, posts will be moved with their original message content. That said, due to technical challenges both emoji reactions and file attachments are currently unable to be properly moved to the new posts. This is something I hope to address soon.

Q: I would very much like some UI please.
A: That isn't a question, but I hear you. It's actively in the works.

Q: It would be awesome if Wrangler could do this other thing! Is that coming any time soon?
A: Please open a GitHub issue and I would be happy to see if we can implement it.
