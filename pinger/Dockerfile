FROM golang:latest AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app

FROM alpine:latest
COPY --from=builder /app/app /app/
CMD ["/app/app"]