<p align="center">
  <img width="192" height="192" alt="android-chrome-192x192_round" src="https://github.com/user-attachments/assets/f2973dcc-90d3-4046-b9e0-fd34b6697fa5" />
</p>
<h3 align="center">Self-hosted cross-platform downloader and media viewer for video and audio from 1,000+ platforms (YouTube, Instagram, TikTok, Twitch, etc.).</h3>

---

<p align="center">
  <img width="1867" height="915" alt="elengrab_interface" src="https://github.com/user-attachments/assets/4eb2c205-0115-42bc-b40a-3e29e5bdb57f" />
</p>

# Elengrab

**Fast cross-platform application for downloading and watching video and audio with flexible format and quality options. Integrates with media processing utilities such as [yt-dlp](https://github.com/yt-dlp/yt-dlp) and [ffmpeg](https://github.com/FFmpeg/FFmpeg). Supports downloading from more than 1,000 websites and platforms, including YouTube, Facebook, Instagram, Twitter/X, Twitch, Pinterest, Reddit, VK Video, Rutube, and more.**

Elengrab provides a simple and very fast web interface for downloading and watching media directly in the browser, with a YouTube-like viewing experience, thumbnails, video previews, and a dedicated interface optimized for mobile devices. Supports Grid View, short links for sharing, and flexible playback options. On mobile devices, videos can be shared directly from apps such as YouTube and Instagram to Elengrab, where they are automatically picked up and downloaded.

The project is fully written in **Go**, with a frontend built using plain **HTML** and **CSS**, **HTMX**, and vanilla **JavaScript** without external libraries. It can run as a single binary on Windows, Linux, and macOS, or in a lightweight Docker container for easy deployment. Elengrab also includes full **PWA** support, allowing it to be installed and used as a standalone application. Different operating modes allow you to enable authentication and use Elengrab as a full-featured web service accessible over the Internet. This makes it well suited for creating a fast-access home media service as well as a personal media service accessible from anywhere.

### Screenshots of the interface
<details>
<summary><strong>Desktop — Light and dark themes</strong></summary>
  &nbsp;&nbsp;
  <p align="center">
    <img width="1842" height="909" alt="Screenshot_160" src="https://github.com/user-attachments/assets/2153a3c9-0423-483b-bb0f-6e510df3505f" />
    &nbsp;&nbsp;
    <img width="1867" height="915" alt="Screenshot_158" src="https://github.com/user-attachments/assets/0ae7c011-cdba-42ae-b344-65ffaaf59dac" />
  </p>
</details>
<details>
<summary><strong>Desktop — Grid view</strong></summary>
  &nbsp;&nbsp;
  <p align="center">
    <img width="1323" height="625" alt="image" src="https://github.com/user-attachments/assets/5b5a5db6-ab0d-4eb3-adfe-ea7167e1238d" />
    &nbsp;&nbsp;
    <img width="1327" height="619" alt="image" src="https://github.com/user-attachments/assets/4a4f52dd-ac82-4fd1-83cd-11bc7a1f9a00" />
  </p>
</details>
<details>
<summary><strong>Desktop — Menu and Editing</strong></summary>
  &nbsp;&nbsp;
  <p align="center">
    <img width="1235" height="664" alt="image" src="https://github.com/user-attachments/assets/90d0b894-91c3-47e6-9b32-ea08e8a491e8" />
    &nbsp;&nbsp;&nbsp;
    <img width="1230" height="661" alt="image" src="https://github.com/user-attachments/assets/caac8f79-58e6-48ba-acdc-95baa0867dd4" />
  </p>
</details>

<details>
<summary><strong>Mobile — Light and dark themes</strong></summary>
  &nbsp;&nbsp;
  <p align="center">
    <img width="360" height="740" alt="image" src="https://github.com/user-attachments/assets/b59b44ed-0037-4751-a53f-8b0a145c8fda" />
    &nbsp;&nbsp;&nbsp;
    <img width="360" height="740" alt="image" src="https://github.com/user-attachments/assets/ee80e7a2-af13-4e4f-8ea7-5d3b47256b62" />
  </p>
  &nbsp;&nbsp;&nbsp;
  <p align="center">
  <img width="360" height="740" alt="image" src="https://github.com/user-attachments/assets/eaf4df5d-cf7d-41aa-bfe4-e28aeeed4063" />
    &nbsp;&nbsp;&nbsp;
  <img width="360" height="740" alt="image" src="https://github.com/user-attachments/assets/a1fd3bc7-76d4-4305-b4ce-8af14e49b481" />
  </p>
</details>
<details>
<summary><strong>Mobile — Menu and Editing</strong></summary>
  &nbsp;&nbsp;
  <p align="center">
    <img width="360" height="740" alt="image" src="https://github.com/user-attachments/assets/2ff66a7b-3e7c-48e4-9b41-d2127d70f5c0" />
    &nbsp;&nbsp;&nbsp;
    <img width="360" height="740" alt="image" src="https://github.com/user-attachments/assets/d0f7c3fc-4264-4f72-b36f-da6062a20eab" />
  </p>
</details>

---

## Features

* Support for downloading video and audio from **1,000+ websites**, powered by [yt-dlp](https://github.com/yt-dlp/yt-dlp) extractors.
* Configurable format and quality settings.
* Cookie-based authentication support for YouTube and other supported websites.
* Task queue for media processing.
* Concurrent processing of multiple tasks.
* Instant addition and removal of tasks from the queue.
* Built-in media search.
* **In-browser video and audio playback.**
* **Video thumbnails and media type indicators.**
* **Video preview on desktop and mobile devices.**
* **Watch progress tracking with resume playback and view statistics.**
* **Grid View and list view for browsing media.**
* **Short links for sharing media.**
* **Mobile sharing support for importing videos directly from apps such as YouTube and Instagram.**
* **PWA support for installing Elengrab as a standalone application.**
* **Multiple access modes for anonymous, read-only, session-isolated, and authenticated access.**
* **Configurable media visibility and public access.**
* Animated status indicators for download and processing tasks.
* Real-time task and media updates via SSE.
* Channel and website icons displayed in the interface.
* Responsive interface optimized for desktop and mobile devices.
* Dark theme.
* Cross-platform support: Windows, macOS, and Linux.
* Docker support for easy deployment.

---

## Media Content

Elengrab is designed to work with video and audio content from a wide range of websites and platforms through [yt-dlp](https://github.com/yt-dlp/yt-dlp).

Users are responsible for ensuring that their use of Elengrab and downloaded content complies with applicable laws, copyright requirements, and the terms of service of the respective platforms.

---

## Requirements

### Minimum

* **CPU:** 1–2 cores
* **Memory:** ~512 MB RAM
* **Concurrent downloads:** 1 worker
* **Dependencies:** [yt-dlp](https://github.com/yt-dlp/yt-dlp) and [ffmpeg](https://github.com/FFmpeg/FFmpeg) must be installed
* **Optional:** [Deno](https://github.com/denoland/deno) is required when using cookie support

This configuration is suitable for low-resource servers. To limit concurrency, set:

```
ELENGRAB_DOWNLOAD_WORKERS=1
```

### Recommended

* **CPU:** 4 cores
* **Memory:** 4 GB RAM
* **Concurrent downloads:** 3 workers (default)
* **Dependencies:** [yt-dlp](https://github.com/yt-dlp/yt-dlp) and [ffmpeg](https://github.com/FFmpeg/FFmpeg) must be installed
* **Optional:** [Deno](https://github.com/denoland/deno) is required when using cookie support

By default, the Docker container is configured to use this setup.

### Notes

The main resource consumers are `yt-dlp` and `ffmpeg`, especially during video downloading, merging, and transcoding.

Resource usage increases with the number of concurrent workers. The number of concurrent workers can be adjusted using the following environment variable:

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
| `ELENGRAB_LOG_LEVEL` | `warn` | Logging level. Options: `debug`, `info`, `warn`, `error`. |
| `ELENGRAB_SQLITE_DATA_DIR` | `sqlite/data` | Directory where SQLite database files are stored. |
| `ELENGRAB_SQLITE_BACKUPS_DIR` | `sqlite/backups` | Directory where SQLite backup files are stored. |
| `ELENGRAB_ROOT_DIR` | *(not set)* | Application base directory. If empty, defaults to: `~/.elengrab` on Linux/macOS, `%LOCALAPPDATA%\Elengrab` on Windows. |
| `ELENGRAB_DOWNLOADER_BIN_DIR` | `/usr/local/bin` | Directory containing yt-dlp binary. |
| `ELENGRAB_ASSETS_DIR` | `assets` | Directory containing application assets. |
| `ELENGRAB_DOWNLOADS_DIR` | `downloads` | Directory where downloaded files are stored inside the container. Must be mapped to a host volume. |
| `ELENGRAB_COOKIES_DIR` | `cookies` | Directory where cookie files are stored. Files must be in Netscape cookies.txt format (compatible with yt-dlp). |
| `ELENGRAB_DOWNLOAD_WORKERS` | 3 | Number of concurrent workers used for processing YouTube video and audio tasks in parallel. |
| `ELENGRAB_MODE` | `public` | Controls access mode. Possible values: `public` (anonymous, full access, shared download history), `public_readonly` (anonymous, read-only access to public media only), `guest` (anonymous, session-isolated download history), `authenticated` (login required, permission-based access). |
| `ELENGRAB_BASE_URL` | *(not set)* | Base URL of the Elengrab service, for example `https://example.com`. Used, for example, to generate short links. |
| `ELENGRAB_ALLOW_COOKIES` | `false` | Enables use of cookies for multimedia sources. Requires [Deno](https://github.com/denoland/deno). The `youtube.txt` file should be located in the directory specified by `ELENGRAB_COOKIES_DIR`. |
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

### Run as a standalone application

Elengrab is distributed as a single binary for Windows, Linux, and macOS. Download the latest release for your operating system and run the application directly without Docker.

Elengrab requires [yt-dlp](https://github.com/yt-dlp/yt-dlp) and [ffmpeg](https://github.com/FFmpeg/FFmpeg), which are not included in the standard releases.

#### Windows

##### Standard release

Download the latest Windows release from the [GitHub Releases](https://github.com/neosy/elengrab/releases) page and run `elengrab.exe`.

After starting Elengrab, open `http://localhost:8080` in your browser.

##### Portable Windows package

A portable Windows package containing the latest `elengrab.exe`, `yt-dlp.exe`, `ffmpeg.exe`, and `ffprobe.exe` is also available [here](https://nc.n-hub.ru/index.php/s/XyTm8HqginkwECT).

After starting Elengrab, open `http://localhost:2380` in your browser.

#### Linux

Download the latest Linux binary for your architecture from the [GitHub Releases](https://github.com/neosy/elengrab/releases) page.

Make the binary executable and run it:

```bash
chmod +x elengrab
./elengrab
```

Install [yt-dlp](https://github.com/yt-dlp/yt-dlp) and [ffmpeg](https://github.com/FFmpeg/FFmpeg) separately.

After starting Elengrab, open `http://localhost:8080` in your browser.

#### macOS

Download the latest macOS binary for your architecture from the [GitHub Releases](https://github.com/neosy/elengrab/releases) page.

Make the binary executable and run it:

```bash
chmod +x elengrab
./elengrab
```

Install [yt-dlp](https://github.com/yt-dlp/yt-dlp) and [ffmpeg](https://github.com/FFmpeg/FFmpeg) separately.

After starting Elengrab, open `http://localhost:8080` in your browser.

---

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

## Support

If you find Elengrab useful, you can support its development with a donation:

**[Support Elengrab](https://n-hub.ru/donate)**

---

# License

Copyright (C) 2025–2026 Yury Savonin (Neosy)

This project is licensed under the GNU Affero General Public License v3.0 or later (AGPL-3.0-or-later).

See the [`LICENSE`](LICENSE) file for the full license text.
