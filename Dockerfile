FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

CMD CGO_ENABLED=0 go test -tags integration -v ./...