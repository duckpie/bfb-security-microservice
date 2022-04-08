FROM golang:1.18-alpine

WORKDIR /app

COPY . /app/

RUN go mod download
RUN go build -o app main.go

ENTRYPOINT [ "/app/app" ]