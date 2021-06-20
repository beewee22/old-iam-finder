FROM golang:1.15-buster as builder

WORKDIR /tmp/old-iam-finder
COPY . .

RUN go mod tidy \
    && go get -u -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s -w' -o main ./

FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=builder /tmp/old-iam-finder /
CMD ["/main"]
