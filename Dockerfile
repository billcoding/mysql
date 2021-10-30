FROM golang:latest as builder
WORKDIR /go/src/app
COPY . .
ENV GOPROXY "https://goproxy.cn,direct"
RUN go mod tidy
RUN go build -o mysql
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/app/mysql /mysql
CMD ["./mysql"]