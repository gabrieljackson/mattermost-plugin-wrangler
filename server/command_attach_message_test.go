package main

import (
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAttachMessageCommand(t *testing.T) {
	team1 := &model.Team{
		Id:   model.NewId(),
		Name: "team-1",
	}
	channel1 := &model.Channel{
		Id:     model.NewId(),
		TeamId: team1.Id,
		Name:   "channel1",
		Type:   model.CHANNEL_OPEN,
	}
	postToBeAttached := &model.Post{
		Id:        model.NewId(),
		ChannelId: channel1.Id,
	}
	postToAttachTo := &model.Post{
		Id:        model.NewId(),
		ChannelId: channel1.Id,
	}
	rootID := model.NewId()
	postInThreadAlready := &model.Post{
		Id:        model.NewId(),
		ChannelId: channel1.Id,
		RootId:    rootID,
		ParentId:  rootID,
	}

	postInAnotherChannel := &model.Post{
		Id:        model.NewId(),
		ChannelId: model.NewId(),
	}

	api := &plugintest.API{}
	api.On("GetPost", postToBeAttached.Id).Return(postToBeAttached, nil)
	api.On("GetPost", postToAttachTo.Id).Return(postToAttachTo, nil)
	api.On("GetPost", postInThreadAlready.Id).Return(postInThreadAlready, nil)
	api.On("GetPost", postInAnotherChannel.Id).Return(postInAnotherChannel, nil)
	api.On("GetPost", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(nil, model.NewAppError("where", model.NewId(), nil, "not found", 0))
	api.On("CreatePost", mock.Anything, mock.Anything).Return(mockGeneratePost(), nil)
	api.On("DeletePost", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(nil)
	api.On("LogInfo",
		mock.AnythingOfTypeArgument("string"),
		mock.AnythingOfTypeArgument("string"),
		mock.AnythingOfTypeArgument("string"),
		mock.AnythingOfTypeArgument("string"),
		mock.AnythingOfTypeArgument("string"),
		mock.AnythingOfTypeArgument("string"),
		mock.AnythingOfTypeArgument("string"),
	).Return(nil)

	var plugin Plugin
	plugin.SetAPI(api)

	t.Run("no args", func(t *testing.T) {
		resp, isUserError, err := plugin.runAttachMessageCommand([]string{}, &model.CommandArgs{ChannelId: channel1.Id})
		require.NoError(t, err)
		assert.True(t, isUserError)
		assert.Contains(t, resp.Text, "Error: missing arguments")
	})

	t.Run("one arg", func(t *testing.T) {
		resp, isUserError, err := plugin.runAttachMessageCommand([]string{"id1"}, &model.CommandArgs{ChannelId: channel1.Id})
		require.NoError(t, err)
		assert.True(t, isUserError)
		assert.Contains(t, resp.Text, "Error: missing arguments")
	})

	t.Run("post to be attached invalid", func(t *testing.T) {
		resp, isUserError, err := plugin.runAttachMessageCommand([]string{model.NewId(), postToAttachTo.Id}, &model.CommandArgs{ChannelId: model.NewId()})
		require.NoError(t, err)
		assert.True(t, isUserError)
		assert.Contains(t, resp.Text, "Error: unable to get message with ID")
	})

	t.Run("post to be attached to invalid", func(t *testing.T) {
		resp, isUserError, err := plugin.runAttachMessageCommand([]string{postToBeAttached.Id, model.NewId()}, &model.CommandArgs{ChannelId: model.NewId()})
		require.NoError(t, err)
		assert.True(t, isUserError)
		assert.Contains(t, resp.Text, "Error: unable to get message with ID")
	})

	t.Run("not in channel with messages", func(t *testing.T) {
		resp, isUserError, err := plugin.runAttachMessageCommand([]string{postToBeAttached.Id, postToAttachTo.Id}, &model.CommandArgs{ChannelId: model.NewId()})
		require.NoError(t, err)
		assert.True(t, isUserError)
		assert.Contains(t, resp.Text, "Error: the attach command must be run from the channel containing the messages")
	})

	t.Run("attach to message in another channel", func(t *testing.T) {
		resp, isUserError, err := plugin.runAttachMessageCommand([]string{postToBeAttached.Id, postInAnotherChannel.Id}, &model.CommandArgs{ChannelId: channel1.Id})
		require.NoError(t, err)
		assert.True(t, isUserError)
		assert.Contains(t, resp.Text, "Error: unable to attach message to a thread in another channel")
	})

	t.Run("attach message already in another thread", func(t *testing.T) {
		resp, isUserError, err := plugin.runAttachMessageCommand([]string{postInThreadAlready.Id, postToAttachTo.Id}, &model.CommandArgs{ChannelId: channel1.Id})
		require.NoError(t, err)
		assert.True(t, isUserError)
		assert.Contains(t, resp.Text, "Error: the message to be attached is already part of a thread")
	})

	t.Run("attach message successfully", func(t *testing.T) {
		plugin.setConfiguration(&configuration{MoveThreadToAnotherTeamEnable: true})
		require.NoError(t, plugin.configuration.IsValid())

		resp, isUserError, err := plugin.runAttachMessageCommand([]string{postToBeAttached.Id, postToAttachTo.Id}, &model.CommandArgs{ChannelId: channel1.Id})
		require.NoError(t, err)
		assert.False(t, isUserError)
		assert.Contains(t, resp.Text, "Message successfully attached to thread")
	})
}
