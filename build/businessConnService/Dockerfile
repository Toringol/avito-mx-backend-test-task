FROM golang:latest

ENV GO111MODULE=on

WORKDIR /avitoservice

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build /avitoservice/cmd/businessConnService/businessConnService.go

EXPOSE 8080

ENTRYPOINT ["/avitoservice/businessConnService"]
