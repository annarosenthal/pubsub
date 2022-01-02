FROM golang:1.17-alpine AS build
COPY ./ /src
WORKDIR /src
ENV CGO_ENABLED=0
RUN go build -o pubsub ./cmd/pubsub
RUN go test ./pkg/pubsub

FROM alpine
COPY --from=build /src/pubsub /usr/local/bin/pubsub
ENTRYPOINT ["/usr/local/bin/pubsub"]
