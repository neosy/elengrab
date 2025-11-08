#!/bin/sh

# Set ownership of downloads folder to non-root user
chown -R elengrab:elengrab /app_n/downloads

# Start cron as root
crond -f &

# Switch to non-root user and start main app
exec su-exec elengrab /app_n/bin/elengrab
