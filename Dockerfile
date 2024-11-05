# stage: build ---------------------------------------------------------

FROM golang:1.22-alpine as build

RUN apk add --no-cache gcc musl-dev linux-headers

WORKDIR /go/src/github.com/flashbots/slogth

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -o bin/slogth -ldflags "-s -w" github.com/flashbots/slogth/cmd

# stage: run -----------------------------------------------------------

FROM alpine

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=build /go/src/github.com/flashbots/slogth/bin/slogth ./slogth

ENTRYPOINT ["/app/slogth"]
