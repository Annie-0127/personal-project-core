FROM golang:1.23-alpine as builder

ENV APP_ENV development

RUN apk add bash ca-certificates curl git gcc g++ libc-dev unzip make nodejs npm
RUN mkdir -p /personal-project-core


WORKDIR /personal-project-core
ADD . .

RUN npm install -g @redocly/cli@latest
RUN go install github.com/githubnemo/CompileDaemon@latest
RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

ENTRYPOINT CompileDaemon -build="go build -o build/app ./main.go" -command="./build/app"