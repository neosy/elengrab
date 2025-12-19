<p align="center">
  <img width="192" height="192" alt="android-chrome-192x192_round" src="https://github.com/user-attachments/assets/f2973dcc-90d3-4046-b9e0-fd34b6697fa5" />
</p>
<h3 align="center">Self-hosted web interface for working with YouTube media.</h3>

<p align="center">
  <img width="1248" height="517" alt="Screenshot_126" src="https://github.com/user-attachments/assets/e6f9cfa6-2a26-4330-8897-ca32c5e97b88" />
</p>

# Elengrab

**Fast web interface for your own YouTube video/audio processing server with format and quality options. The project integrates with existing open-source media processing utilities
(such as [yt-dlp](https://github.com/yt-dlp/yt-dlp)) as an optional backend component.**

Elengrab provides a simple and **very fast** web interface to work with YouTube videos and audio tracks, allowing selection of formats and quality settings. The project is fully written in **Go**, and the frontend is powered by **HTMX**, ensuring high responsiveness and minimal overhead. It runs in a lightweight Docker container, is easy to deploy, and serves as a self-hosted frontend for managing your personal media server.

---

## Features

- Support for YouTube video content in various formats.
- Support for audio tracks from YouTube content.
- Configurable format and quality settings.
- Task queue for media processing.
- Instant addition/removal from queue.
- Concurrent processing of multiple tasks (3 by default).
- Dark theme.

---

## YouTube Content Context

Elengrab is designed to work with media hosted on YouTube.
It operates as a self-hosted interface for organizing and processing personal video and audio content.
Users are responsible for ensuring compliance with applicable laws and platform terms.

---

## Quick Start
### Run with default settings

```
docker run -d \
  --name elengrab \
  -v elengrab_downloads:/app_n/downloads \
  -p 8080:8080 \
  neosy/elengrab:latest
```

---

## Requirements

### Minimum

* **CPU:** 1–2 cores
* **Memory:** ~1 GB RAM
* **Concurrent downloads:** 1 worker

This configuration is suitable for low-resource servers. To limit concurrency, set:

```
ELENGRAB_DOWNLOAD_WORKERS=1
```

### Recommended

* **CPU:** 4-6 cores
* **Memory:** 4 GB RAM
* **Concurrent downloads:** 3 workers (default)

By default, the Docker container is configured to use this setup.

### Notes

The main resource consumers are **yt-dlp** and **ffmpeg**, especially during video downloading, merging, and transcoding.

The number of concurrent workers can be adjusted using the following environment variable:

| Variable                  | Default | Description                                                                                |
| ------------------------- | ------- | ------------------------------------------------------------------------------------------ |
| ELENGRAB_DOWNLOAD_WORKERS | 3       | Number of concurrent workers used for processing YouTube video and audio tasks in parallel |


---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `LOG_LEVEL` | `warn` | Logging level. Options: `debug`, `info`, `warn`, `error`. |
| `SQLITE_DATA_DIR` | `/app_n/sqlite/data` | Directory where SQLite database files are stored. |
| `SQLITE_BACKUPS_DIR` | `/app_n/sqlite/backups` | Directory where SQLite backup files are stored. |
| `ELENGRAB_DOWNLOADER_BIN_DIR` | `/usr/local/bin` | Directory containing yt-dlp binary. |
| `ELENGRAB_ASSETS_DIR` | `/app_n/assets` | Directory containing application assets. |
| `ELENGRAB_DOWNLOADS_DIR` | `/app_n/downloads` | Directory where downloaded files are stored inside the container. Must be mapped to a host volume. |
| `ELENGRAB_DOWNLOAD_WORKERS` | 3 | Number of concurrent workers used for processing YouTube video and audio tasks in parallel. |
| `ELENGRAB_LOAD_HISTORY` | `true` | Whether to display the download history of files. Set to true to show the history, or false to hide it. |
| `ELENGRAB_MAINTENANCE_ENABLE_MOVE_UNMATCHED_FILES` | `false`   | Enables the periodic operation that moves files not present in the database tables from the download folder to the `.lost` folder. Default is `false` (disabled). |

---

## Usage

### Run with default settings

```
docker run -d \
  --name elengrab \
  -v elengrab_db:/app_n/sqlite/data \
  -v elengrab_db_backups:/app_n/sqlite/backups \
  -v elengrab_downloads:/app_n/downloads \
  -p 8080:8080 \
  neosy/elengrab:latest
```

---

### Docker Compose

Create a file `docker-compose.yml` with the following content:

```
version: "3.8"

services:
  elengrab:
    image: neosy/elengrab:latest
    container_name: elengrab
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      TZ: "Europe/Moscow"   # your time zone
    volumes:
      - elengrab_db:/app_n/sqlite/data
      - elengrab_db_backups:/app_n/sqlite/backups
      - elengrab_downloads:/app_n/downloads
```

Then:

```
docker-compose up -d
```

After this, Elengrab will be accessible at http://localhost:8080

---

### Docker Compose for Docker Swarm

Create a file `docker-compose.yml` with the following content:

```
version: '3.5'

services:
  server:
    image: neosy/elengrab
    ports:
      - 8080:8080
    environment:
      TZ: "Europe/Moscow"
      ELENGRAB_DOWNLOAD_WORKERS: "3"
    volumes:
      - db:/app_n/sqlite/data
      - db_backups:/app_n/sqlite/backups
      - downloads:/app_n/downloads
    deploy:
      mode: replicated
      replicas: 1
      resources:
        limits:
          cpus: "4.0"
          memory: "4G"

volumes:
  db:
  db_backups:
  downloads:
```

---

### Deploying in Docker Swarm

1. Initialize Docker Swarm (if not already initialized):

```
docker swarm init
```

2. Deploy the stack using your `docker-compose.yml`:

```
docker stack deploy -c docker-compose.yml elengrab
```

3. Check the running services:

```
docker service ls
```

4. Access Elengrab at `http://<your-host-ip>:8080` and start managing your personal video and audio content.
