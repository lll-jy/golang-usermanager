FROM golang:1.16

WORKDIR /go/src/app

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

RUN go build -o usermanager cmd/usermanager/*.go

ENV PORT 8080

EXPOSE 8080

CMD ["./usermanager"]

#### FINAL
### setup db
# docker run --name=db -p 3306:3306 -e MYSQL_ROOT_HOST='%' -d mysql/mysql-server:latest --port=3306
# docker exec -ti db bash
# mysql -uroot -p <GENERATED> (see log)
# ALTER USER 'root'@'localhost' IDENTIFIED BY 'password';
# ALTER USER 'root'@'%' IDENTIFIED BY 'password';
# CREATE DATABASE IF NOT EXISTS entryTask;
### launch server
# docker build --tag server  -f deployments/app.dockerfile . --no-cache
# docker run -p 8080:8080 --name server --link db:db -d server
### run tests
# docker exec -ti server bash