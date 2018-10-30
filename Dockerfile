ARG GO_IMPORT_PATH=github.com/thelegy/koma-bot

FROM golang:1.11-alpine AS builder

RUN apk add --no-cache \
    bash \
    git \
    npm

RUN npm install --global \
    autoprefixer \
    postcss-cli \
    sass

ARG GO_IMPORT_PATH
WORKDIR src/$GO_IMPORT_PATH
COPY ./ ./

RUN SASS_STYLE=compressed ./processSass.sh

RUN mkdir _build \
 && cp -R templates _build/templates \
 && cp -R static _build/static \
 && rm _build/static/**/*.scss

RUN go get -v -d ./...

RUN CGO_ENABLED=0 GOOS=linux go build \
      -a -v -o _build/koma-bot -tags netgo \
      -ldflags " \
        -s \
        -X main._versionDate=$(date -Iseconds) \
        -X main._versionGitDirty=$(git diff-index --quiet HEAD -- && echo f || echo t) \
        -X main._versionGitBranch=$(git symbolic-ref HEAD 2>/dev/null|sed 's|^refs/heads/||'||true) \
        -X main._versionGitHash=$(git rev-parse --short --verify HEAD) \
      " .

FROM alpine

MAINTAINER The Legy

WORKDIR /opt/koma-bot/

EXPOSE 8000

COPY koma_bot.example.yaml koma_bot_sounds.example.yaml /etc/koma_bot/
COPY ca-certificates.crt /etc/ssl/certs/

VOLUME "/etc/koma_bot"

ARG GO_IMPORT_PATH
COPY --from=builder /go/src/$GO_IMPORT_PATH/_build/ ./

RUN mv koma-bot /bin/

USER nobody

ENTRYPOINT ["/bin/koma-bot", "--docker"]
