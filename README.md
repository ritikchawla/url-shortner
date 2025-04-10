# URL Shortener Service

A high-performance URL shortening service built with Go, featuring PostgreSQL for persistent storage and Redis for caching.

## Features

- Create short URLs with custom expiration
- Redirect short URLs to original URLs
- View URL statistics (visits, creation date, expiration)
- Redis caching for improved performance
- Docker containerization for easy deployment

## Tech Stack

- Go (Gin web framework)
- PostgreSQL (persistent storage)
- Redis (caching)
- Docker & Docker Compose

## API Endpoints

- `POST /api/shorten` - Create a short URL
- `GET /:shortCode` - Redirect to original URL
- `GET /api/stats/:shortCode` - Get URL statistics

## Setup

1. Clone the repository:
```bash
git clone https://github.com/ritikchawla/url-shortner.git
cd url-shortner
```

2. Start the services using Docker Compose:
```bash
docker-compose up --build
```

The service will be available at `http://localhost:8080`

## API Usage

### Create Short URL

```bash
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d '{"long_url": "https://example.com", "expires_at": "2024-12-31T23:59:59Z"}'
```

Response:
```json
{
  "id": "1",
  "long_url": "https://example.com",
  "short_code": "abc123",
  "visits": 0,
  "created_at": "2024-04-09T13:20:00Z",
  "expires_at": "2024-12-31T23:59:59Z"
}
```

### Get URL Statistics

```bash
curl http://localhost:8080/api/stats/abc123
```

### Redirect to Original URL

Simply visit: `http://localhost:8080/abc123`

## Development

To run tests:
```bash
go test ./tests -v
```

## Architecture

- Uses PostgreSQL for persistent storage of URLs and their metadata
- Redis caching for frequently accessed URLs
- Asynchronous visit counting
- Base64-encoded short codes for URL-safe identifiers
- Containerized services for easy deployment and scaling

## Best Practices

- Input validation for URLs
- Proper error handling and logging
- Caching for improved performance
- Database indexing for quick lookups
- Asynchronous statistics updates
- Containerized deployment
- Comprehensive test coverage