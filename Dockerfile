FROM golang:1.22.4 AS builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app
COPY . .

RUN go mod tidy

WORKDIR /app/cmd
RUN CGO_ENABLED=0 GOOS=linux go build  -ldflags="-w -s" -o ../main


WORKDIR /app
RUN mkdir publish  \
    && cp main publish  \
    && cp -r conf publish

FROM busybox:1.28.4

WORKDIR /app
COPY --from=builder /app/publish .

ENV GIN_MODE=release
EXPOSE 5001

ENTRYPOINT ["./main"]
