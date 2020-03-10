FROM golang:stretch AS builder
ENV GO111MODULE=on
WORKDIR /go/src/
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
run CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app .
FROM alpine:3.8
COPY --from=builder /go/src/app /app
WorkDir / 
CMD ["/app"]


