FROM golang:1.23 as builder

WORKDIR /Forum01

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o forum ./prg/web/* 

FROM golang:1.23-alpine

RUN apk update && apk upgrade

RUN apk add --no-cache sqlite

WORKDIR /Work

COPY --from=builder /Forum01/forum .

EXPOSE  8080

ENTRYPOINT [ "./forum" ]