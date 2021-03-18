package main

import (
	"context"

	"github.com/trackit/jsonlog"
)

func taskWorker(ctx context.Context) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Running task 'worker'.", nil)

	logger.Info("Subscribing to SQS Queue.", nil)
	// TODO: Subscribe to SQS Queue

	for {


		break
	}

	return nil
}
