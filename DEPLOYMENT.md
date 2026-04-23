# Deployment Guide

This guide covers different deployment scenarios for the Real-Time Forum application.

## Table of Contents

- [Docker Deployment](#docker-deployment)
- [Manual Deployment](#manual-deployment)
- [Environment Configuration](#environment-configuration)

## Docker Deployment

### Using Docker Compose (Recommended)

The easiest way to deploy the application is using Docker Compose.

1. **Clone the repository**

```bash
git clone <repository-url>
cd real-time-forum
```

2. **Configure environment**

```bash
cd backend
cp .env.example .env
# Edit .env with production values
```

3. **Start the application**

```bash
docker-compose up -d
```

4. **Verify deployment**

```bash
docker-compose ps
docker-compose logs -f
```

5. **Check health**

```bash
curl http://localhost:8000/health
```

### Using Docker Only

1. **Build the image**

```bash
cd backend
docker build -t forum-backend:latest .
```

2. **Run the container**

```bash
docker run -d \
  --name forum-backend \
  -p 8000:8000 \
  -e ENVIRONMENT=production \
  -e DATABASE_PATH=/app/data/forum.db \
  -e RATE_LIMIT_ENABLED=true \
  -v forum-data:/app/data \
  -v $(pwd)/../frontend:/app/frontend:ro \
  forum-backend:latest
```

3. **View logs**

```bash
docker logs -f forum-backend
```

## Manual Deployment

### Prerequisites

- Go 1.21 or higher
- SQLite3

### Build from Source

1. **Clone and build**

```bash
git clone <repository-url>
cd real-time-forum/backend
make build
```

2. **Configure environment**

```bash
cp .env.example .env
# Edit /etc/forum/.env
```

3. **Run the application**

```bash
make run
```

4. **Check health**

```bash
curl http://localhost:8000/health
```

## Environment Configuration

### Production Environment Variables

```bash
# Environment
ENVIRONMENT=production

# Server
PORT=8000
SERVER_READ_TIMEOUT=15s
SERVER_WRITE_TIMEOUT=15s
SERVER_IDLE_TIMEOUT=60s

# Database
DATABASE_PATH=/app/data/forum.db

# Session
SESSION_DURATION=24h

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_RPM=100

# CORS (set to your domain)
CORS_ALLOWED_ORIGINS=https://forum.example.com

# WebSocket
WS_READ_BUFFER_SIZE=1024
WS_WRITE_BUFFER_SIZE=1024
WS_PING_PERIOD=54s
WS_PONG_WAIT=60s
WS_WRITE_WAIT=10s
```