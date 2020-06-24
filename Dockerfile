FROM golang as builder

WORKDIR /go/src/github.com/kotaroooo0/zyankenbot

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o main

FROM scratch

COPY .env /
EXPOSE 3000
COPY --from=builder /go/src/github.com/kotaroooo0/zyankenbot/main /main

ENTRYPOINT ["/main"]
