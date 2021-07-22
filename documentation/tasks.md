# Tasks

Tasks are a main concept of TrackIt.

They are parts of the main executable which can be executed independently. The first argument passed to the main executable is the name of the task which is going to be run. For example, `./main server` will launch the task named `server`.

Examples of tasks are:
* `server`: starts the REST API
* `ingest`: processes AWS bills from S3
* `process-account`: fetches resources status from AWS API
* `update-tags`: updates the tagging data for TagBot using the information retrieved by `process-account`

## How to run a task locally
You can use the `task.sh` script to run a task.
For example: `./task.sh process-account 1` will run the task process-account on the local environment for the AWS account with ID 1.

Tasks to run in order to get an account ready:
- `ingest {AWS ID} {BILL REPOSITORY ID}`
- `process-account {AWS ID}`
- `process-account-plugins {AWS ID}`

## How to create a new task
In order to create a task, you must add a function which takes a context. Context as parameter in the map named `task` in the file [`server/server.go`](https://github.com/trackit/trackit/blob/master/server/server.go#L60).

A task should log when it starts, ends or encounters and error. See [Logging](./logging.md).

They also report status and errors in the SQL database. See [Models](./models.md)
