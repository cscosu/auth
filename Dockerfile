FROM golang:1.23.2-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY templates templates
COPY migrations migrations
COPY *.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
COPY --from=0 /usr/local/go/lib/time/zoneinfo.zip /
ENV ZONEINFO=/zoneinfo.zip
WORKDIR /app
COPY static/ /app/static/
COPY --from=0 /app/main main
COPY keys/sp-cert.pem keys/sp-cert.pem
COPY keys/sp-key.pem keys/sp-key.pem
COPY .env .env
ENV ENV=prod
CMD ["/app/main"]
