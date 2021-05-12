FROM golang:1.16

WORKDIR /go/src/app

# COPY go.mod .
# COPY go.sum .

# RUN go mod download

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

RUN go build -o usermanager cmd/usermanager/*.go

ENV PORT 8080

EXPOSE 8080

CMD ["./usermanager"]

# docker run --name=db -d mysql/mysql-server:latest

# docker run --name=db -p 3306:3306 -v mysql-volume:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=root -d mysql/mysql-server:8.0.20

# docker run --name=db -p3306:3306 -v mysql-volume:/var/lib/mysql -e MYSQL_ALLOW_EMPTY_PASSWORD=yes -d mysql/mysql-server:8.0.20
# docker exec -it db bash
# mysql -u root -p

# docker build --tag server  -f deployments/app.dockerfile .
# docker run --name server --link db:db -p 8080:8080 -d server


# docker build --tag server  -f deployments/app.dockerfile .
# docker run -p 8080:8080 --name server -d server
# docker exec -it db mysql -uroot -p

# docker container run -it --detach --name db --env MYSQL_RANDOM_ROOT_PASSWORD=no mysql:latest

# update user set password=password("123456") where user="root";