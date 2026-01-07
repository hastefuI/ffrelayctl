FROM golang:1.23-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -a \
    -installsuffix cgo \
    -ldflags='-w -s -extldflags "-static"' \
    -trimpath \
    -o ffrelayctl .

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /build/ffrelayctl /usr/local/bin/ffrelayctl

USER nonroot:nonroot

ENTRYPOINT ["/usr/local/bin/ffrelayctl"]
