FROM golang:1.22.5-alpine3.20 AS builder
WORKDIR /app
COPY . /app
RUN go build -o main .

FROM alpine:3.20
ENV PORT=4000
WORKDIR /app
COPY --from=builder /app/main /app/main
COPY --from=builder /app/docs /app/docs
RUN chmod +x /app/main
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

EXPOSE 4000

CMD ["/app/main"]
