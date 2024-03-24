FROM node:16.14.0-alpine AS frontend
LABEL maintainer="Statping-ng (https://github.com/statping-ng)"
ARG BUILDPLATFORM
WORKDIR /statping
COPY ./frontend/package.json .
COPY ./frontend/yarn.lock .
RUN yarn install --pure-lockfile --network-timeout 1000000
COPY ./frontend .
RUN yarn build && yarn cache clean

# Statping Golang BACKEND building from source
# Creates "/go/bin/statping" and "/usr/local/bin/sass" for copying
FROM golang:1.20.0-alpine AS backend
LABEL maintainer="Statping-NG (https://github.com/statping-ng)"
ARG VERSION
ARG COMMIT
ARG BUILDPLATFORM
ARG TARGETARCH
RUN apk add --no-cache libstdc++ gcc g++ make git autoconf \
    libtool ca-certificates linux-headers wget curl jq && \
    update-ca-certificates

WORKDIR /root
RUN git clone --depth 1 --branch 3.6.2 https://github.com/sass/sassc.git
RUN . sassc/script/bootstrap && make -C sassc -j4
# sassc binary: /root/sassc/bin/sassc

WORKDIR /go/src/github.com/statping-ng/statping-ng
ADD go.mod go.sum ./
RUN go mod download
ENV GO111MODULE on
ENV CGO_ENABLED 1
COPY cmd ./cmd
COPY database ./database
COPY handlers ./handlers
COPY notifiers ./notifiers
COPY source ./source
COPY types ./types
COPY utils ./utils
COPY --from=frontend /statping/dist/ ./source/dist/
RUN go install github.com/GeertJohan/go.rice/rice@latest
RUN cd source && rice embed-go
RUN go build -a -ldflags "-s -w -extldflags -static -X main.VERSION=$VERSION -X main.COMMIT=$COMMIT" -o statping --tags "netgo linux" ./cmd
RUN chmod a+x statping && mv statping /go/bin/statping
# /go/bin/statping - statping binary
# /root/sassc/bin/sassc - sass binary
# /statping - Vue frontend (from frontend)

# Statping main Docker image that contains all required libraries
FROM alpine:latest

RUN apk --no-cache add libgcc libstdc++ ca-certificates curl jq && update-ca-certificates

COPY --from=backend /go/bin/statping /usr/local/bin/
COPY --from=backend /root/sassc/bin/sassc /usr/local/bin/
COPY --from=backend /usr/local/share/ca-certificates /usr/local/share/

WORKDIR /app
VOLUME /app

ENV IS_DOCKER=true
ENV SASS=/usr/local/bin/sassc
ENV STATPING_DIR=/app
ENV PORT=8080
ENV BASE_PATH=""

EXPOSE $PORT

HEALTHCHECK --interval=60s --timeout=10s --retries=3 CMD if [ -z "$BASE_PATH" ]; then HEALTHPATH="/health"; else HEALTHPATH="/$BASE_PATH/health" ; fi && curl -s "http://localhost:${PORT}$HEALTHPATH" | jq -r -e ".online==true"

CMD statping --port $PORT
