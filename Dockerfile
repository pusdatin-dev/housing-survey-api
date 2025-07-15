# ---------- Builder Stage ----------
FROM golang:1.21-alpine AS builder

# Install necessary packages
RUN apk update && apk add --no-cache git curl

WORKDIR /app

# Copy Go module files first for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the full source code
COPY . .

# Build the binary
RUN go build -o server ./cmd/

# ---------- Air Dev Image (Optional) ----------
FROM golang:1.21-alpine AS dev

RUN apk add --no-cache git curl

WORKDIR /app

# Install Air
RUN curl -L https://github.com/air-verse/air/releases/download/v1.48.0/air_1.48.0_linux_amd64.tar.gz -o air.tar.gz \
  && tar -xzf air.tar.gz \
  && mv air /usr/local/bin/air \
  && chmod +x /usr/local/bin/air \
  && rm air.tar.gz

COPY . .

CMD ["air"]

# ---------- Final Prod Image ----------
FROM alpine:3.18 AS prod

WORKDIR /app

# Install only required dependencies
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]
