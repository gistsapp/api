FROM golang:1.22.5-alpine3.20 AS builder
WORKDIR /app
COPY . /app
RUN go build -o main .

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main /app/main
RUN chmod +x /app/main
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

CMD ["/app/main"]
