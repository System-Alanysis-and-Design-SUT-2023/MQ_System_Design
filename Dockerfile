FROM alpine:3.19 AS root-certs
RUN apk add -U --no-cache ca-certificates
RUN addgroup -g 1001 app
RUN adduser app -D -u 1001 -G app /home/app

FROM golang:1.21 AS builder
WORKDIR /go/src/server
ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GGOOS=linux
ENV GOARCH=arm64
COPY --from=root-certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY ../.. .
RUN go build -o /go/bin/app cmd/main.go

FROM scratch AS final
COPY --from=root-certs /etc/passwd /etc/passwd
COPY --from=root-certs /etc/group /etc/group
COPY --chown=1001:1001 --from=root-certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --chown=1001:1001 --from=builder /go/bin/app /opt/server/
COPY --chown=1001:1001 app.env /
USER app
ENTRYPOINT [ "/opt/server/app" ]