package main

import (
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMessagelListCommand(t *testing.T) {
	testChannel := model.Channel{
		Id:   model.NewId(),
		Name: "test-channel",
	}

	testPostList := mockGeneratePostList(3, testChannel.Id, false)

	api := &plugintest.API{}
	api.On("GetPostsForChannel", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(testPostList, nil)

	var plugin Plugin
	plugin.SetAPI(api)

	t.Run("list messages successfully", func(t *testing.T) {
		resp, isUserError, err := plugin.runListMessagesCommand([]string{}, &model.CommandArgs{ChannelId: testChannel.Id})
		require.NoError(t, err)
		assert.False(t, isUserError)
		for _, post := range testPostList.ToSlice() {
			assert.Contains(t, resp.Text, post.Id)
			assert.Contains(t, resp.Text, post.Message)
		}
	})

	t.Run("specify valid count", func(t *testing.T) {
		resp, isUserError, err := plugin.runListMessagesCommand([]string{"--count=50"}, &model.CommandArgs{ChannelId: testChannel.Id})
		require.NoError(t, err)
		assert.False(t, isUserError)
		assert.Contains(t, resp.Text, "The last 50 messages in this channel")
		for _, post := range testPostList.ToSlice() {
			assert.Contains(t, resp.Text, post.Id)
			assert.Contains(t, resp.Text, post.Message)
			assert.Contains(t, resp.Text, post.Message)
		}
	})

	t.Run("specify count that is too low", func(t *testing.T) {
		_, isUserError, err := plugin.runListMessagesCommand([]string{"--count=-1"}, &model.CommandArgs{ChannelId: testChannel.Id})
		require.Error(t, err)
		assert.True(t, isUserError)
		assert.Contains(t, err.Error(), "count (-1) must be between 1 and 100")
	})

	t.Run("specify count that is too high", func(t *testing.T) {
		_, isUserError, err := plugin.runListMessagesCommand([]string{"--count=120"}, &model.CommandArgs{ChannelId: testChannel.Id})
		require.Error(t, err)
		assert.True(t, isUserError)
		assert.Contains(t, err.Error(), "count (120) must be between 1 and 100")
	})

	t.Run("list messages successfully with system", func(t *testing.T) {
		testPostList := mockGeneratePostList(3, testChannel.Id, true)

		api := &plugintest.API{}
		api.On("GetPostsForChannel", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(testPostList, nil)

		var plugin Plugin
		plugin.SetAPI(api)

		resp, isUserError, err := plugin.runListMessagesCommand([]string{}, &model.CommandArgs{ChannelId: testChannel.Id})
		require.NoError(t, err)
		assert.False(t, isUserError)
		assert.Contains(t, resp.Text, "[     system message     ] - <skipped>")
	})
}

func TestTrimMessage(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    string
	}{
		{
			name:    "short message",
			message: "short",
			want:    "short",
		},
		{
			name:    "max (50) characters",
			message: "12345678901234567890123456789012345678901234567890",
			want:    "12345678901234567890123456789012345678901234567890",
		},
		{
			name:    "max+1 (51) characters",
			message: "123456789012345678901234567890123456789012345678901",
			want:    "12345678901234567890123456789012345678901234567890...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, trimMessage(tt.message))
		})
	}
}
