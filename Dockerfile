# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR /seyna/emailService

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /seyna-email-service 

CMD [ "/seyna-email-service" ]