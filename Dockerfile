# syntax=docker/dockerfile:1

ARG GOVERSION=1.26.5

FROM golang:${GOVERSION}-alpine AS dev
RUN go install "github.com/air-verse/air@latest" && \
    go install "github.com/a-h/templ/cmd/templ@latest"
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download && go mod verify
CMD ["air", "-c", ".air.toml"]


FROM --platform=$BUILDPLATFORM golang:${GOVERSION}-alpine AS build-prod
ARG TARGETOS
ARG TARGETARCH
WORKDIR /src
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /bin/server ./cmd/server
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /bin/manage ./cmd/manage

FROM alpine AS prod
WORKDIR /src
COPY --from=build-prod /usr/local/go/lib/time/zoneinfo.zip /
ENV ZONEINFO=/zoneinfo.zip
COPY --from=build-prod /bin/server /bin/manage /bin/
COPY . .
EXPOSE 80
ENV DEBUG=false
ARG SOURCE_COMMIT
ENV SOURCE_COMMIT=${SOURCE_COMMIT}
ENV SERVER_PORT=80
# Migrations are embedded in the manage binary (db/migrate.go); apply them, then serve.
ENTRYPOINT ["/bin/sh", "-c", "manage -c migrate-up && server"]
