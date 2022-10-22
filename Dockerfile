FROM golang:1.19.2-alpine3.16

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /go-rest-api

EXPOSE 8080

CMD [ "/go-rest-api" ]
