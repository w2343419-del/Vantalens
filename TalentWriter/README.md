# Vantalens Writer - Refactored Project

## Project Structure

```
TalentWriter/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── models/
│   │   └── models.go        # Data structures
│   ├── config/
│   │   └── config.go        # Configuration management
│   ├── auth/
│   │   └── auth.go          # JWT authentication
│   ├── comment/
│   │   └── comment.go       # Comment service
│   ├── email/
│   │   └── email.go         # Email notification service
│   ├── handlers/
│   │   └── handlers.go      # HTTP handlers
│   └── utils/               # Utility functions
├── web/
│   ├── templates/           # HTML templates
│   └── static/              # Static assets
├── go.mod
└── README.md
```

## Features

- JWT-based authentication
- Comment management with moderation
- Email notifications for pending comments
- Hugo static site integration
- CORS support for API endpoints
- Dual-backend architecture:
	- Control backend (site-wide control and operational checks)
	- Writer backend (article editing and content workflows)

## Build

```bash
go build -o talentwriter ./cmd/server
```

## Run

```bash
HUGO_PATH=/path/to/hugo ADMIN_TOKEN=your-token ./talentwriter
```

### Run Dual Backends (Recommended)

Control backend (default port `9090`):

```bash
go run ./cmd/control
```

Writer backend (default port `9091`):

```bash
go run ./cmd/writer
```

Optional environment variables:

- `CONTROL_PORT` for control backend port
- `WRITER_PORT` for writer backend port
- `ADMIN_TOKEN` or `ADMIN_PASSWORD` for admin login password

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| /api/login | POST | Admin login |
| /api/posts | GET | Get all posts |
| /api/comments | GET | Get comments for a post |
| /api/comments/add | POST | Add a new comment |
| /api/comments/approve | POST | Approve a comment (auth required) |
| /api/comments/delete | DELETE | Delete a comment (auth required) |
| /api/settings | GET | Get comment settings |
| /api/settings/save | POST | Save settings (auth required) |

### Control Backend Endpoints (`mode=control`)

- `/api/login`
- `/api/control/status`
- `/api/control/command`
- `/platform/control`

### Writer Backend Endpoints (`mode=writer`)

- `/api/login`
- `/api/posts`
- `/api/get_content`
- `/api/save_content`
- `/api/delete_post`
- `/api/create_post`
- `/api/comments`
- `/api/settings`
- `/platform/backend`
