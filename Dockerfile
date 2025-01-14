FROM golang:latest AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
 
RUN go build -o /goserver
 
EXPOSE ${SERVER_PORT}
 
CMD ["/goserver"]
