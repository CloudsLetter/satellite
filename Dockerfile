FROM golang:1.19-alpine3.16 AS builder
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