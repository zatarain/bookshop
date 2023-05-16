FROM golang:1.20.4
WORKDIR /api
COPY . .
RUN go mod tidy
RUN go test -v ./...
CMD ["go", "run", "main.go"]
