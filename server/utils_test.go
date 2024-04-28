package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeBotDM(t *testing.T) {
	tests := []struct {
		name     string
		base     string
		postLink string
		executor string
		expected string
	}{
		{
			name:     "all empty",
			base:     "",
			postLink: "",
			executor: "",
			expected: "",
		},
		{
			name:     "base only",
			base:     "test message",
			postLink: "",
			executor: "",
			expected: "test message",
		},
		{
			name:     "no replacements",
			base:     "test message",
			postLink: "https://domain.com/path",
			executor: "user1",
			expected: "test message",
		},
		{
			name:     "replace post link only",
			base:     "test message to {postLink}",
			postLink: "https://domain.com/path",
			executor: "user1",
			expected: "test message to https://domain.com/path",
		},
		{
			name:     "replace executor only",
			base:     "test message from {executor}",
			postLink: "https://domain.com/path",
			executor: "user1",
			expected: "test message from user1",
		},
		{
			name:     "both replaced (default)",
			base:     "@{executor} wrangled a thread you started to a new channel for you: {postLink}",
			postLink: "https://domain.com/path",
			executor: "user1",
			expected: "@user1 wrangled a thread you started to a new channel for you: https://domain.com/path",
		},
		{
			name:     "multiple replace",
			base:     "User: @{executor} @{executor} @{executor} | Link: {postLink} {postLink}",
			postLink: "https://domain.com/path",
			executor: "user1",
			expected: "User: @user1 @user1 @user1 | Link: https://domain.com/path https://domain.com/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := makeBotDM(tt.base, tt.postLink, tt.executor)
			assert.Equal(t, tt.expected, message)
		})
	}
}

func TestIsInputMessageLink(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid link",
			input:    "https://test.sampledomain.com/team_id/pl/8w89igrsffyt3ghmwsmsgyeoqe",
			expected: true,
		},
		{
			name:     "invalid link",
			input:    "https://invalid_link",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isInputMessageLink(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetMessageIDFromLink(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid link",
			input:    "https://test.sampledomain.com/team-1/pl/8w89igrsffyt3ghmwsmsgyeoqe",
			expected: "8w89igrsffyt3ghmwsmsgyeoqe",
		},
		{
			name:     "invalid link",
			input:    "https://invalid_link",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getMessageIDFromLink(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCleanInputID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid link",
			input:    "https://test.sampledomain.com/team2/pl/8w89igrsffyt3ghmwsmsgyeoqe",
			expected: "8w89igrsffyt3ghmwsmsgyeoqe",
		},
		{
			name:     "valid id",
			input:    "8w89igrsffyt3ghmwsmsgyeoqe",
			expected: "8w89igrsffyt3ghmwsmsgyeoqe",
		},
		{
			name:     "invalid link",
			input:    "https://invalid_link",
			expected: "https://invalid_link",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanInputID(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
