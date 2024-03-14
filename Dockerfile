FROM golang:latest AS GO_BUILD
COPY . .
CMD go run ./cmd/main.go
ENV GOPATH=/



