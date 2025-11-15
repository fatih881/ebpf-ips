package main

import (
	"log"

	"go.uber.org/zap"
)

var logger *zap.Logger

func main() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("Can't create logger: %v", err)
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			log.Fatalf("Can't stop logger: %v", err)
		}
	}(logger)
}
