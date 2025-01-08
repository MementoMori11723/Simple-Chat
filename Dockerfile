FROM golang:1.23.3-alpine
WORKDIR /app
COPY . .
EXPOSE 11000
CMD ["go", "run", "."]
