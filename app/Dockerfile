FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./main.go

FROM alpine:3.20
WORKDIR /app
# Create non-root user (uid/gid 65532 similar to distroless nonroot)
RUN addgroup -g 65532 appgroup \
    && adduser -D -H -u 65532 -G appgroup appuser
COPY --from=builder /app/server /app/server
EXPOSE 8000
USER 65532:65532
ENTRYPOINT ["/app/server"]


