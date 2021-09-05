FROM golang:alpine as build

WORKDIR /build
ADD . /build

RUN go build -o app ./examples/client

FROM golang:alpine

COPY --from=build /build/app /srv/app

WORKDIR /srv

ENV ADDR = "0.0.0.0:8080"

CMD ["./app"]





