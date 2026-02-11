FROM golang:1.24-alpine AS builder

WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /dtu-discourse ./cmd/dtu-discourse

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=builder /dtu-discourse /usr/local/bin/dtu-discourse
EXPOSE 4200
ENV PORT=4200
ENTRYPOINT ["dtu-discourse"]
