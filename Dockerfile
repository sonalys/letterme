FROM golang:1.17 as build

WORKDIR /src
COPY ./ .

ENV CGO_ENABLED=0
ARG TARGETOS
ARG TARGETARCH
RUN go build -o bin -mod vendor -ldflags '-w' .

FROM alpine:3.12

WORKDIR /app

COPY --from=build /src/bin .

CMD /app/bin