# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o go_translate ./cmd/.

# Final image
FROM alpine:3.19

WORKDIR /app

# Copy Go binary
COPY --from=builder /app/go_translate .

# Copy configs and env
COPY configs ./configs
COPY .env . 

# Install Node.js + npm + migrate-mongo
RUN apk add --no-cache nodejs npm bash curl
RUN npm install -g migrate-mongo

# Copy migrate-mongo config and scripts
COPY ./migrations ./migrations
COPY ./migrate-mongo-config.js .
COPY ./entrypoint.sh .
RUN chmod +x ./entrypoint.sh

EXPOSE 8080

CMD ["./go_translate"]
