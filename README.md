<img src="https://github.com/gabrieljackson/mattermost-plugin-wrangler/blob/master/assets/profile.png?raw=true" width="75" height="75" alt="wrangler">

# Mattermost Wrangler Plugin

A [Mattermost](https://mattermost.com) plugin for managing messages.

[![CircleCI](https://circleci.com/gh/gabrieljackson/mattermost-plugin-wrangler.svg?style=shield)](https://circleci.com/gh/gabrieljackson/mattermost-plugin-wrangler)

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

## Usage

Type `/wrangler` for a list of all Wrangler commands.
