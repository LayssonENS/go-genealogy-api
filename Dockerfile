FROM golang:1.19.5-alpine

WORKDIR /usr/src/app

RUN go install github.com/cosmtrek/air@latest

COPY . .

# Execute o "go mod tidy"
RUN go mod tidy \
    && rm -rf /var/cache/apk/*