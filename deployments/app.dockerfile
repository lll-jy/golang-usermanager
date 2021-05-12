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
# docker build --tag first  -f deployments/Dockerfile .
# docker run -p 8080:8080 --name first -d first