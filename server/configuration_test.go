package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigurationIsValid(t *testing.T) {
	baseConfiguration := configuration{
		AllowedEmailDomain: "mattermost.com",
		MoveThreadMaxCount: "10",
	}

	t.Run("valid", func(t *testing.T) {
		require.NoError(t, baseConfiguration.IsValid())
	})

	t.Run("MaxThreadCountMoveSize", func(t *testing.T) {
		config := baseConfiguration

		t.Run("invalid integer", func(t *testing.T) {
			config.MoveThreadMaxCount = "twenty"
			require.Error(t, config.IsValid())
		})

		t.Run("negative integer", func(t *testing.T) {
			config.MoveThreadMaxCount = "-10"
			err := config.IsValid()
			if err == nil {
				t.Log("WTF")
			}
			t.Log(config.MaxThreadCountMoveSizeInt())
			require.Error(t, config.IsValid())
		})

		t.Run("unset value", func(t *testing.T) {
			config.MoveThreadMaxCount = ""
			require.NoError(t, config.IsValid())
		})
	})
}
