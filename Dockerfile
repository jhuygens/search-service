FROM golang:alpine as builder

RUN mkdir /build
ADD ./ /build 
WORKDIR /build
RUN env GOOS=linux GOARCH=386 go build -o main .

FROM alpine:latest

RUN mkdir -p /app && adduser -S -D -H -h /app miku && chown -R miku /app
COPY --from=builder /build/main /app/
USER miku
EXPOSE 9091
WORKDIR /app
CMD ["./main"]