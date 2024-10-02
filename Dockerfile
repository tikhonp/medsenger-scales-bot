# syntax=docker/dockerfile:1

ARG GO_VERSION=1.23.1

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build

WORKDIR /src

RUN go install "github.com/air-verse/air@latest"
RUN go install "github.com/pressly/goose/v3/cmd/goose@latest"

ADD --chmod=111 'https://github.com/apple/pkl/releases/download/0.26.3/pkl-alpine-linux-amd64' /bin/pkl

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . ./

ARG TARGETARCH
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/get_db_string ./cmd/get_db_string/

CMD ["air"]