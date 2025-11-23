# Deployment Guide

This guide covers different deployment scenarios for the Real-Time Forum application.

## Table of Contents

- [Docker Deployment](#docker-deployment)
- [Manual Deployment](#manual-deployment)
- [Cloud Deployment](#cloud-deployment)
- [Environment Configuration](#environment-configuration)
- [Database Management](#database-management)
- [Monitoring and Logging](#monitoring-and-logging)
- [Security Considerations](#security-considerations)

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
- Systemd (for service management)

### Build from Source

1. **Clone and build**

```bash
git clone <repository-url>
cd real-time-forum/backend
make build
```

2. **Configure environment**

```bash
cp .env.example /etc/forum/.env
# Edit /etc/forum/.env
```

3. **Create systemd service**

Create `/etc/systemd/system/forum.service`:

```ini
[Unit]
Description=Real-Time Forum Backend
After=network.target

[Service]
Type=simple
User=forum
Group=forum
WorkingDirectory=/opt/forum
EnvironmentFile=/etc/forum/.env
ExecStart=/opt/forum/bin/forum-backend
Restart=on-failure
RestartSec=5s

# Security hardening
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/forum/data

[Install]
WantedBy=multi-user.target
```

4. **Deploy the application**

```bash
# Create user
sudo useradd -r -s /bin/false forum

# Copy files
sudo mkdir -p /opt/forum/{bin,data}
sudo cp bin/forum-backend /opt/forum/bin/
sudo cp -r ../frontend /opt/forum/
sudo chown -R forum:forum /opt/forum

# Start service
sudo systemctl daemon-reload
sudo systemctl enable forum
sudo systemctl start forum
```

5. **Check status**

```bash
sudo systemctl status forum
sudo journalctl -u forum -f
```

### Nginx Reverse Proxy

Create `/etc/nginx/sites-available/forum`:

```nginx
upstream forum_backend {
    server localhost:8000;
}

server {
    listen 80;
    server_name forum.example.com;

    # Redirect to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name forum.example.com;

    # SSL Configuration
    ssl_certificate /etc/letsencrypt/live/forum.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/forum.example.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    # Logging
    access_log /var/log/nginx/forum_access.log;
    error_log /var/log/nginx/forum_error.log;

    # WebSocket support
    location /ws {
        proxy_pass http://forum_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_read_timeout 86400;
    }

    # API requests
    location /api {
        proxy_pass http://forum_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Static files (frontend)
    location / {
        root /opt/forum/frontend;
        try_files $uri $uri/ /index.html;
        
        # Cache static assets
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
    }
}
```

Enable the site:

```bash
sudo ln -s /etc/nginx/sites-available/forum /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## Cloud Deployment

### AWS EC2

1. **Launch EC2 instance**
   - AMI: Ubuntu 22.04 LTS
   - Instance type: t3.small or larger
   - Security group: Allow ports 22, 80, 443

2. **Connect to instance**

```bash
ssh -i your-key.pem ubuntu@your-instance-ip
```

3. **Install Docker**

```bash
sudo apt update
sudo apt install -y docker.io docker-compose
sudo usermod -aG docker ubuntu
```

4. **Deploy application**

```bash
git clone <repository-url>
cd real-time-forum
docker-compose up -d
```

5. **Set up SSL** (using Let's Encrypt)

```bash
sudo apt install -y certbot python3-certbot-nginx
sudo certbot --nginx -d forum.example.com
```

### Google Cloud Run

1. **Build and push image**

```bash
cd backend
gcloud builds submit --tag gcr.io/PROJECT_ID/forum-backend
```

2. **Deploy to Cloud Run**

```bash
gcloud run deploy forum-backend \
  --image gcr.io/PROJECT_ID/forum-backend \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars ENVIRONMENT=production
```

### Heroku

1. **Create Heroku app**

```bash
heroku create forum-app
```

2. **Deploy**

```bash
git push heroku main
```

3. **Configure environment**

```bash
heroku config:set ENVIRONMENT=production
heroku config:set DATABASE_PATH=/app/data/forum.db
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

## Database Management

### Backup

```bash
# Docker
docker exec forum-backend sqlite3 /app/data/forum.db ".backup /app/data/backup.db"
docker cp forum-backend:/app/data/backup.db ./backup-$(date +%Y%m%d).db

# Manual
sqlite3 /opt/forum/data/forum.db ".backup backup-$(date +%Y%m%d).db"
```

### Restore

```bash
# Docker
docker cp backup.db forum-backend:/app/data/
docker exec forum-backend sqlite3 /app/data/forum.db ".restore /app/data/backup.db"

# Manual
sqlite3 /opt/forum/data/forum.db ".restore backup.db"
```

### Automated Backups

Add to crontab:

```bash
# Daily backup at 2 AM
0 2 * * * /usr/local/bin/forum-backup.sh
```

Create `/usr/local/bin/forum-backup.sh`:

```bash
#!/bin/bash
BACKUP_DIR="/var/backups/forum"
DATE=$(date +%Y%m%d)
mkdir -p $BACKUP_DIR

# Backup database
docker exec forum-backend sqlite3 /app/data/forum.db ".backup /app/data/backup.db"
docker cp forum-backend:/app/data/backup.db $BACKUP_DIR/forum-$DATE.db

# Keep only last 30 days
find $BACKUP_DIR -name "forum-*.db" -mtime +30 -delete
```

## Monitoring and Logging

### Health Checks

```bash
# Basic health check
curl http://localhost:8000/health

# Detailed monitoring with cron
*/5 * * * * curl -f http://localhost:8000/health || systemctl restart forum
```

### Log Management

```bash
# View Docker logs
docker-compose logs -f

# View systemd logs
sudo journalctl -u forum -f

# Rotate logs
sudo logrotate /etc/logrotate.d/forum
```

### Metrics

Consider integrating:
- Prometheus for metrics
- Grafana for visualization
- ELK stack for log aggregation

## Security Considerations

### Checklist

- [ ] Use HTTPS in production
- [ ] Set strong CORS origins (not *)
- [ ] Enable rate limiting
- [ ] Regular security updates
- [ ] Backup database regularly
- [ ] Use environment variables for secrets
- [ ] Implement monitoring and alerting
- [ ] Use firewall to restrict ports
- [ ] Keep Go and dependencies updated
- [ ] Review and audit code regularly

### Firewall Configuration

```bash
# Allow SSH, HTTP, HTTPS
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

### SSL/TLS Configuration

Use Let's Encrypt for free SSL certificates:

```bash
sudo certbot --nginx -d forum.example.com
sudo certbot renew --dry-run
```

## Troubleshooting

### Service won't start

```bash
# Check logs
sudo journalctl -u forum -n 50
docker-compose logs

# Check configuration
cat /etc/forum/.env

# Test manually
cd /opt/forum
./bin/forum-backend
```

### Database issues

```bash
# Check database file
ls -la /opt/forum/data/forum.db

# Verify database integrity
sqlite3 /opt/forum/data/forum.db "PRAGMA integrity_check;"
```

### Performance issues

```bash
# Check resource usage
docker stats forum-backend

# Check connections
netstat -an | grep :8000

# Monitor logs
tail -f /var/log/nginx/forum_access.log
```

## Updating the Application

### Docker

```bash
# Pull latest changes
git pull origin main

# Rebuild and restart
docker-compose build
docker-compose up -d
```

### Manual

```bash
# Pull latest changes
cd /opt/forum-source
git pull origin main

# Build
cd backend
make build

# Stop service
sudo systemctl stop forum

# Update binary
sudo cp bin/forum-backend /opt/forum/bin/

# Start service
sudo systemctl start forum
```

## Rollback

### Docker

```bash
# Stop current version
docker-compose down

# Checkout previous version
git checkout <previous-commit>

# Start
docker-compose up -d
```

### Manual

```bash
# Stop service
sudo systemctl stop forum

# Restore previous binary
sudo cp /opt/forum/backups/forum-backend-previous /opt/forum/bin/forum-backend

# Restore database if needed
sqlite3 /opt/forum/data/forum.db ".restore /var/backups/forum/forum-YYYYMMDD.db"

# Start service
sudo systemctl start forum
```

## Support

For deployment issues:
1. Check logs
2. Verify configuration
3. Review this guide
4. Open an issue on GitHub

