FROM golang:alpine as build
ENV CGO_ENABLED=0

ADD . /build
WORKDIR /build

RUN cd app && go build -o /build/kktAPI

FROM umputun/baseimage:scratch-latest
ENV TZ=Europe/Moscow
ENV ADDR="0.0.0.0:8080"

COPY --from=build /build/kktAPI /srv/app

WORKDIR /srv
CMD ["/srv/app"]