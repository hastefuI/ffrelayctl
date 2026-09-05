FROM golang:1.27-alpine AS builder

WORKDIR /build

ARG VERSION=dev

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -o ffrelayctl .

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /build/ffrelayctl /usr/local/bin/ffrelayctl

USER nonroot:nonroot

ENTRYPOINT ["/usr/local/bin/ffrelayctl"]
