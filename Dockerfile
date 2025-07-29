# syntax=docker/dockerfile:1

ARG GOVERSION=1.24.5

FROM golang:${GOVERSION}-alpine AS dev
RUN go install "github.com/air-verse/air@latest" && \
    go install "github.com/a-h/templ/cmd/templ@latest" &&\
    go install "github.com/pressly/goose/v3/cmd/goose@latest"
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download && go mod verify
CMD ["air", "-c", ".air.toml"]

FROM golang:${GOVERSION}-alpine AS prod
RUN go install "github.com/pressly/goose/v3/cmd/goose@latest"
ARG TARGETARCH
WORKDIR /src
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    GOARCH=$TARGETARCH go build -o /bin/server ./cmd/server
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    GOARCH=$TARGETARCH go build -o /bin/get_db_string ./cmd/get_db_string/
COPY . .
