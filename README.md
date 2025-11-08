# Elengrab

**Fast YouTube video/audio grabber with format and quality options.**

Elengrab allows you to quickly download YouTube videos or audio tracks with customizable format and quality options. Runs in a lightweight Docker container and is easy to deploy.

---

## Features

- Download YouTube videos in various formats.
- Download only audio tracks.
- Configure download quality.

---

## Requirements

- Docker 20.10+  
- Optional: Docker Compose for multi-container setups.

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `LOG_LEVEL` | `warn` | Logging level. Options: `debug`, `info`, `warn`, `error`. |
| `ELENGRAB_DOWNLOADS_DIR` | `/app_n/downloads` | Directory where downloaded files are stored inside the container. Must be mapped to a host volume. |

---

## Usage

### Run with default settings

```bash
docker run -d \
  --name elengrab \
  -v elengrab_downloads:/app_n/downloads \
  -p 8080:8080 \
  neosy/elengrab:latest
