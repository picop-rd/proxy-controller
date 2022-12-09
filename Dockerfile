# syntax=docker/dockerfile:1
FROM golang:1.19-bullseye AS builder

WORKDIR /go/src/github.com/hiroyaonoe/bcop-proxy-controller/
COPY . .

RUN mkdir -p -m 0600 ~/.ssh \
        && ssh-keyscan github.com >> ~/.ssh/known_hosts \
        && git config --global url."git@github.com:".insteadOf "https://github.com/"
RUN --mount=type=ssh CGO_ENABLED=0 go build -o /bcop-proxy-controller ./cmd/controller/main.go


FROM scratch

COPY --from=builder /bcop-proxy-controller /bin/bcop-proxy-controller
ENTRYPOINT [ "/bin/bcop-proxy-controller" ]
CMD [ "--port", "8080", "--mysql", "user:password@tcp(localhost:3306)/db?parseTime=true&collation=utf8mb4_bin" ]

