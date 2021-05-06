FROM golang:1.16

WORKDIR /go/src/app

# COPY go.mod .
# COPY go.sum .

# RUN go mod download

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

RUN go build cmd/main.go

ENV PORT 8080

EXPOSE 8080

CMD ["./main"]


# docker build --tag first  -f deployments/Dockerfile .
# docker run -p 8080:80 --name first -d first