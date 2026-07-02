<p align="center">
  <img width="192" height="192" alt="android-chrome-192x192_round" src="https://github.com/user-attachments/assets/f2973dcc-90d3-4046-b9e0-fd34b6697fa5" />
</p>
<h3 align="center">Self-hosted web interface for downloading media from multiple platforms (YouTube, Instagram, TikTok, Twitch, etc.).</h3>

---

<p align="center">
  <img width="1248" height="517" alt="Screenshot_126" src="https://github.com/user-attachments/assets/e6f9cfa6-2a26-4330-8897-ca32c5e97b88" />
</p>

# Elengrab

**Fast self-hosted downloader for video and audio with flexible format and quality options. Integrates with media processing utilities (such as [yt-dlp](https://github.com/yt-dlp/yt-dlp)) as an optional backend. Supports downloading from multiple platforms, including YouTube, Facebook, Instagram, Twitter/X, Twitch, Pinterest, Reddit, VK Video, Rutube, and more.**

Elengrab provides a simple and **very fast** web interface to quickly download with videos and audio tracks, allowing selection of formats and quality settings. The project is fully written in **Go**, and the frontend is powered by **HTMX**, ensuring high responsiveness and minimal overhead. It runs in a lightweight Docker container, is easy to deploy, and serves as a self-hosted frontend for managing your personal media server.

---

## Features

- Support for downloading video and audio from multiple websites, powered by yt-dlp extractors.
- Support for YouTube content in various formats.
- Configurable format and quality settings.
- Cookie-based authentication support for YouTube (optional, requires Deno).
- Task queue for media processing.
- Concurrent processing of multiple tasks (3 by default).
- Instant addition/removal from queue.
- User authentication and access control.
- Built-in media search.
- In-browser video and audio playback.
- Option to customize how download history is displayed: globally or per user.
- Animated status indicators for download and processing tasks.
- Channel and website icons displayed in the interface.
- Dark theme.
- Cross-platform support: Windows, macOS/Linux.

---

## YouTube Content Context

Elengrab is designed to work with media hosted on YouTube.
It operates as a self-hosted interface for organizing and processing personal video and audio content.
Users are responsible for ensuring compliance with applicable laws and platform terms.

---

## Requirements

### Minimum

* **CPU:** 1–2 cores
* **Memory:** ~1 GB RAM
* **Concurrent downloads:** 1 worker
* **Dependencies:** `yt-dlp` and `ffmpeg` must be installed

This configuration is suitable for low-resource servers. To limit concurrency, set:

```
ELENGRAB_DOWNLOAD_WORKERS=1
```

### Recommended

* **CPU:** 4-6 cores
* **Memory:** 4 GB RAM
* **Concurrent downloads:** 3 workers (default)
* **Dependencies:** `yt-dlp` and `ffmpeg` must be installed

By default, the Docker container is configured to use this setup.

### Notes

The main resource consumers are `yt-dlp` and `ffmpeg`, especially during video downloading, merging, and transcoding.

The number of concurrent workers can be adjusted using the following environment variable:

```
ELENGRAB_DOWNLOAD_WORKERS=3
```

---

## Quick Start
### Run docker with minimum settings

```
docker run -d \
  --name elengrab \
  -v elengrab_downloads:/app_n/downloads \
  -p 8080:8080 \
  neosy/elengrab:latest
```

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `LOG_LEVEL` | `warn` | Logging level. Options: `debug`, `info`, `warn`, `error`. |
| `SQLITE_DATA_DIR` | `sqlite/data` | Directory where SQLite database files are stored. |
| `SQLITE_BACKUPS_DIR` | `sqlite/backups` | Directory where SQLite backup files are stored. |
| `ELENGRAB_ROOT_DIR` | *(not set)* | Application base directory. If empty, defaults to: `~/.elengrab` on Linux/macOS, `%LOCALAPPDATA%\Elengrab` on Windows. |
| `ELENGRAB_DOWNLOADER_BIN_DIR` | `/usr/local/bin` | Directory containing yt-dlp binary. |
| `ELENGRAB_ASSETS_DIR` | `assets` | Directory containing application assets. |
| `ELENGRAB_DOWNLOADS_DIR` | `downloads` | Directory where downloaded files are stored inside the container. Must be mapped to a host volume. |
| `ELENGRAB_COOKIES_DIR` | `cookies` | Directory where cookie files are stored. Files must be in Netscape cookies.txt format (compatible with yt-dlp). |
| `ELENGRAB_DOWNLOAD_WORKERS` | 3 | Number of concurrent workers used for processing YouTube video and audio tasks in parallel. |
| `ELENGRAB_MODE` | `public` | Controls access mode. Possible values: `public` (anonymous, full access, shared download history), `public_readonly` (anonymous, read-only access to public media only), `guest` (anonymous, session-isolated download history), `authenticated` (login required, permission-based access). |
| `ELENGRAB_ALLOW_COOKIES` | `false` | Enables use of cookies for multimedia sources. Requires **Deno**. The `youtube.txt` file should be located in the directory specified by `ELENGRAB_COOKIES_DIR`. |
| `ELENGRAB_MAINTENANCE_ENABLE_MOVE_UNMATCHED_FILES` | `false`   | Enables the periodic operation that moves files not present in the database tables from the download folder to the `.lost` folder. Default is `false` (disabled). |

---

## Volumes

| Volume                | Description                                                               |
| --------------------- | ------------------------------------------------------------------------- |
| `db:/app_n/sqlite/data`         | Stores SQLite database files.                                   |
| `db_backups:/app_n/sqlite/backups` | Stores backups of SQLite databases.                          |
| `downloads:/app_n/downloads`  | Stores downloaded files.                                          |
| `cookies:/app_n/cookies`    | Stores cookie files used by the application (see `ELENGRAB_COOKIES_DIR`). |
| `media:/app_n/media`    | Stores auxiliary media files such as thumbnails and other resources.    |

---

## Usage

### Run Docker with default settings

```
docker run -d \
  --name elengrab \
  -v elengrab_db:/app_n/sqlite/data \
  -v elengrab_db_backups:/app_n/sqlite/backups \
  -v elengrab_downloads:/app_n/downloads \
  -v elengrab_cookies:/app_n/cookies \
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
      TZ: "Europe/Moscow"
      ELENGRAB_DOWNLOAD_WORKERS: "3"
    volumes:
      - elengrab_db:/app_n/sqlite/data
      - elengrab_db_backups:/app_n/sqlite/backups
      - elengrab_downloads:/app_n/downloads
      - elengrab_cookies:/app_n/cookies
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
      - cookies:/app_n/cookies
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
  cookies:
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

---

# License

Copyright (C) 2025–2026 Yury Savonin (Neosy)

This project is licensed under the GNU Affero General Public License v3.0 or later (AGPL-3.0-or-later).

See the LICENSE file for the full license text.
