FROM golang:1.23-alpine

## install sqlite
RUN apk add --no-cache gcc musl-dev sqlite sqlite-dev

WORKDIR /Forum01

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

EXPOSE  8228

ENTRYPOINT [ "./main" ]