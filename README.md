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

```
docker run -d \
  --name elengrab \
  -v elengrab_downloads:/app_n/downloads \
  -p 8080:8080 \
  neosy/elengrab:latest
```

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
    volumes:
      - downloads:/app_n/downloads
    deploy:
      mode: replicated
      replicas: 1
      resources:
        limits:
          memory: "512M"
          cpus: "2.0"

volumes:
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

4. Access Elengrab at `http://<your-host-ip>:8080` and enjoy downloading videos and audio.
