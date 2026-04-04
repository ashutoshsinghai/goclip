# ---- Build stage ----
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o goclip .

# ---- Runtime stage ----
FROM alpine:3.19

RUN apk add --no-cache xclip

WORKDIR /root

COPY --from=builder /app/goclip /usr/local/bin/goclip

# History is stored in ~/.goclip — mount a volume to persist it
VOLUME ["/root/.goclip"]

ENTRYPOINT ["goclip"]
CMD ["help"]
