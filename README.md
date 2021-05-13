# entry-task

## Table of Content
1. [Run on macOS](#run-on-macos)
    1. [Prerequisites to run on macOS](#prerequisites-to-run-on-macos)
    1. [Launch server on macOS](#launch-server-on-macos)
    1. [Test on macOS](#test-on-macos)
1. [Run on Virtual Machine using Docker](#run-on-virtual-machine-using-docker)
    1. [Prerequisites to run on Docker](#prerequisites-to-run-on-docker)
    1. [Launch server on Docker](#launch-server-on-docker)
    1. [Test on Docker](#test-on-docker)

## Run on macOS

### Prerequisites to run on macOS

`Go` and `MySQL` are installed on the computer.

### Launch server on macOS

Before running, make sure the database configuration in the call of 
`sql.Open("mysql", "{user}:{password}@{host}/{database}")` in `cmd/usermanager/main.go#setDb()` matches the settings on 
the local device. In particular, make sure the database exists.

In addition, check the paths in `cmd/paths/paths.go` under the switch case of `main` are valid. In particular, one needs
to check the `FileBasePath` directs to some valid place from the temporary file default directory.

In the `entry-task` directory, run `go run cmd/usermanager/*.go`. Do make sure the port :8080 is available. Then, open a 
browser, e.g. Chrome, and direct to `http://localhost:8080/`, one should see the following page.

![Index page](docs/screenshots/index.png)

The server can be stopped by simply press `Ctrl` + `C` in the terminal where the server is running.

### Test on macOS

#### Unit Testing

Before running, make sure the database configuration in the call of `sql.Open(...)` (same as above) in 
`test/server_helpers/util.go#SetupDb(t *testing.T)` matches the settings on the local device. 

In addition, check the paths in `cmd/paths/paths.go` under the switch case of `test` are valid. In particular, one needs
to check the `FileBasePath` directs to some valid place from the temporary file default directory.

In the `entry-task` directory, run `go test test/v*.go` for validator testing,
and run `go test test/h*.go` for handlers testing.

### Massive Testing for Performance

First, also in the `entry-task` directory, start the server by
`go run cmd/usermanager/*.go`.

Run `go test test/l*.go -parallel 100` for 1000 login requests, and 
`go test test/m*.go -parallel 100` for 1000 different request types.
One can watch the terminal where server is running to see logs.

## Run on Virtual Machine using Docker

### Prerequisites to run on Docker

`Docker` is installed on the computer.

### Launch server on Docker

#### Setup database

First, run `docker run --name=db -p 3306:3306 -e MYSQL_ROOT_HOST='%' -d mysql/mysql-server:latest --port=3306` to start 
a MySQL server on Docker. Then, run `docker exec -ti db bash` to start the bash, and run `mysql -uroot -p` in the bash 
start the MySQL command line tool. Use the password generated in logs to log in.

Then, change the password of localhost by running `ALTER USER 'root'@'localhost' IDENTIFIED BY 'password';` in MySQL. 
Also, make sure that the user, password, and host in this database matches the `sql.Open` in `cmd/usermanager/main.go`, 
for example, run `ALTER USER 'root'@'%' IDENTIFIED BY 'password';`. In addition, create the database wanted if it does 
not exist, for example, run `CREATE DATABASE IF NOT EXISTS entryTask;`

#### Launch server

After the above steps are done, we can build the image to run the web app. Before this, make sure all the settings in 
the source code match the actual setting in the Docker environment, including the two calls of `sql.Open` in 
`cmd/usermanager/main.go` and `test/server_helpers/util.go`, and the paths.

The paths on the Docker virtual machine environment for the temp file and relative path are recommended to be

```go
FileBasePath = ""
FileBaseRelativePath := "../../../../tmp"
```

Build the image called `server` by running `docker build --tag server  -f deployments/app.dockerfile . --no-cache` in 
terminal in the `entry-task` directory. Then, run the container also called `server` by 
`docker run -p 8080:8080 --name server --link db:db -d server`.

Then, open the browser,and direct to `http://localhost:8080/`, one should see the index page. Do make sure the port 
:8080 is available. Stop other servers if they are using this port.

To stop the server, one can run `docker stop server` in the terminal.

### Test on Docker

Open bash on the server by running `docker exec -ti server bash` in the terminal. Then do the same thing as on macOS to 
do the tests.
