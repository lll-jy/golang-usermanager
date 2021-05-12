# entry-task

## Prerequisites

`Go` and `MySQL` are installed on the computer

## Run

In the `entry-task` directory, run `go run cmd/usermanager/*.go`. Then, open a browser, e.g. Chrome, and
direct to `http://localhost:8080/`, one should see the following page.
![Index page](docs/screenshots/index.png)

## Test

### Unit Testing

In the `entry-task` directory, run `go test test/v*.go` for validator testing,
and run `go test test/h*.go` for handlers testing.

### Massive Testing for Performance

First, also in the `entry-task` directory, start the server by
`go run cmd/usermanager/*.go`.

Run `go test test/l*.go -parallel 100` for 1000 login requests, and 
`go test test/m*.go -parallel 100` for 1000 different request types.
One can watch the terminal where server is running to see logs.

## Extension

Run using Docker. (Problem: I cannot connect to MySQL on docker.)