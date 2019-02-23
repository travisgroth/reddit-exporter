FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN apk add --no-cache git gcc libc-dev ca-certificates
RUN go mod download
RUN go build cmd/reddit-exporter/*.go

FROM alpine
RUN adduser -S -D -H -h /app appuser
RUN apk add --no-cache ca-certificates
USER appuser
COPY --from=builder /build/reddit-exporter /app/
ENTRYPOINT ["/app/reddit-exporter"]
