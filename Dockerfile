FROM ghcr.io/umputun/baseimage/buildgo:latest as build
WORKDIR /build
ADD . /build

RUN \
    revision=$(/script/git-rev.sh) && \
    echo "revision=${revision}" && \
    go build -o app -ldflags "-X main.revision=$revision -s -w" ./examples/client

FROM ghcr.io/umputun/baseimage/app:lastest

COPY --from=build /build/app /srv/app

EXPOSE 8080
WORKDIR /srv
ENV LISTEN = "0.0.0.0:8080"

CMD ["/srv/app"]





