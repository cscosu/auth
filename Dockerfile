FROM golang:1.23.2-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY auth_provider.go main.go saml_auth_provider.go tls_certs.go /app
RUN go build -o main .

FROM alpine:latest
WORKDIR /app
COPY static/ /app/static/
COPY --from=0 /app/main main
COPY keys/sp-cert.pem keys/sp-cert.pem
COPY keys/sp-key.pem keys/sp-key.pem
ENV ENV=prod
CMD ["/app/main"]
