FROM golang:1.22-bullseye

WORKDIR /app

COPY go.mod go.sum /app/
COPY ./src /app/src

WORKDIR /app/src/app/main/cmd

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main .

EXPOSE 8000

CMD ["/main"]
