FROM golang:alpine as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    APP_PORT=3001 \
    MYSQL_HOST=mysql \
    MYSQL_DATABASE=transaction \
    MYSQL_USER=dev \
    MYSQL_PASSWORD=dev \
    MYSQL_PORT=3306

WORKDIR /build
COPY . .
RUN go mod download && go build -a --installsuffix cgo --ldflags="-s" -o main

FROM scratch

COPY --from=builder /build .

ENTRYPOINT ["./main"]

EXPOSE 3001