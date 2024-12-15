FROM golang:1.23-alpine AS builder
WORKDIR /app
RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/tg-server

FROM scratch
# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server /server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/server"]