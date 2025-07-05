# Accountability Backend - Docker Setup

This Go Fiber backend application has been dockerized and includes PostgreSQL database integration.

## Prerequisites

- Docker
- Docker Compose

## Quick Start

1. **Clone the repository** (if not already done)

2. **Create environment file**
   ```bash
   cp .env.example .env
   ```
   
   Edit the `.env` file with your preferred configuration:
   ```env
   DB_HOST=db
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=password123
   DB_NAME=accountability_db
   JWT_SECRET=your-jwt-secret-key-here
   PORT=5000
   ```

3. **Build and run with Docker Compose**
   ```bash
   docker-compose up --build
   ```

   Or run in detached mode:
   ```bash
   docker-compose up -d --build
   ```

4. **Access the application**
   - Backend API: http://localhost:5000
   - PostgreSQL: localhost:5432

## Available Commands

### Development
```bash
# Start services
docker-compose up

# Start services in background
docker-compose up -d

# Rebuild and start
docker-compose up --build

# Stop services
docker-compose down

# View logs
docker-compose logs backend
docker-compose logs db

# Access backend container shell
docker-compose exec backend sh

# Access PostgreSQL container
docker-compose exec db psql -U postgres -d accountability_db
```

### Production Build
```bash
# Build only the backend image
docker build -t accountability-backend .

# Run the backend container with external database
docker run -p 5000:5000 \
  -e DB_HOST=your-db-host \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=your-password \
  -e DB_NAME=accountability_db \
  -e JWT_SECRET=your-jwt-secret \
  accountability-backend
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | Database host | `db` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database username | `postgres` |
| `DB_PASSWORD` | Database password | `password123` |
| `DB_NAME` | Database name | `accountability_db` |
| `JWT_SECRET` | JWT signing secret | `your-jwt-secret-key-here` |
| `PORT` | Server port | `5000` |

### Docker Compose Services

- **backend**: Go Fiber application
- **db**: PostgreSQL 15 database with persistent volume

## Database

The PostgreSQL database is automatically configured with:
- Database: `accountability_db`
- User: `postgres`
- Password: `password123`
- Port: `5432`

Data is persisted in a Docker volume named `postgres_data`.

## Health Checks

The PostgreSQL service includes health checks to ensure the database is ready before starting the backend service.

## CORS Configuration

The backend is configured to allow requests from `http://localhost:3000` by default. Update the CORS configuration in `cmd/main.go` if needed.

## Troubleshooting

### Database Connection Issues
```bash
# Check if database is running
docker-compose ps

# Check database logs
docker-compose logs db

# Test database connection
docker-compose exec db pg_isready -U postgres
```

### Backend Issues
```bash
# Check backend logs
docker-compose logs backend

# Restart backend service
docker-compose restart backend
```

### Clean Start
```bash
# Stop and remove containers, networks, and volumes
docker-compose down -v

# Rebuild everything
docker-compose up --build
```
