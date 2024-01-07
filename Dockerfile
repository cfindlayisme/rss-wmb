FROM golang:1.21.3 AS builder

WORKDIR /app

COPY . ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /application

FROM scratch

COPY --from=builder /application /application

# Run
CMD ["/application"]