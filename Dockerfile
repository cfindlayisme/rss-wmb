FROM golang:1.21.3 AS builder

WORKDIR /app

COPY . ./
RUN go mod download

RUN apt-get update && apt-get install -y ca-certificates

RUN CGO_ENABLED=0 GOOS=linux go build -o /application

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /application /application

# Run
CMD ["/application"]