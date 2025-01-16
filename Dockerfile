FROM --platform=$BUILDPLATFORM tonistiigi/xx:1.6.1 AS xx


FROM --platform=$BUILDPLATFORM node:16.14.0-alpine AS frontend
RUN mkdir /build
WORKDIR /build
COPY ./frontend/package.json .
COPY ./frontend/yarn.lock .
RUN yarn install --pure-lockfile --network-timeout 1000000
COPY ./frontend .
RUN yarn build && yarn cache clean


FROM --platform=$BUILDPLATFORM golang:1.20.0-alpine AS backend
COPY --from=xx / /
RUN apk add --no-cache clang lld
ARG TARGETPLATFORM
RUN xx-apk add --no-cache musl-dev gcc

RUN mkdir /build
WORKDIR /build

ADD go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY database ./database
COPY handlers ./handlers
COPY notifiers ./notifiers
COPY source ./source
COPY types ./types
COPY utils ./utils

COPY --from=frontend /build/dist/ ./source/dist/
RUN go install github.com/GeertJohan/go.rice/rice@latest
RUN cd source && rice embed-go

ENV CGO_ENABLED=1
ARG VERSION
ARG COMMIT
RUN xx-go build -a -ldflags "-s -w -extldflags -static -X main.VERSION=$VERSION -X main.COMMIT=$COMMIT" -o statping --tags "netgo linux" ./cmd
RUN xx-verify --static statping


# Statping main Docker image that contains all required libraries
FROM alpine:latest
LABEL maintainer="Statping-NG (https://github.com/statping-ng)"

RUN apk --no-cache add libgcc libstdc++ ca-certificates curl jq sassc && update-ca-certificates

COPY --from=backend /build/statping /usr/local/bin/
COPY --from=backend /usr/local/share/ca-certificates /usr/local/share/

WORKDIR /app
VOLUME /app

ENV IS_DOCKER=true
ENV SASS=/usr/bin/sassc
ENV STATPING_DIR=/app
ENV PORT=8080
ENV BASE_PATH=""

EXPOSE $PORT
HEALTHCHECK --interval=60s --timeout=10s --retries=3 CMD if [ -z "$BASE_PATH" ]; then HEALTHPATH="/health"; else HEALTHPATH="/$BASE_PATH/health" ; fi && curl -s "http://localhost:${PORT}$HEALTHPATH" | jq -r -e ".online==true"

CMD statping --port $PORT
