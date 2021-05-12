# entry-task

## Prerequisites

`Go` and `MySQL` are installed on the computer

## Run

In the `entry-task` directory, run `go run cmd/usermanager/*.go`.

## Test

### Unit Testing

In the `entry-task` directory, run `go test test/v*.go` for validator testing,
and run `go test test/h*.go` for handlers testing.

### Massive Testing for Performance

First, also in the `entry-task` directory, start the server by
`go run cmd/usermanager/*.go`.

Run `go test test/l*.go -parallel 100` for 1000 login requests, and 
`go test test/m*.go -parallel 100` for 1000 different request types. 

## Extension

Run using Docker. (Problem: I cannot connect to MySQL on docker.)