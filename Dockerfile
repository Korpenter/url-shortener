# Building stage for shortener application
FROM golang:1.21-alpine as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app cmd/shortener/main.go

# Final stage
FROM alpine:latest
COPY --from=builder /app/app /app/
CMD ["/app/app"]
