FROM golang:1.23-alpine
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download && go mod verify

COPY main.go ./

RUN go build

EXPOSE 8080

ENTRYPOINT ./cup-stacker