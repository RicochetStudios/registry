package library

import (
	"fmt"
	"testing"
)

// TestgetEnvWithDefault tests the getEnvWithDefault function.
func TestGetEnvWithDefault(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := getEnvWithDefault("", "default")
		want := "default"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got := getEnvWithDefault("value", "default")
		want := "value"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

// TestGetGameSettings tests the GetGameSettings function.
func TestGetGameSettings(t *testing.T) {
	matrix := []struct {
		name    string
		input   string
		want    []GameSetting
		wantErr error
	}{
		{"empty", "", []GameSetting{}, fmt.Errorf("game not found")},
		{"non-empty", "minecraft_java", []GameSetting{
			{Name: "VERSION", Value: "MINECRAFT_VERSION"},
			{Name: "EULA", Value: "EULA"},
			{Name: "MODLOADER", Value: "TYPE"},
			{Name: "PLAYERS", Value: "MAX_PLAYERS"},
			{Name: "NAME", Value: "MOTD"},
		}, nil},
		{"incorrect", "not_a_game", []GameSetting{}, nil},
	}

	t.Run("empty", func(t *testing.T) {
		got, err := GetGameSettings("")
		if err == nil {
			t.Errorf("got %v, want error", got)
		}
	})

	t.Run("non-empty", func(t *testing.T) {
		got, err := GetGameSettings("minecraft_java")
		if err != nil {
			t.Errorf("got error, want %v", got)
		}
	})

	t.Run("incorrect", func(t *testing.T) {
		got, err := GetGameSettings("not_a_game")
		if err == nil {
			t.Errorf("got %v, want error", got)
		}
	})
}
