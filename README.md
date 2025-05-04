# Halooid Platform

Halooid is a multi-product platform consisting of five integrated products designed to help businesses manage their operations efficiently.

## Products

1. **Taskake**: Task management system
2. **Qultrix**: Human Resource Management Software
3. **AdminHub**: Internal monitoring and administration
4. **CustomerConnect**: CRM for customer interactions
5. **Invantray**: Inventory and asset management software

## Documentation

Comprehensive documentation is available at [https://yourusername.github.io/halooid](https://yourusername.github.io/halooid).

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

For detailed installation instructions, see the [Installation Guide](https://yourusername.github.io/halooid/getting-started/installation/).

Quick start with Docker Compose:

```bash
git clone https://github.com/yourusername/halooid.git
cd halooid
cp .env.example .env
# Edit .env file with your configuration
docker-compose up -d
```

## Development

For development instructions, see the [Development Guide](https://yourusername.github.io/halooid/development/).

## Contributing

We welcome contributions to the Halooid platform! Please see our [Contributing Guide](https://yourusername.github.io/halooid/contributing/) for more information.

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
