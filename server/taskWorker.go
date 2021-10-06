//   Copyright 2021 MSolution.IO
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/awsSession"
	"github.com/trackit/trackit/config"
	"github.com/trackit/trackit/db"
	"github.com/trackit/trackit/es"
)

type (
	MessageData struct {
		TaskName   string   `json:"task_name"`
		Parameters []string `json:"parameters"`
		LogStream  string   `json:"cloudwatch_log_stream"`
	}

	LogInput struct {
		Message   string
		Timestamp int64
	}
)

const visibilityTimeoutInHours = 10
const visibilityTimeoutTaskFailedInMinutes = 20
const retryTaskOnFailure = true
const waitTimeInSeconds = 20
const esHealthcheckTimeoutInMinutes = 10
const esHealthcheckWaitTimeRetryInMinutes = 10

func taskWorker(ctx context.Context) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Running task 'worker'.", nil)
	sqsq := sqs.New(awsSession.Session)
	cwl := cloudwatchlogs.New(awsSession.Session)

	logger.Info("Retrieving SQS Queue URL.", nil)
	queueUrl, err := retrieveQueueUrl(ctx, sqsq)

	if err != nil {
		return err
	}

	logger.Info("Waiting for messages.", nil)

	var logsBuffer bytes.Buffer
	writer := io.MultiWriter(os.Stdout, &logsBuffer)
	ctx = jsonlog.ContextWithLogger(ctx, logger.WithWriter(writer))
	logger = jsonlog.LoggerFromContextOrDefault(ctx)

	timeoutStr := strconv.Itoa(esHealthcheckTimeoutInMinutes) + "m"

	for {
		message, receiptHandle, err := getNextMessage(ctx, sqsq, queueUrl)
		if err != nil {
			logsBuffer.Reset()
			continue
		}

		logger.Info("Received message, checking ES health and executing task.", map[string]interface{}{
			"message": message,
		})

		if res, err := es.Client.ClusterHealth().Timeout(timeoutStr).Do(ctx); err != nil || res.TimedOut {
			logger.Error("ES is not reachable.", map[string]interface{}{
				"timeout": timeoutStr,
				"error":   err.Error(),
			})
			logsBuffer.Reset()
			_ = changeMessageVisibility(ctx, sqsq, queueUrl, receiptHandle, 0)
			time.Sleep(time.Minute * esHealthcheckWaitTimeRetryInMinutes)
			continue
		}

		if err := db.OpenWorker(); err != nil {
			logger.Error("Database is not reachable.", map[string]interface{}{
				"error": err.Error(),
			})
			_ = changeMessageVisibility(ctx, sqsq, queueUrl, receiptHandle, 60*visibilityTimeoutTaskFailedInMinutes)
			continue
		}

		// TODO: Perhaps should avoid using built-in string type as key, so as to avoid collisions
		ctx = context.WithValue(ctx, "taskParameters", message.Parameters)

		if task, ok := tasks[message.TaskName]; ok {
			err = executeTask(ctx, task)
			if err != nil {
				logger.Error("Error while executing task. Setting visibility timeout to task failed timeout.", map[string]interface{}{
					"message":                    message,
					"taskFailedTimeoutInMinutes": 60 * visibilityTimeoutTaskFailedInMinutes,
					"error":                      err.Error(),
				})
				if retryTaskOnFailure {
					_ = changeMessageVisibility(ctx, sqsq, queueUrl, receiptHandle, 60*visibilityTimeoutTaskFailedInMinutes)
				} else {
					_ = acknowledgeMessage(ctx, sqsq, queueUrl, receiptHandle)
				}
			} else {
				logger.Info("Task done, acknowledging.", nil)
				_ = acknowledgeMessage(ctx, sqsq, queueUrl, receiptHandle)
			}
			_ = flushCloudwatchLogEvents(ctx, cwl, message, &logsBuffer, err == nil)
		} else {
			logger.Error("Unable to find requested task.", map[string]interface{}{
				"task_name": message.TaskName,
			})
			logsBuffer.Reset()
			_ = acknowledgeMessage(ctx, sqsq, queueUrl, receiptHandle)
		}
		if err := db.Close(); err != nil {
			logger.Error("Could not close connection to database.", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}
}

func executeTask(ctx context.Context, task func(context.Context) error) (err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	defer func() {
		rec := recover()
		if rec != nil {
			logger.Error("Task crashed, handling error.", map[string]interface{}{
				"error": rec,
			})
			err = errors.New("task: crashed")
		}
	}()
	return task(ctx)
}
func getNextMessage(ctx context.Context, sqsq *sqs.SQS, queueUrl *string) (MessageData, *string, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	messageData := MessageData{}
	msgResult, err := sqsq.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            queueUrl,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(60 * 60 * visibilityTimeoutInHours),
		WaitTimeSeconds:     aws.Int64(waitTimeInSeconds),
	})

	if err != nil {
		logger.Error("An error occurred while waiting for message.", map[string]interface{}{
			"queueUrl": *queueUrl,
			"error":    err.Error(),
		})
		return messageData, nil, err
	}
	if len(msgResult.Messages) == 0 {
		return messageData, nil, errors.New("queue: no message in queue.")
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

func flushCloudwatchLogEvents(ctx context.Context, cwl *cloudwatchlogs.CloudWatchLogs, message MessageData, logsBuffer *bytes.Buffer, success bool) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)

	defer func() {
		logsBuffer.Reset()
	}()

	if config.Environment != "prod" && config.Environment != "stg" {
		return nil
	}

	logGroup := config.Environment + "/task-logs/" + message.TaskName

	logStreamPrefix := message.LogStream
	if len(logStreamPrefix) == 0 {
		logStreamPrefix = "generic"
	}
	var logStreamSuffix string
	if success {
		logStreamSuffix = "succeeded"
	} else {
		logStreamSuffix = "failed"
	}
	logStream := logStreamPrefix + "/" + uuid.New().String() + "/" + logStreamSuffix

	_, err := cwl.CreateLogStream(&cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(logGroup),
		LogStreamName: aws.String(logStream),
	})
	if err != nil {
		logger.Error("Unable to create log stream. Skipping logs sending.", map[string]interface{}{
			"logGroupName":  logGroup,
			"logStreamName": logStream,
			"error":         err.Error(),
		})
		return err
	}

	logInputs := decodeLogBuffer(*logsBuffer)

	var logEvents []*cloudwatchlogs.InputLogEvent
	for _, logInput := range logInputs {
		logEvents = append(logEvents, &cloudwatchlogs.InputLogEvent{
			Message:   aws.String(logInput.Message),
			Timestamp: aws.Int64(logInput.Timestamp),
		})
	}

	_, err = cwl.PutLogEvents(&cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  aws.String(logGroup),
		LogStreamName: aws.String(logStream),
		LogEvents:     logEvents,
	})

	if err != nil {
		logger.Error("Unable to put log stream events.", map[string]interface{}{
			"logGroupName":  logGroup,
			"logStreamName": logStream,
			"error":         err.Error(),
		})
		return err
	}

	return nil
}

func decodeLogBuffer(logsBuffer bytes.Buffer) []LogInput {
	lines := strings.Split(logsBuffer.String(), "\n")

	var logInputs []LogInput

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		logInputs = append(logInputs, LogInput{
			Message:   line,
			Timestamp: getLogTimestamp(line),
		})
	}

	return logInputs
}

func getLogTimestamp(message string) int64 {
	var logJson map[string]interface{}

	err := json.Unmarshal([]byte(message), &logJson)
	if err != nil {
		return time.Now().UnixNano() / int64(time.Millisecond)
	}

	logTime, err := time.Parse(time.RFC3339, logJson["time"].(string))
	if err != nil {
		return time.Now().UnixNano() / int64(time.Millisecond)
	}

	return logTime.UnixNano() / int64(time.Millisecond)
}

func changeMessageVisibility(ctx context.Context, sqsq *sqs.SQS, queueUrl *string, receiptHandle *string, timeout int64) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	_, err := sqsq.ChangeMessageVisibility(&sqs.ChangeMessageVisibilityInput{
		ReceiptHandle:     receiptHandle,
		QueueUrl:          queueUrl,
		VisibilityTimeout: aws.Int64(timeout),
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

func acknowledgeMessage(ctx context.Context, sqsq *sqs.SQS, queueUrl *string, receiptHandle *string) error {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	_, err := sqsq.DeleteMessage(&sqs.DeleteMessageInput{
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

func retrieveQueueUrl(ctx context.Context, sqsq *sqs.SQS) (queueUrl *string, err error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	urlResult, err := sqsq.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(config.SQSQueueName),
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

func paramsFromContextOrArgs(ctx context.Context) []string {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	args, ok := ctx.Value("taskParameters").([]string)
	if !ok {
		logger.Info("Retrieving task parameters from args as context does not contain any parameters.", nil)
		args = flag.Args()
	}
	return args
}
