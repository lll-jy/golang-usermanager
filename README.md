# entry-task

## Run

First check `entry-task/cmd/handlers/pageHandlers.go` has the variable `var TemplateFileNameFormat = "templates/%s.html"`.

Then, in the `entry-task` directory, run `go run cmd/usermanager/*.go`.

## Test

First check `entry-task/cmd/handlers/pageHandlers.go` has the variable `var TemplateFileNameFormat = "../templates/%s.html"`.

Then, in the `entry-task` directory, run `go test test/*.go -parallel 1000`.

If logs are wanted to be shown, run with `-v` tag.
