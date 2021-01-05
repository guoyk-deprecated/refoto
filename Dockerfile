FROM golang:1.14 AS builder
ENV CGO_ENABLED 0
WORKDIR /src
ADD . .
RUN go build -mod vendor -o /refoto

FROM acicn/alpine:3.12
WORKDIR /work
ADD views views
COPY --from=builder /refoto refoto
CMD ["/work/refoto"]