# Build first
FROM golang:1.17-alpine AS builder
RUN apk add --no-cache git musl-dev
COPY . /opt
WORKDIR /opt
RUN go build -v -o bin/matrix-room-directory-server

# The actual image (which is lightweight)
FROM alpine
COPY --from=builder /opt/bin/matrix-room-directory-server /usr/local/bin/
RUN apk add --no-cache \
        su-exec \
        ca-certificates
ENTRYPOINT "/usr/local/bin/matrix-room-directory-server"
EXPOSE 8080
