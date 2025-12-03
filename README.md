# GoUrlShortener

A production-ready URL shortening service written in Go.
Features a clean architecture, Basic Auth, SQLite storage, structured logging, a REST API, and fully automated CI/CD deployment via GitHub Actions and SSH to a VPS.

This project was built as a backend engineering exercise focusing on real-world architecture, testing, and deployment workflows.

## Features

- Create short URLs
- Redirect to long URLs
- Delete existing shortcuts
- Basic Auth for write operations
- SQLite persistent storage
- Structured logging (slog)
- Clean Architecture (handlers → service → storage)
- Configuration via YAML and .env
- Unit tests (handlers, random generator)
- Automated deployment to VPS using systemd and GitHub Actions

## Project Structure

```
cmd/url-shortener
    main.go
internal/
    config/
    http-server/
        handlers/
        middleware/logger/
lib/
    api/
    logger/
    random/
    storage/
config/
deployment/
storage/
```

## Tech Stack

- Go (chi router, slog)
- SQLite
- REST API
- systemd
- GitHub Actions (CI and Deployment)
- SSH and rsync deployment to VPS
- Clean Architecture

## Authentication (Basic Auth)

All write operations require Basic Auth.

Default local credentials (from local.yaml):

```
username: admin
password: password
```

## API Endpoints

### Create short URL

```
POST /url/save
Auth: Basic
Content-Type: application/json
```

Request:

```json
{
  "url": "https://github.com/MaleevFedor/GoURLShortener",
  "alias": "source_code"
}
```

Response:

```json
{
  "status": "ok",
  "alias": "source_code"
}
```

### Redirect

```
GET /{alias}
```

Redirects with HTTP 302.

### Delete alias

```
DELETE /url/{alias}
Auth: Basic
```

## Alias Rules

- Cannot contain characters that break routing (`/`, `:`, spaces)
- Otherwise mostly unrestricted
- If omitted, system generates a random alias

## Running Locally

1. Clone the repository.
2. Create a `.env` file:

```
CONFIG_PATH=./config/local.yaml
```

3. Run:

```
go run ./cmd/url-shortener
```

### Example Requests

Create:

```
curl -u admin:password -X POST http://localhost:8082/url/save   -H "Content-Type: application/json"   -d '{"url":"https://example.com", "alias":"test"}'
```

Redirect:

```
curl -v http://localhost:8082/test
```

Delete:

```
curl -u admin:password -X DELETE http://localhost:8082/url/test
```

## Testing

Run unit tests:

```
go test ./...
```

## Deployment (CI/CD)

Includes a full GitHub Actions pipeline:

- Runs tests
- Builds the binary
- Deploys to a VPS via SSH and rsync
- Replaces systemd unit
- Restarts the service

Deployment workflow:

```
workflow_dispatch → build → rsync → replace systemd → restart service
```

## Production Configuration

Example prod.yaml:

```
env: "prod"
storage_path: "/root/apps/url-shortener/storage/storage.db"
http_server:
  host: "0.0.0.0"
  port: 8082
  timeout: 4s
  idle_timeout: 60s
  user: "<basic-auth-user>"
  password: "<basic-auth-pass>"
```

Systemd service location:

```
/etc/systemd/system/url-shortener.service
```

## Developer Notes

Skills practiced in this project:

- Clean architecture in Go
- Structured logging with slog
- Writing testable handler logic
- Dependency injection via interfaces
- CI/CD pipelines in GitHub Actions
- Deploying Go binaries to VPS via SSH
- Managing systemd services
- SQLite integration
