# entry-task

## Prerequisites

`Go` and `MySQL` are installed on the computer

## Run

First check `entry-task/cmd/handlers/pageHandlers.go` has the variable `var TemplateFileNameFormat = "templates/%s.html"`.

Then, in the `entry-task` directory, run `go run cmd/usermanager/*.go`.

## Test

First check `entry-task/cmd/handlers/pageHandlers.go` has the variable `var TemplateFileNameFormat = "../templates/%s.html"`. And check `entry-task/cmd/paths/paths.go`
has the variable `var FileBasePath = "../../../../Desktop/EntryTask/entry-task/test/data/upload"`

Then, in the `entry-task` directory, run `go test test/*.go -parallel 1000`.

If logs are wanted to be shown, run with `-v` tag.

## Extension

Run using Docker. (Problem: I cannot connect to MySQL on docker.)