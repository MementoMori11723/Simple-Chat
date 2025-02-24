FROM golang:1.23.3-alpine
WORKDIR /app
RUN apk add --no-cache gcc musl-dev
ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64
COPY . .
RUN go mod download
CMD ["go", "run", "."]
