FROM golang:alpine as build

WORKDIR /build
ADD . /build

RUN go build -o app ./examples/client

FROM golang:alpine

COPY --from=build /build/app /srv/app

WORKDIR /srv

EXPOSE 8080

ENV LISTEN = "0.0.0.0:8080"

CMD ["./app"]





