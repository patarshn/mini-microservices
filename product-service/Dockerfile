FROM golang:1.21
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
COPY .env.docker ../.env
RUN go build -o /product-service ./cmd/main.go
EXPOSE 8082
CMD ["/product-service"]