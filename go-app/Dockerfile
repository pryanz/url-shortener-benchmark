FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener-benchmark .

FROM alpine:3.19

WORKDIR /root/

COPY --from=builder /app/url-shortener-benchmark .

EXPOSE 8081

CMD ["./url-shortener-benchmark"]