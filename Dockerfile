FROM golang:1.21 AS builder
WORKDIR /go/app
COPY . .
RUN go build -o bin cmd/main.go

FROM alpine:3.19
WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/app/bin .
EXPOSE 8080
CMD [ "/app/bin" ]