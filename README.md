# entry-task

## Prerequisites

`Go` and `MySQL` are installed on the computer

## Run

In the `entry-task` directory, run `go run cmd/usermanager/*.go`.

## Test

In the `entry-task` directory, run `go test test/*.go -parallel 1000`. If different test want to be run independently, use the following commands.

### Validator Unit Testing

Run `go test test/validation_test.go`.

### Handlers Unit Testing

Run `go test test/server*.go`.

### Massive Testing for Performance

Run `go test`. 

## Extension

Run using Docker. (Problem: I cannot connect to MySQL on docker.)