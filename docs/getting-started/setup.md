# Development Environment Setup

This guide will help you set up your development environment for the Halooid platform.

## Prerequisites

Before you begin, make sure you have the following tools installed:

- **Go** (version 1.20 or later)
- **Node.js** (version 18 or later)
- **Flutter** (version 3.0 or later)
- **Docker** and **Docker Compose**
- **PostgreSQL** (version 14 or later)
- **Redis** (version 6 or later)
- **Git**

## Automated Setup

We provide a setup script that will check for required tools and help set up your development environment.

```bash
# Clone the repository
git clone https://github.com/Jerinji2016/halooid.git
cd halooid

# Run the setup script
./scripts/setup-dev-env.sh

# Start the development services
docker-compose up -d
```

The setup script will:

1. Check if all required tools are installed
2. Set up the Docker Compose environment
3. Prepare the project structure

## Manual Setup

If you prefer to set up your environment manually, follow these steps:

### 1. Clone the Repository

```bash
git clone https://github.com/Jerinji2016/halooid.git
cd halooid
```

### 2. Start the Development Services

```bash
docker-compose up -d
```

This will start PostgreSQL and Redis services in Docker containers.

### 3. Set Up the Backend

```bash
cd backend
go mod download
```

### 4. Set Up the Web Frontend

```bash
cd web
npm install
```

### 5. Set Up the Mobile App

```bash
cd mobile
flutter pub get
```

## Verifying Your Setup

To verify that your development environment is set up correctly:

1. Check that the Docker containers are running:

```bash
docker ps
```

You should see containers for PostgreSQL and Redis.

2. Check that you can connect to the database:

```bash
psql -h localhost -p 5432 -U halooid -d halooid
```

The password is `halooid_password`.

3. Check that you can connect to Redis:

```bash
redis-cli -h localhost -p 6379 ping
```

You should see `PONG` as the response.

## Next Steps

Now that your development environment is set up, you can:

- Explore the [project structure](../development/project-structure.md)
- Learn about the [architecture](../architecture/index.md)
- Start working on [your first feature](../development/first-feature.md)

## Troubleshooting

If you encounter any issues during setup, check the following:

- Make sure all required tools are installed and in your PATH
- Check that Docker is running
- Ensure that ports 5432 (PostgreSQL) and 6379 (Redis) are not already in use
- Check the Docker logs for any errors:

```bash
docker logs halooid-postgres
docker logs halooid-redis
```

If you still have issues, please [open an issue](https://github.com/Jerinji2016/halooid/issues) on GitHub.
