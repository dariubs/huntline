# HuntLine

**HuntLine** is a multi-platform launch tracker that tracks top products from various launch platforms like ProductHunt, altern.ai, and tinylaunch daily. The first version focuses on ProductHunt, with a protocol designed to easily add other launch platforms.

## What is it?

HuntLine provides a clean, modern interface to discover and track the best products launched across multiple platforms. It aggregates daily top products, allowing you to:

- Browse top products by date with intuitive day/month navigation
- View historical archives organized by month
- Discover best products of the week and month
- Track products across multiple launch platforms
- Enjoy a responsive, dark-mode enabled interface

## Features

- **Multi-platform support**: Built with an extensible platform protocol
- **Daily tracking**: Automatically fetches top products daily
- **Beautiful UI**: Modern, minimal design with dark mode support
- **Historical data**: Browse archives and best-of collections
- **Date navigation**: Easy browsing through different time periods
- **Responsive design**: Works seamlessly on desktop and mobile

## Prerequisites

- Go 1.23 or higher
- PostgreSQL database
- ProductHunt API key (or API keys for other platforms you want to use)

## Installation

1. **Clone the repository:**

```bash
git clone https://github.com/dariubs/huntline.git
cd huntline
```

2. **Install dependencies:**

```bash
make install
# or
go mod download
go mod tidy
```

3. **Set up environment variables:**

Create a `.env` file in the root directory:

```env
# Database configuration
PG_USER=your_db_user
PG_PASS=your_db_password
PG_NAME=huntline
PG_HOST=localhost
PG_PORT=5432

# ProductHunt API
PH_API_KEY=your_producthunt_api_key

# Server configuration
HL_PORT=8080
HL_NAME=HuntLine
HL_URL=http://localhost:8080
HL_LOGO=
HL_FAVICON=
HL_CDN=
HL_X=
HL_GITHUB=
```

4. **Run database migrations:**

```bash
make migrate
# or
go run app/main/migrate/migrate.go
```

## How to Run

### Running the Web Server

Start the web server to view the HuntLine interface:

```bash
make run-server
# or
go run app/main/huntline/main.go
```

The server will start on `http://localhost:8080` (or the port specified in `HL_PORT`).

### Running the Receiver

The receiver fetches product data from launch platforms. Run it to update yesterday's data:

```bash
make run-receiver
# or
go run app/main/receiver/main.go
```

#### Receiver Options

- **Fetch data for a specific date:**
  ```bash
  make receiver-date DATE=2025-01-15
  # or
  go run app/main/receiver/main.go -date 2025-01-15
  ```

- **Run on a daily schedule:**
  ```bash
  make receiver-repeat
  # or
  go run app/main/receiver/main.go -repeat=true
  ```

- **Backfill historical data:**
  ```bash
  make receiver-historical
  # or
  go run app/main/receiver/main.go -historical=true
  ```

- **Update last month's data:**
  ```bash
  make receiver-last-month
  # or
  go run app/main/receiver/main.go -last-month=true
  ```

### Development Commands

```bash
# Build all binaries
make build

# Build specific components
make build-server
make build-receiver
make build-migrate

# Clean build artifacts
make clean

# Run tests (if available)
make test
```

## Project Structure

```
huntline/
├── app/
│   ├── db/              # Database connection
│   ├── handler/         # HTTP handlers
│   │   └── huntline/    # HuntLine-specific handlers
│   ├── main/            # Application entry points
│   │   ├── huntline/    # Web server
│   │   ├── receiver/    # Data fetcher
│   │   └── migrate/     # Database migrations
│   ├── model/           # Database models
│   ├── platform/        # Platform protocol implementations
│   └── types/           # Shared types and utilities
├── view/                # HTML templates
├── assets/              # Static assets
├── Makefile            # Build automation
└── go.mod              # Go dependencies
```

## Adding a New Launch Platform

HuntLine uses a protocol-based design to easily add new platforms:

1. Implement the `LaunchPlatform` interface in `app/platform/platform.go`
2. Create your platform package (e.g., `app/platform/yourplatform/`)
3. Register the platform in the receiver

See `app/platform/producthunt/` for a reference implementation.

## Timezone

All date calculations use **San Francisco timezone (Pacific Time)** to ensure consistency across the application. This matches ProductHunt's timezone.

## How to Contribute

We welcome contributions! Here's how you can help:

### Reporting Bugs

1. Check existing issues to see if the bug has been reported
2. Create a new issue with:
   - Clear description of the bug
   - Steps to reproduce
   - Expected vs actual behavior
   - Environment details (OS, Go version, etc.)

### Suggesting Features

1. Check existing issues/PRs for similar suggestions
2. Open an issue describing:
   - The feature and its use case
   - How it would benefit users
   - Any implementation ideas you have

### Submitting Code

1. **Fork the repository**

2. **Create a feature branch:**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make your changes:**
   - Follow Go conventions and best practices
   - Add comments for complex logic
   - Ensure code is clean and readable
   - Test your changes

4. **Commit your changes:**
   ```bash
   git commit -m "Add: descriptive commit message"
   ```
   
   Use clear commit messages:
   - `Add:` for new features
   - `Fix:` for bug fixes
   - `Update:` for improvements
   - `Refactor:` for code restructuring

5. **Push to your fork:**
   ```bash
   git push origin feature/your-feature-name
   ```

6. **Open a Pull Request:**
   - Provide a clear description of your changes
   - Reference any related issues
   - Wait for code review and feedback

### Code Style Guidelines

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Add comments for exported functions and complex logic
- Keep functions focused and small
- Handle errors explicitly (don't ignore them)
- Use consistent error handling patterns

### Testing

While tests aren't required yet, they're encouraged. If you add features, consider adding tests to verify they work correctly.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

- **Issues**: [GitHub Issues](https://github.com/dariubs/huntline/issues)
- **Discussions**: [GitHub Discussions](https://github.com/dariubs/huntline/discussions)

## Acknowledgments

- Built with [Gin](https://gin-gonic.com/) web framework
- Uses [GORM](https://gorm.io/) for database operations
- Styled with [Tailwind CSS](https://tailwindcss.com/)
- Powered by [Alpine.js](https://alpinejs.dev/)

