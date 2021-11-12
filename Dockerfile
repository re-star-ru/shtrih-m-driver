FROM golang:alpine as build
ENV CGO_ENABLED=0

ADD . /build
WORKDIR /build

RUN apk add --no-cache --update git

RUN \
  version=$(git rev-parse --abbrev-ref HEAD)-$(git log -1 --format=%h)-$(date +%Y%m%dT%H:%M:%S) && \
  echo "version=$version" && \
  cd app && go build -ldflags "-X main.version=${version} -s -w" -o /build/kktAPI

FROM umputun/baseimage:scratch-latest
ENV TZ=Europe/Moscow
ENV ADDR="0.0.0.0:8080"

COPY --from=build /build/kktAPI /srv/app

WORKDIR /srv
CMD ["/srv/app"]