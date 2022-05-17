# Syncovery prometheus exporter
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/andrzejressel/syncovery-exporter)

Prometheus exporter for Syncovery backup software: https://www.syncovery.com/

## Quick start

```yaml
# docker-compose.yml
version: "3.9"
services:
  syncovery-exporter:
    image: ghcr.io/andrzejressel/syncovery-exporter:0.0.2
    command:
      - --url=SYNCOVERY_URL
```

## Setup Syncovery

Exporter using endpoint `/profile.json` which is available after adding `SkipProfileListAuth=1` to Syncovery config file. Because of latest changes related to encryption it has to be added a following way.

1. Create local file `syncovery.ini` with content `SkipProfileListAuth=1`
2. Login to Syncovery instance. Go to Program Settings
3. Select `Import Config Lines (INI Style)...` and choose created file.

You can limit IPs that endpoint is available to by using `SkipProfileListAuthForIP` option.
