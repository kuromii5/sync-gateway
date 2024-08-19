# SYNC-GATEWAY

Gateway microservice for social network "SYNC". This project is currently under development.

## Installation

## Configuration setup

### .env example

Make sure you have a .env file in the root directory with the following content:

```env
PORT=8080
ENV=local
LOG_LEVEL=info

AUTH_SERVICE_ENDPOINT=localhost:44044

# not ready
USER_SERVICE_ENDPOINT=
FEED_SERVICE_ENDPOINT=
MESSAGE_SERVICE_ENDPOINT=
NOTIFICATION_SERVICE_ENDPOINT=
MUSIC_SERVICE_ENDPOINT=
VIDEO_SERVICE_ENDPOINT=
GROUP_SERVICE_ENDPOINT=
```

## Running the app

Run the next command to run service:

```bash
go run cmd/server/main.go
```
