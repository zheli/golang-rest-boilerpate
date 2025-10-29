# Golang REST Boilerplate

A production-ready REST API boilerplate built with Golang, PostgreSQL, JWT authentication, and Google OAuth support. It includes Docker configurations for local development and production, automated testing scaffolding, and CI/CD via GitHub Actions.

## Features

- Gin-based HTTP server with modular architecture
- PostgreSQL database integration via GORM
- JWT authentication with refreshable configuration
- Google OAuth 2.0 sign-in flow
- User registration, login, and CRUD management endpoints
- Health check endpoint (`/health`)
- Dockerfile and Compose setup for dev/prod
- Makefile for common tasks (run, test, build, docker compose)
- GitHub Actions CI pipeline running formatting and tests
- Example unit test covering authentication service

## Getting Started

### Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- [Docker & Docker Compose](https://docs.docker.com/get-docker/)
- [Make](https://www.gnu.org/software/make/)

### Environment Variables

Copy the provided template and update values as needed:

```bash
cp .env.example .env
```

Key variables:

- `APP_PORT`: HTTP server port (default `8080`).
- `DATABASE_URL`: PostgreSQL connection string.
- `JWT_SECRET`: Secret key for signing JWTs.
- `JWT_ISSUER`: Issuer claim embedded in JWTs.
- `TOKEN_EXPIRE_MINUTES`: Access token lifetime.
- `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `GOOGLE_REDIRECT_URL`: Google OAuth configuration (optional).

### Local Development

Use the development profile which mounts source code and runs `go run` inside the container:

```bash
docker compose --profile dev up --build
```

The API becomes available at `http://localhost:8080`. Update code locally and the container automatically picks it up.

Alternatively, run locally without Docker:

```bash
make run
```

Run tests:

```bash
make test
```

Format code:

```bash
make fmt
```

### Production Build

Build the optimized binary using the included multi-stage Dockerfile:

```bash
docker compose --profile prod up --build -d
```

This profile builds the Go binary and runs it in a lightweight Alpine container. Customize environment variables in `docker-compose.yml` or via an external `.env` file when deploying.

### Database

The Compose file provisions a PostgreSQL 15 container with default credentials (`postgres/postgres`). Persisted data is stored in the named Docker volume `postgres_data`.

You can connect to the database using any Postgres client:

```bash
psql postgres://postgres:postgres@localhost:5432/app
```

### API Overview

| Method | Endpoint | Description | Auth |
| ------ | -------- | ----------- | ---- |
| GET    | `/health` | Service health check | None |
| POST   | `/api/v1/auth/register` | Register a new user | None |
| POST   | `/api/v1/auth/login` | Email/password login | None |
| GET    | `/api/v1/auth/google/login` | Start Google OAuth flow | None |
| GET    | `/api/v1/auth/google/callback` | Google OAuth callback | None |
| GET    | `/api/v1/users` | List users | Bearer token |
| GET    | `/api/v1/users/:id` | Get a user by ID | Bearer token |
| PUT    | `/api/v1/users/:id` | Update user name | Bearer token |
| DELETE | `/api/v1/users/:id` | Delete user | Bearer token |

The JWT token should be sent in the `Authorization: Bearer <token>` header for protected routes.

### Google OAuth Setup

1. Create an OAuth 2.0 Client ID in the [Google Cloud Console](https://console.cloud.google.com/).
2. Set the authorized redirect URI to `http://localhost:8080/api/v1/auth/google/callback` (or your deployment URL).
3. Populate `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, and `GOOGLE_REDIRECT_URL` environment variables.
4. Call `GET /api/v1/auth/google/login` to receive the authorization URL and `state` token. Redirect the user there to complete the login.
5. Handle the callback and exchange the returned JWT token for API access.

### Running Tests in CI

The GitHub Actions workflow automatically runs `go fmt` (as a check) and `go test ./...` on every push or pull request targeting the `main` branch.

### Makefile Targets

- `make run`: start the API locally.
- `make test`: execute unit tests.
- `make build`: compile the binary into `bin/server`.
- `make fmt`: format the Go source code.
- `make docker-up`: run the dev Docker Compose profile.
- `make docker-down`: stop containers.

### Folder Structure

```
├── cmd/server           # Application entry point
├── internal             # Application code (config, db, HTTP handlers, services)
├── pkg                  # Shared helpers
├── .github/workflows    # CI pipeline definition
├── docker-compose.yml   # Docker services
├── Dockerfile           # Multi-stage build
├── Makefile             # Developer tasks
└── README.md            # Project documentation
```

### Notes

- The unit test uses an in-memory SQLite database to exercise authentication logic without needing a running PostgreSQL instance.
- Before deploying, ensure you set secure values for secrets and consider integrating a secrets manager.

## License

This project is provided as boilerplate under the MIT License. Customize as needed for your applications.
