FROM golang:1.20.4

ENV GOMOD=/api/go.mod

WORKDIR /api
COPY . .

RUN go install github.com/joho/godotenv/cmd/godotenv@latest
RUN go mod tidy
RUN ENVIRONMENT=test godotenv -f "${ENVIRONMENT}.env" go test -v ./...

CMD godotenv -f ${ENVIRONMENT}.env go run main.go
