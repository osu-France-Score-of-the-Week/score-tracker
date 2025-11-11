FROM golang:1.25.3-alpine3.21

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/mitranim/gow@latest

ENV PATH="/go/bin:${PATH}"

CMD ["go", "run", "src/main.go"]