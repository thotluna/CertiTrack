# CertiTrack

[![Go](https://img.shields.io/badge/Go-1.23-blue.svg)](https://golang.org/)
[![React](https://img.shields.io/badge/React-18.x-61DAFB.svg)](https://reactjs.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

CertiTrack is a comprehensive certification management system designed to help organizations track and manage employee and equipment certifications efficiently. The system provides automated alerts for upcoming expirations, document management, and detailed reporting capabilities.

## Features

- **Certification Management**: Track all types of certifications with expiration dates
- **Person & Equipment Tracking**: Associate certifications with both employees and equipment
- **Automated Alerts**: Get notified before certifications expire
- **Document Storage**: Attach and manage certification documents
- **Role-based Access Control**: Secure access based on user roles
- **Reporting & Analytics**: Generate reports on certification statuses

## Tech Stack

### Backend
- **Language**: Go 1.23
- **Framework**: Gin
- **ORM**: GORM
- **Database**: PostgreSQL
- **Containerization**: Docker

### Frontend
- **Framework**: React 18 with Next.js
- **UI Library**: Chakra UI
- **State Management**: React Query
- **Form Handling**: React Hook Form

## Prerequisites

- Go 1.23+
- Node.js 18+
- Docker & Docker Compose
- PostgreSQL 14+

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/certitrack.git
cd certitrack
```

### 2. Set Up Environment Variables

Copy the example environment files and update them with your configuration:

```bash
# Backend
cp backend/.env.example backend/.env

# Frontend
cp frontend/.env.example frontend/.env.local
```

### 3. Start the Development Environment

Using Docker Compose (recommended):

```bash
docker-compose up -d
```

Or run services manually:

#### Backend

```bash
cd backend
go mod download
go run cmd/certitrack-backend/main.go
```

#### Frontend

```bash
cd frontend
npm install
npm run dev
```

### 4. Access the Application

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Database**: PostgreSQL on port 5432

## Project Structure

```
certitrack/
├── backend/               # Backend services
│   ├── cmd/              # Application entry points
│   ├── feature/          # Feature modules
│   ├── shared/           # Shared packages
│   └── migrations/       # Database migrations
├── frontend/             # Frontend application
│   ├── public/           # Static files
│   └── src/              # Source code
├── docs/                 # Documentation
└── docker-compose.yml    # Docker Compose configuration
```

## API Documentation

API documentation is available at `/api/v1/docs` when running the backend in development mode.

## Development

### Code Style

- **Go**: Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- **JavaScript/TypeScript**: Follow [Airbnb Style Guide](https://github.com/airbnb/javascript)
- **Git**: Follow [Conventional Commits](https://www.conventionalcommits.org/)

### Testing

Run tests for the backend:

```bash
cd backend
go test ./...
```

Run tests for the frontend:

```bash
cd frontend
npm test
```

## Deployment

### Production Build

```bash
# Build and start production containers
docker-compose -f docker-compose.prod.yml up -d --build
```

### Environment Variables

Required environment variables are documented in `.env.example` files in each service directory.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support, please open an issue in the GitHub repository.
