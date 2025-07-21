FROM golang:1.24 AS builder

WORKDIR /app


COPY . .

RUN go build -o subs .

FROM debian:stable-slim

RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /app

COPY --from=builder /app/subs /app/subs
COPY --from=builder /app/pkg/dbwork/migrations /app/pkg/dbwork/migrations
EXPOSE 8080

CMD ["/app/subs"]
