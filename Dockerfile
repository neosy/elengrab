# ======================
# Builder stage: compile Go application
# ======================
FROM golang:1.25-alpine AS builder

# Set working directory and copy source
WORKDIR /build
COPY . .

# Install git, remove unnecessary files, build binary
RUN apk add --no-cache git \
    && rm -f go.work go.work.sum \
    && CGO_ENABLED=0 go build -o elengrab ./cmd/elengrab

# ======================
# Final image
# ======================
FROM alpine:latest

ARG APP_DIR=/app_n

RUN apk add --no-cache tzdata su-exec curl dcron python3 ffmpeg \
    # Create necessary directories and install dependencies
    && mkdir ${APP_DIR} \
    && cd ${APP_DIR} \
        && mkdir -p bin assets downloads sqlite/data sqlite/backups \
    # Download yt-dlp binary
    && curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp \
    && chmod a+rx /usr/local/bin/yt-dlp \
    # time Zone
    && ln -sf /usr/share/zoneinfo/Europe/Moscow /etc/localtime \
    # Remove packages to keep image small
    && apk del curl \
    && rm -rf /var/cache/apk/* \
    # Add cron job to update yt-dlp once a day
    && echo "0 1 * * * /usr/local/bin/yt-dlp -U >> /var/log/yt-dlp-update.log 2>&1" >> /etc/crontabs/root

# Copy the compiled Go binary from the builder stage
COPY --from=builder /build/elengrab ${APP_DIR}/bin/elengrab
COPY entrypoint.sh ${APP_DIR}/entrypoint.sh
COPY internal/api/rest/server/assets ${APP_DIR}/assets/

# Create a non-root user and set ownership
RUN adduser -D -h ${APP_DIR} elengrab \
    && chown -R elengrab:elengrab ${APP_DIR} \
    && chown root:root ${APP_DIR}/entrypoint.sh \
    && chmod +x ${APP_DIR}/entrypoint.sh

# Set working directory
WORKDIR ${APP_DIR}

# Declare downloads folder as a volume
VOLUME ["/app_n/downloads"]
VOLUME ["/app_n/sqlite/data"]
VOLUME ["/app_n/sqlite/backups"]

# Enviroments
ENV TZ=Europe/Moscow
ENV ELENGRAB_APP_DIR=${APP_DIR}

# Expose web server port
EXPOSE 8080

# Start cron in background and run main app
# ENTRYPOINT ["/bin/sh", "-c", "crond -f & exec /app_n/bin/elengrab"]
ENTRYPOINT ["/app_n/entrypoint.sh"]
