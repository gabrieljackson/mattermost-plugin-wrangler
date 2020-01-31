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

	testPostList := mockGeneratePostList(3, testChannel.Id)

	api := &plugintest.API{}
	api.On("GetPostsForChannel", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(testPostList, nil)

	var plugin Plugin
	plugin.SetAPI(api)

	t.Run("list channels successfully", func(t *testing.T) {
		resp, isUserError, err := plugin.runListMessagesCommand([]string{}, &model.CommandArgs{ChannelId: testChannel.Id})
		require.NoError(t, err)
		assert.False(t, isUserError)
		for _, post := range testPostList.ToSlice() {
			assert.Contains(t, resp.Text, post.Id)
			assert.Contains(t, resp.Text, post.Message)
		}
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
