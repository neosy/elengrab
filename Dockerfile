# ======================
# Builder stage
# ======================
FROM golang:1.25-alpine AS builder

# Set working directory and copy source
WORKDIR /build
COPY . .

# Install git, remove unnecessary files, build binary
RUN apk add --no-cache git \
    && rm -f go.work go.work.sum \
    && CGO_ENABLED=0 go build -o elengrab ./cmd/elengrab/main.go
 
# ======================
# Final image
# ======================
FROM alpine:latest

# Create necessary directories and install dependencies
RUN mkdir -p /app_n/{bin,assets,downloads} \
    && apk add --no-cache su-exec curl dcron python3 py3-pip ffmpeg \
    # Download yt-dlp binary
    && curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp \
    && chmod a+rx /usr/local/bin/yt-dlp \
    # Remove curl to keep image small
    && apk del curl \
    # Add cron job to update yt-dlp once a day
    && echo "0 1 * * * /usr/local/bin/yt-dlp -U >> /var/log/yt-dlp-update.log 2>&1" >> /etc/crontabs/root

# Copy application binaries and assets from builder stage
COPY --from=builder /build/elengrab /app_n/bin/elengrab
COPY entrypoint.sh /app_n/entrypoint.sh
COPY internal/api/htmx/assets /app_n/assets/

# Create a non-root user and set ownership
RUN adduser -D -h /app_n elengrab \
    && chown -R elengrab:elengrab /app_n \
    && chown root:root /app_n/entrypoint.sh \
    && chmod +x /app_n/entrypoint.sh

# Set working directory
WORKDIR /app_n

# Declare downloads folder as a volume
VOLUME ["/app_n/downloads"]

# Expose web server port
EXPOSE 8080

# Start cron in background and run main app
# ENTRYPOINT ["/bin/sh", "-c", "crond -f & exec /app_n/bin/elengrab"]
ENTRYPOINT ["/app_n/entrypoint.sh"]
