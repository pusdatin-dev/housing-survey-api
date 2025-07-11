FROM golang:1.21-alpine

# Install required packages
RUN apk update && apk add --no-cache git curl unzip

WORKDIR /app

# Install air (hot reload)
# RUN go install github.com/air-verse/air@v1.48.0
RUN curl -L https://github.com/air-verse/air/releases/download/v1.48.0/air_1.48.0_linux_amd64.tar.gz -o air.tar.gz \
  && tar -xzf air.tar.gz \
  && mv air /usr/local/bin/air \
  && chmod +x /usr/local/bin/air \
  && rm air.tar.gz

# Copy Go mod files and download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy all source
COPY . .

# Run air on container start
CMD ["air"]
