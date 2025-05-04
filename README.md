# Halooid Platform

Halooid is a multi-product platform consisting of five integrated products designed to help businesses manage their operations efficiently.

## Products

1. **Taskake**: Task management system
2. **Qultrix**: Human Resource Management Software
3. **AdminHub**: Internal monitoring and administration
4. **CustomerConnect**: CRM for customer interactions
5. **Invantray**: Inventory and asset management software

## Documentation

Comprehensive documentation is available in the [docs](./docs) directory.

## Technical Stack

### Backend

- **Language**: Go
- **API Approach**: Hybrid (REST with OpenAPI + gRPC for internal communication)
- **Database**: PostgreSQL with Redis for caching
- **Authentication**: JWT-based with role-based access control

### Web Frontend

- **Framework**: Svelte with SvelteKit
- **Styling**: Tailwind CSS
- **State Management**: Svelte stores
- **Build System**: Vite

### Mobile

- **Framework**: Flutter
- **State Management**: Provider/Riverpod
- **API Communication**: Dio

## Getting Started

### Prerequisites

- Go (version 1.20 or later)
- Node.js (version 18 or later)
- Flutter (version 3.0 or later)
- Docker and Docker Compose
- PostgreSQL (version 14 or later)
- Redis (version 6 or later)
- Git

### Installation

For detailed installation instructions, see the [Installation Guide](./docs/getting-started/setup.md).

Quick start with our setup script:

```bash
git clone https://github.com/Jerinji2016/halooid.git
cd halooid
./scripts/setup-dev-env.sh
docker-compose up -d
```

This will:

1. Check if all required tools are installed
2. Set up the development environment with Docker Compose
3. Start PostgreSQL and Redis services
4. Initialize the database schema

### Project Structure

```text
/
├── backend/                # Backend services
│   ├── cmd/                # Application entry points
│   ├── internal/           # Private application code
│   ├── pkg/                # Reusable packages
│   ├── api/                # API definitions
│   ├── scripts/            # Build and deployment scripts
│   └── configs/            # Configuration files
├── web/                    # Web frontend
│   ├── packages/           # Shared packages
│   └── apps/               # SvelteKit applications
├── mobile/                 # Mobile applications
│   ├── packages/           # Shared packages
│   └── apps/               # Flutter applications
└── docs/                   # Documentation
```

## Development

For development instructions, see the [Development Guide](./docs/development/index.md).

## Contributing

We welcome contributions to the Halooid platform! Please see our [Contributing Guide](./docs/contributing/index.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [Go](https://golang.org/)
- [Svelte](https://svelte.dev/)
- [SvelteKit](https://kit.svelte.dev/)
- [Flutter](https://flutter.dev/)
- [PostgreSQL](https://www.postgresql.org/)
- [Redis](https://redis.io/)
- [Tailwind CSS](https://tailwindcss.com/)
- [Docker](https://www.docker.com/)
- [MkDocs](https://www.mkdocs.org/)
- [Material for MkDocs](https://squidfunk.github.io/mkdocs-material/)
