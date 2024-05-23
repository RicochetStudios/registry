package client

import "context"

type Sample struct {
	// The name of the game running on the server.
	GameName string `json:"gameName"`
}

func (s *Sample) Start() error {
	return nil
}

func (s *Sample) Wait() error {
	return nil
}

func (s *Sample) Serve(context.Context) error {
	return nil
}

func (s *Sample) Stop() {}

func (s *Sample) Status() (string, error) {
	return "OK", nil
}

func (s *Sample) Healthy() (bool, error) {
	return true, nil
}

func (s *Sample) Logs() (string, error) {
	return "Logs", nil
}
