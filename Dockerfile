FROM golang:1.25.0-alpine
WORKDIR /app
COPY . .
RUN go mod download
CMD ["go", "run", "cmd/server/main.go"]

