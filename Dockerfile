# Build Stage
FROM golang:1.19.2-alpine3.16 AS builder

RUN apk update && \
    apk upgrade && \
    apk --no-cache add ca-certificates build-base gcc pkgconf && \
    rm -rf /var/cache/apk/*

WORKDIR /app

COPY . .

RUN go build -o main main.go


# Run Stage
FROM alpine:3.16

WORKDIR /app

COPY --from=builder /app/main .

COPY app.env .

EXPOSE 8080

CMD /app/main