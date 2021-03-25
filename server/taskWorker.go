package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"
	"github.com/trackit/trackit/config"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/awsSession"
)

type (
	MessageData struct {
		TaskName	string			`json:"task_name"`
		ParamsData	[]string		`json:"params_data"`
	}
)

const Visibility = -1

func taskWorker(ctx context.Context) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	logger.Info("Running task 'worker'.", nil)
	svc := sqs.New(awsSession.Session)

	logger.Info("Retrieving SQS Queue URL.", nil)
	queueUrl, err := retrieveQueueUrl(ctx, svc)

	if err != nil {
		return err
	}

	for {
		message, receiptHandle, err := getNextMessage(ctx, svc, queueUrl)

		if err != nil {
			continue
		}

		logger.Info("Received message, executing task.", map[string]interface{}{
			"message": message,
		})

		ctx := context.WithValue(ctx, "paramsData", message.ParamsData)

		if task, ok := tasks[config.Task]; ok {
			err = task(ctx)

			if err != nil {
				logger.Error("Error while executing task. Resetting visibility timeout to 0.", map[string]interface{}{
					"message": message,
					"error":   err.Error(),
				})
				_ = resetMessageTimeout(ctx, svc, queueUrl, receiptHandle)
				continue
			}

			logger.Info("Task done, acknowledging...", nil)

			_ = acknowledgeMessage(ctx, svc, queueUrl, receiptHandle)
		}
	}
}

func getNextMessage(ctx context.Context, svc *sqs.SQS, queueUrl *string) (MessageData, *string, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	messageData := MessageData{}

	logger.Info("Waiting for messages.", nil)

	msgResult, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl: queueUrl,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout: aws.Int64(60 * 15),
	})

	if err != nil {
		logger.Error("Unable to read message.", map[string]interface{}{
			"queueUrl": *queueUrl,
			"error":    err.Error(),
		})
		return messageData, nil, err
	}

	if len(msgResult.Messages) == 0 {
		logger.Info("No message in queue.", nil)
		return messageData, nil, errors.New("no message in queue.")
	}

	err = json.Unmarshal([]byte(*msgResult.Messages[0].Body), &messageData)

	if err != nil {
		logger.Error("Unable to decode message.", map[string]interface{}{
			"messageBody": *msgResult.Messages[0].Body,
			"error":       err.Error(),
		})
		return messageData, nil, err
	}

	return messageData, msgResult.Messages[0].ReceiptHandle, nil
}

func resetMessageTimeout(ctx context.Context, svc *sqs.SQS, queueUrl *string, receiptHandle *string) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	_, err := svc.ChangeMessageVisibility(&sqs.ChangeMessageVisibilityInput{
		ReceiptHandle:     receiptHandle,
		QueueUrl:          queueUrl,
		VisibilityTimeout: aws.Int64(0),
	})

	if err != nil {
		logger.Error("Unable to reset timeout.", map[string]interface{}{
			"messageHandle": *receiptHandle,
			"error":         err.Error(),
		})
		return err
	}

	return nil
}

func acknowledgeMessage(ctx context.Context, svc *sqs.SQS, queueUrl *string, receiptHandle *string) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	_, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      queueUrl,
		ReceiptHandle: receiptHandle,
	})

	if err != nil {
		logger.Error("Unable to acknowledge message.", map[string]interface{}{
			"messageHandle": *receiptHandle,
			"error":         err.Error(),
		})
		return err
	}

	return nil
}

func retrieveQueueUrl(ctx context.Context, svc *sqs.SQS) (queueUrl *string, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	urlResult, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &config.SQSQueueName,
	})

	if err != nil {
		logger.Error("Unable to get queue URL from name.", map[string]interface{}{
			"queueName": config.SQSQueueName,
			"error":     err.Error(),
		})
		return nil, err
	}

	return urlResult.QueueUrl, nil
}

func paramsFromContextOrArgs(ctx context.Context) ([]string) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	args, ok := ctx.Value("paramsData").([]string)

	if !ok {
		logger.Info("Retrieving task parameters from args as context does not contain any parameters.", nil)
		args = flag.Args()
	}

	return args
}
