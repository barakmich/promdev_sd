FROM golang:1.21 as builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
ADD . .
RUN go build -tags netgo -ldflags="-w -s" ./cmd/promdev_reporter
RUN go build -tags netgo -ldflags="-w -s" ./cmd/promdev_server


FROM scratch

WORKDIR /
COPY --from=builder /src/promdev_reporter /promdev_reporter
COPY --from=builder /src/promdev_server /promdev_server

EXPOSE 9111

CMD ["./promdev_server"]
