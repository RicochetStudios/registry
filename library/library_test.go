package library

import (
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
		name  string
		input string
		want GameSettings{}
		wantErr error
	}{
		{"empty", "", },
		{"non-empty", "minecraft_java", GameSettings{
			Name
		}},
		{"incorrect", "not_a_game", nil},
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
