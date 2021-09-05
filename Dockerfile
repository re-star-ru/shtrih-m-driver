FROM golang:alpine as build
ENV CGO_ENABLED=0

WORKDIR /build
ADD . /build

RUN go build -o app ./examples/client

FROM umputun/baseimage:scratch-latest
ENV TZ=Europe/Moscow
ENV ADDR="0.0.0.0:8080"

COPY --from=build /build/app /srv/app

WORKDIR /srv
CMD ["/srv/app"]