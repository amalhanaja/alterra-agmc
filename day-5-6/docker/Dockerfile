# Build executable binary
FROM golang:1.19.1-alpine3.16 AS builder
LABEL maintainer="Alfian Akmal Hanantio<amalhanaja@gmail.com>"
RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base
WORKDIR /app
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
COPY . .
RUN go build -o /alterra-agmc cmd/restapi/main.go

# build small iamge
FROM alpine:3.16.0
WORKDIR /app
COPY --from=builder alterra-agmc .
EXPOSE 8080
CMD [ "./alterra-agmc" ]
