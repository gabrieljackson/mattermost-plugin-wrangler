package main

import (
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCopyThreadCommand(t *testing.T) {
	team1 := &model.Team{
		Id:   model.NewId(),
		Name: "team-1",
	}
	originalChannel := &model.Channel{
		Id:     model.NewId(),
		TeamId: team1.Id,
		Name:   "original-channel",
		Type:   model.CHANNEL_OPEN,
	}
	privateChannel := &model.Channel{
		Id:     model.NewId(),
		TeamId: team1.Id,
		Name:   "private-channel",
		Type:   model.CHANNEL_PRIVATE,
	}
	directChannel := &model.Channel{
		Id:     model.NewId(),
		TeamId: team1.Id,
		Name:   "direct-channel",
		Type:   model.CHANNEL_DIRECT,
	}
	groupChannel := &model.Channel{
		Id:     model.NewId(),
		TeamId: team1.Id,
		Name:   "group-channel",
		Type:   model.CHANNEL_GROUP,
	}

	targetTeam := &model.Team{
		Id:   model.NewId(),
		Name: "target-team",
	}
	targetChannel := &model.Channel{
		Id:     model.NewId(),
		TeamId: targetTeam.Id,
		Name:   "target-channel",
	}

	reactions := []*model.Reaction{
		{
			UserId: model.NewId(),
			PostId: model.NewId(),
		},
	}

	executor := &model.User{
		Nickname: "executing user",
	}

	config := &model.Config{
		ServiceSettings: model.ServiceSettings{
			SiteURL: NewString("test.sampledomain.com"),
		},
	}

	generatedPosts := mockGeneratePostList(3, originalChannel.Id, false)

	api := &plugintest.API{}
	api.On("GetChannel", originalChannel.Id).Return(originalChannel, nil)
	api.On("GetChannel", privateChannel.Id).Return(privateChannel, nil)
	api.On("GetChannel", directChannel.Id).Return(directChannel, nil)
	api.On("GetChannel", groupChannel.Id).Return(groupChannel, nil)
	api.On("GetChannel", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(targetChannel, nil)
	api.On("GetPostThread", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(generatedPosts, nil)
	api.On("GetChannelMember", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(mockGenerateChannelMember(), nil)
	api.On("GetDirectChannel", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(directChannel, nil)
	api.On("GetTeam", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(targetTeam, nil)
	api.On("GetUser", mock.Anything).Return(executor, nil)
	api.On("CreatePost", mock.Anything, mock.Anything).Return(mockGeneratePost(), nil)
	api.On("DeletePost", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(nil)
	api.On("GetReactions", mock.AnythingOfType("string")).Return(reactions, nil)
	api.On("AddReaction", mock.Anything).Return(nil, nil)
	api.On("GetConfig", mock.Anything).Return(config)
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
		resp, isUserError, err := plugin.runCopyThreadCommand([]string{}, &model.CommandArgs{ChannelId: originalChannel.Id})
		require.NoError(t, err)
		assert.True(t, isUserError)
		assert.Contains(t, resp.Text, "Error: missing arguments")
	})

	t.Run("one arg", func(t *testing.T) {
		resp, isUserError, err := plugin.runCopyThreadCommand([]string{"id1"}, &model.CommandArgs{ChannelId: originalChannel.Id})
		require.NoError(t, err)
		assert.True(t, isUserError)
		assert.Contains(t, resp.Text, "Error: missing arguments")
	})

	t.Run("private channel", func(t *testing.T) {
		t.Run("disabled", func(t *testing.T) {
			plugin.setConfiguration(&configuration{MoveThreadFromPrivateChannelEnable: false})
			require.NoError(t, plugin.configuration.IsValid())

			resp, isUserError, err := plugin.runCopyThreadCommand([]string{"id1", "id2"}, &model.CommandArgs{ChannelId: privateChannel.Id})
			require.NoError(t, err)
			assert.False(t, isUserError)
			assert.Contains(t, resp.Text, "Wrangler is currently configured to not allow moving posts from private channels")
		})
	})

	t.Run("direct channel", func(t *testing.T) {
		t.Run("disabled", func(t *testing.T) {
			plugin.setConfiguration(&configuration{MoveThreadFromDirectMessageChannelEnable: false})
			require.NoError(t, plugin.configuration.IsValid())

			resp, isUserError, err := plugin.runCopyThreadCommand([]string{"id1", "id2"}, &model.CommandArgs{ChannelId: directChannel.Id})
			require.NoError(t, err)
			assert.False(t, isUserError)
			assert.Contains(t, resp.Text, "Wrangler is currently configured to not allow moving posts from direct message channels")
		})
	})

	t.Run("group channel", func(t *testing.T) {
		t.Run("disabled", func(t *testing.T) {
			plugin.setConfiguration(&configuration{MoveThreadFromGroupMessageChannelEnable: false})
			require.NoError(t, plugin.configuration.IsValid())

			resp, isUserError, err := plugin.runCopyThreadCommand([]string{"id1", "id2"}, &model.CommandArgs{ChannelId: groupChannel.Id})
			require.NoError(t, err)
			assert.False(t, isUserError)
			assert.Contains(t, resp.Text, "Wrangler is currently configured to not allow moving posts from group message channels")
		})
	})

	t.Run("to another team", func(t *testing.T) {
		t.Run("disabled", func(t *testing.T) {
			plugin.setConfiguration(&configuration{MoveThreadToAnotherTeamEnable: false})
			require.NoError(t, plugin.configuration.IsValid())

			resp, isUserError, err := plugin.runCopyThreadCommand([]string{"id1", "id2"}, &model.CommandArgs{ChannelId: originalChannel.Id})
			require.NoError(t, err)
			assert.False(t, isUserError)
			assert.Contains(t, resp.Text, "Wrangler is currently configured to not allow moving messages to different teams")
		})
	})

	t.Run("invalid command run location", func(t *testing.T) {
		plugin.setConfiguration(&configuration{MoveThreadToAnotherTeamEnable: true})

		t.Run("not in thread channel", func(t *testing.T) {
			resp, isUserError, err := plugin.runCopyThreadCommand([]string{"id1", "id2"}, &model.CommandArgs{ChannelId: model.NewId()})
			require.NoError(t, err)
			assert.True(t, isUserError)
			assert.Contains(t, resp.Text, "Error: this command must be run from the channel containing the post")
		})

		postSlice := generatedPosts.ToSlice()
		rootPostID := postSlice[len(postSlice)-1].Id

		t.Run("in thread being copied", func(t *testing.T) {
			t.Run("parentId matches", func(t *testing.T) {
				resp, isUserError, err := plugin.runCopyThreadCommand([]string{"id1", "id2"}, &model.CommandArgs{ChannelId: originalChannel.Id, ParentId: rootPostID})
				require.NoError(t, err)
				assert.True(t, isUserError)
				assert.Contains(t, resp.Text, "Error: this command cannot be run from inside the thread; please run directly in the channel containing the thread")
			})

			t.Run("rootId matches", func(t *testing.T) {
				resp, isUserError, err := plugin.runCopyThreadCommand([]string{"id1", "id2"}, &model.CommandArgs{ChannelId: originalChannel.Id, RootId: rootPostID})
				require.NoError(t, err)
				assert.True(t, isUserError)
				assert.Contains(t, resp.Text, "Error: this command cannot be run from inside the thread; please run directly in the channel containing the thread")
			})
		})
	})

	t.Run("copy thread successfully", func(t *testing.T) {
		require.NoError(t, plugin.configuration.IsValid())

		resp, isUserError, err := plugin.runCopyThreadCommand([]string{"id1", "id2"}, &model.CommandArgs{ChannelId: originalChannel.Id})
		require.NoError(t, err)
		assert.False(t, isUserError)
		assert.Contains(t, resp.Text, "Thread copy complete")
	})

	t.Run("thread is above configuration move-maximum", func(t *testing.T) {
		plugin.setConfiguration(&configuration{MoveThreadMaxCount: "1"})
		require.NoError(t, plugin.configuration.IsValid())
		resp, isUserError, err := plugin.runCopyThreadCommand([]string{"id1", "id2"}, &model.CommandArgs{ChannelId: model.NewId()})
		require.NoError(t, err)
		assert.True(t, isUserError)
		assert.Contains(t, resp.Text, "Error: the thread is 3 posts long, but this command is configured to only move threads of up to 1 posts")
	})
}
