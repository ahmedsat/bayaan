package main

import "github.com/ahmedsat/bayaan"

func main() {

	logger := bayaan.NewLogger(bayaan.WithLevel(bayaan.LoggerLevelInfo))
	defer logger.Close()

	bayaan.WithFields(bayaan.Fields{"fake": "fake"})(logger)

	logger.Info("Hello, world!", bayaan.Fields{"key": "value"})
	logger.Warn("Hello, world!", bayaan.Fields{"key": "value"})
	logger.Error("Hello, world!", bayaan.Fields{"key": "value"})
}
