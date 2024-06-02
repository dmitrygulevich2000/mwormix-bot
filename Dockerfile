FROM golang:1.21

WORKDIR /mwormix-bot
SHELL ["/bin/bash", "-c"]

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/
RUN go build -o ./bin/ ./cmd/notify_once/

COPY .env-docker ./

CMD source .env-docker && ./bin/notify_once
