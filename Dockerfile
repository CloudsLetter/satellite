FROM golang:1.23.4-alpine3.21 AS builder
ADD . /app
WORKDIR /app
RUN go mod tidy
RUN go build -o main .

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/main /app/main
WORKDIR /app
RUN chmod +x main
EXPOSE 8080
CMD ["./main"]