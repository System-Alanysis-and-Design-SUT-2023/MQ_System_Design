FROM golang:1.21 AS builder
WORKDIR /go/src/server
ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GGOOS=linux
COPY ../.. .
RUN go build -o /go/bin/app cmd/main.go

FROM debian:buster-slim
COPY --from=builder /go/bin/app /opt/server/
COPY app.env /
ENTRYPOINT [ "/opt/server/app" ]