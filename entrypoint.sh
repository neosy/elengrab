#!/bin/sh

# Set ownership of downloads folder to non-root user
chown -R elengrab:elengrab /app_n/downloads
chown -R elengrab:elengrab /app_n/sqlite

# Start cron as root
crond -f &

# exec sleep infinity

# Switch to non-root user and start main app
exec su-exec elengrab /app_n/bin/elengrab
