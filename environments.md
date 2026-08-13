# Environment Variables

## Config

Config contains the application configuration.

 - `ELENGRAB_APP_ENV` (default: `production`) - AppEnv defines the application mode: local, develop, test, or production.
 - `ELENGRAB_LOG_FORMAT` (default: `console`) - LogFormat specifies the output format for logs: "json" or "console".
 - `ELENGRAB_LOG_LEVEL` (default: `warn`) - LogLevel sets the minimum log level: "debug", "info", "warn", or "error".
 - `ELENGRAB_ADMIN_SERVER_ENABLE` (default: `false`) - Enables the admin server.
 - `ELENGRAB_ADMIN_SERVER_ADDRESS` (default: `127.0.0.1`) - Address on which the admin server listens.
 - `ELENGRAB_ADMIN_SERVER_PORT` (default: `6060`) - Port on which the admin server listens.
 - `ELENGRAB_ADMIN_SERVER_DEBUG_PPROF` (default: `true`) - Enables the pprof debugging endpoints.
 - `ELENGRAB_ADMIN_SERVER_DEBUG_METRICS` (default: `true`) - Enables the metrics endpoint.
 - `ELENGRAB_ADMIN_SERVER_DEBUG_HEALTH` (default: `true`) - Enables the health check endpoint.
 - `ELENGRAB_MODE` (default: `public`) - Application operating mode.
 - `ELENGRAB_DEMO_MODE` (default: `false`) - Enables demo mode with restricted functionality.
 - `ELENGRAB_BASE_URL` - Base URL used to build absolute application links.
 - `ELENGRAB_SHORT_LINK_PREFIX` (default: `/s`) - Prefix used for generated short links.
 - `ELENGRAB_SHORT_LINK_LENGTH` (default: `6`) - Length of generated short link identifiers.
 - `ELENGRAB_SHORT_LINK_TTL_DAYS` (default: `180`) - Lifetime of generated short links in days.
 - `ELENGRAB_ADMIN_LOGIN` - AdminLogin is used to create the default administrator at first startup
 - `ELENGRAB_ADMIN_PASSWORD` - AdminPassword is used to create the default administrator at first startup
 - `ELENGRAB_DOWNLOADER_BIN_DIR` (default: `/usr/local/bin`) - Directory containing downloader binaries such as yt-dlp.
 - `ELENGRAB_ROOT_DIR` - Root directory for application data.
 - `ELENGRAB_ASSETS_DIR` (default: `assets`) - Directory containing application assets.
 - `ELENGRAB_MEDIA_DIR` (default: `media`) - Directory containing media files (e.g. thumbnails).
 - `ELENGRAB_DOWNLOADS_DIR` (default: `downloads`) - Directory containing downloaded files.
 - `ELENGRAB_COOKIES_DIR` (default: `cookies`) - Directory containing cookies used by the downloader.
 - `ELENGRAB_DOWNLOAD_WORKERS` (default: `3`) - Maximum number of concurrent media downloads.
 - `ELENGRAB_OPERATION_WORKERS` (default: `5`) - Maximum number of concurrent background operations.
 - `ELENGRAB_DELETE_DUPLICATES_UNIQUENESS_SCOPE` (default: `per_user`) - Scope used when checking media uniqueness before deleting duplicates.
 - `ELENGRAB_ALLOW_COOKIES` (default: `false`) - Enables the use of cookies when downloading media.
 - `ELENGRAB_MAINTENANCE_UPDATE_HASH_INTERVAL` (default: `8h`) - Interval for calculating and updating hashes of downloaded files.
File hashes are used to identify duplicate files.
 - `ELENGRAB_MAINTENANCE_DELETE_DUPLICATES_INTERVAL` (default: `1h`) - Interval for deleting duplicate media files.
 - `ELENGRAB_MAINTENANCE_DELETE_MISSING_DOWNLOADS_INTERVAL` (default: `12h`) - Interval for deleting database records whose downloaded files are missing.
 - `ELENGRAB_MAINTENANCE_DELETE_FAILED_DOWNLOADS_INTERVAL` (default: `1h`) - Interval for deleting failed download records.
 - `ELENGRAB_MAINTENANCE_ENABLE_MOVE_UNMATCHED_FILES` (default: `false`) - Enables moving unmatched files during periodic maintenance.
 - `ELENGRAB_MAINTENANCE_DATABASE_BACKUPS_KEEP` (default: `7`) - Number of the latest database backups to keep.
If set to 0, old backup files are not deleted.
 - `ELENGRAB_HTTP_SERVER_ADDRESS` - Address specifies the network address on which the HTTP server listens.
 - `ELENGRAB_HTTP_SERVER_PORT` (default: `8080`) - Port specifies the port on which the HTTP server listens.
 - `ELENGRAB_HTTP_SERVER_COMPRESS` (default: `true`) - Compress enables HTTP response compression.
 - `ELENGRAB_SQLITE_DATA_DIR` (default: `sqlite/data`) - Directory where SQLite database files are stored.
 - `ELENGRAB_SQLITE_BACKUPS_DIR` (default: `sqlite/backups`) - Directory where SQLite database backups are stored.

