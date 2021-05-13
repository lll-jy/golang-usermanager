FROM golang:1.16

WORKDIR /go/src/app

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

RUN go build -o usermanager cmd/usermanager/*.go

ENV PORT 8080

EXPOSE 8080

CMD ["./usermanager"]

# docker run --name=db -d mysql/mysql-server:latest

# docker run --name=db -p 3306:3306 -v mysql-volume:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=password -d mysql/mysql-server:8.0.20

# docker run --name=db -p3306:3306 -v mysql-volume:/var/lib/mysql -e MYSQL_ALLOW_EMPTY_PASSWORD=yes -d mysql/mysql-server:8.0.20
# docker exec -it db bash
# mysql -u root -p

# docker build --tag server  -f deployments/app.dockerfile .
# docker run --name server --link db:db -p 8080:8080 -d server

#### 10:44
# docker run --name=db -p 3306:3306 -v mysql-volume:/var/lib/mysql -e MYSQL_ROOT_PASSWORD="password" -d mysql/mysql-server:8.0.20
# docker exec -ti db bash

#### 11:23
# docker run --name=db -d mysql/mysql-server:latest --port=3306
# docker exec -ti db bash
# ALTER USER 'root'@'localhost' IDENTIFIED BY 'password';
# CREATE DATABASE IF NOT EXISTS entryTask;
# docker build --tag server  -f deployments/app.dockerfile . --no-cache
# docker run -p 8080:8080 --name server --link db:db -d server

### 10:41
# docker build --tag db -f deployments/db.dockerfile .
# docker run --name=db -p 3306:3306 -v mysql-volume:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=password -d db
# docker exec -ti db bash

# docker container run -it --detach --name db --env MYSQL_RANDOM_ROOT_PASSWORD=no mysql:latest

# update user set password=password("123456") where user="root";


#### 13:49
# docker run --name=db -p 3306:3306 -e MYSQL_ROOT_HOST='%' -d mysql/mysql-server:latest --port=3306
# docker exec -ti db bash
# mysql -uroot -p <GENERATED> (see log)
# ALTER USER 'root'@'localhost' IDENTIFIED BY 'password';
# ALTER USER 'root'@'%' IDENTIFIED BY 'password';
# CREATE DATABASE IF NOT EXISTS entryTask;
# docker build --tag server  -f deployments/app.dockerfile . --no-cache
# docker run -p 8080:8080 --name server --link db:db -d server