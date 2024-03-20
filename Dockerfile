FROM golang:1.18-stretch as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .
# FROM debian:stretch-slim
FROM gcr.io/distroless/static-debian11

COPY --from=builder /app/main /
CMD ["/main"]