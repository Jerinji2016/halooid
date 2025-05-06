# Halooid Platform Documentation

Welcome to the Halooid platform documentation. This documentation provides comprehensive information about the platform's architecture, components, and development guidelines.

## Table of Contents

### Architecture

- [Architecture Overview](architecture/overview.md): High-level overview of the platform architecture

### Components

- [Web UI Component Library](components/ui-components.md): Documentation for the shared web UI component library
- [Mobile Widget Library](components/mobile-widgets.md): Documentation for the shared mobile widget library

### Backend Services

- [Authentication Service](backend/authentication.md): JWT-based authentication service
- [RBAC Service](backend/rbac.md): Role-based access control service
- [API Gateway](backend/api-gateway.md): API Gateway for routing requests

### DevOps

- [CI/CD Pipeline](devops/ci-cd.md): Continuous Integration and Deployment pipeline
- [Kubernetes Deployment](devops/kubernetes.md): Kubernetes deployment configuration

## Getting Started

### Prerequisites

- Go 1.20 or later
- Node.js 18 or later
- Flutter 3.10.0 or later
- Docker and Docker Compose
- Kubernetes CLI (kubectl)

### Local Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/Jerinji2016/halooid.git
   cd halooid
   ```

2. Start the backend services:
   ```bash
   cd backend
   docker-compose up -d
   ```

3. Start the web frontend:
   ```bash
   cd web
   npm install
   npm run dev
   ```

4. Start the mobile app:
   ```bash
   cd mobile
   flutter pub get
   flutter run
   ```

## Development Guidelines

### Code Style

- **Go**: Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- **TypeScript/JavaScript**: Follow the ESLint configuration
- **Dart/Flutter**: Follow the Flutter analysis rules

### Git Workflow

1. Create a feature branch from `main`
2. Make changes and commit with descriptive messages
3. Push the branch and create a pull request
4. Address review comments
5. Merge the pull request

### Testing

- Write unit tests for all new code
- Write integration tests for API endpoints
- Write UI tests for frontend components

## Deployment

### Backend Services

The backend services are deployed to Kubernetes using the CI/CD pipeline. See the [CI/CD Pipeline](devops/ci-cd.md) documentation for details.

### Web Frontend

The web frontend is deployed to Kubernetes using the CI/CD pipeline. See the [CI/CD Pipeline](devops/ci-cd.md) documentation for details.

### Mobile App

The mobile app is released to the Google Play Store and Apple App Store using the CI/CD pipeline. See the [CI/CD Pipeline](devops/ci-cd.md) documentation for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes and commit
4. Push the branch and create a pull request
5. Address review comments

## License

This project is licensed under the MIT License - see the LICENSE file for details.
