# Installation Guide

This guide provides instructions for installing the Halooid platform in various environments.

## Installation Options

The Halooid platform can be installed in several ways:

1. **Cloud Hosted (SaaS)**: Use our cloud-hosted solution at [halooid.com](https://halooid.com).
2. **Self-Hosted with Docker Compose**: The simplest way to self-host Halooid.
3. **Self-Hosted with Kubernetes**: For more advanced deployments with high availability and scalability.
4. **Manual Installation**: For complete control over the installation process.

## Cloud Hosted (SaaS)

The easiest way to use Halooid is through our cloud-hosted solution at [halooid.com](https://halooid.com). This option requires no installation or maintenance on your part.

To get started with the cloud-hosted version:

1. Visit [halooid.com](https://halooid.com) and sign up for an account.
2. Choose the products you want to use.
3. Set up your organization and invite team members.
4. Start using the platform.

## Self-Hosted with Docker Compose

Docker Compose is the simplest way to self-host the Halooid platform. This option is suitable for small to medium-sized deployments.

### Prerequisites

- Docker (version 20.10 or later)
- Docker Compose (version 2.0 or later)
- 4GB RAM (minimum)
- 20GB disk space (minimum)
- Internet connection

### Installation Steps

1. Clone the Halooid repository:
   ```bash
   git clone https://github.com/yourusername/halooid.git
   cd halooid
   ```

2. Configure the environment variables:
   ```bash
   cp .env.example .env
   # Edit the .env file with your preferred text editor
   ```

3. Start the services:
   ```bash
   docker-compose up -d
   ```

4. Access the platform:
   Open your web browser and navigate to `http://localhost:8080`.

5. Create an admin account:
   ```bash
   docker-compose exec backend ./create-admin.sh
   ```

### Configuration

The Halooid platform can be configured through environment variables in the `.env` file. Here are some important configuration options:

- `HALOOID_DB_HOST`: PostgreSQL host
- `HALOOID_DB_PORT`: PostgreSQL port
- `HALOOID_DB_NAME`: PostgreSQL database name
- `HALOOID_DB_USER`: PostgreSQL username
- `HALOOID_DB_PASSWORD`: PostgreSQL password
- `HALOOID_REDIS_HOST`: Redis host
- `HALOOID_REDIS_PORT`: Redis port
- `HALOOID_SECRET_KEY`: Secret key for JWT tokens
- `HALOOID_ADMIN_EMAIL`: Admin email address
- `HALOOID_ADMIN_PASSWORD`: Admin password

For a complete list of configuration options, see the [Configuration Guide](../deployment/backend.md#configuration).

## Self-Hosted with Kubernetes

For larger deployments with high availability and scalability requirements, we recommend using Kubernetes.

### Prerequisites

- Kubernetes cluster (version 1.20 or later)
- Helm (version 3.0 or later)
- kubectl (configured to access your cluster)
- 8GB RAM (minimum)
- 40GB disk space (minimum)
- Internet connection

### Installation Steps

1. Add the Halooid Helm repository:
   ```bash
   helm repo add halooid https://charts.halooid.com
   helm repo update
   ```

2. Create a values file:
   ```bash
   cp values.yaml.example values.yaml
   # Edit the values.yaml file with your preferred text editor
   ```

3. Install the Halooid chart:
   ```bash
   helm install halooid halooid/halooid -f values.yaml
   ```

4. Access the platform:
   ```bash
   kubectl get svc halooid-frontend
   ```
   Open your web browser and navigate to the external IP address.

5. Create an admin account:
   ```bash
   kubectl exec -it $(kubectl get pods -l app=halooid-backend -o jsonpath="{.items[0].metadata.name}") -- ./create-admin.sh
   ```

### Configuration

The Halooid platform can be configured through the `values.yaml` file. Here are some important configuration options:

- `global.environment`: Environment (production, staging, development)
- `global.domain`: Domain name for the platform
- `global.tls.enabled`: Enable TLS
- `postgresql.enabled`: Enable PostgreSQL deployment
- `redis.enabled`: Enable Redis deployment
- `backend.replicas`: Number of backend replicas
- `frontend.replicas`: Number of frontend replicas

For a complete list of configuration options, see the [Kubernetes Deployment Guide](../deployment/backend.md#kubernetes).

## Manual Installation

For complete control over the installation process, you can install the Halooid platform manually.

### Prerequisites

- Go (version 1.20 or later)
- Node.js (version 18 or later)
- PostgreSQL (version 14 or later)
- Redis (version 6 or later)
- 4GB RAM (minimum)
- 20GB disk space (minimum)
- Internet connection

### Installation Steps

1. Clone the Halooid repository:
   ```bash
   git clone https://github.com/yourusername/halooid.git
   cd halooid
   ```

2. Install backend dependencies:
   ```bash
   cd backend
   go mod download
   ```

3. Build the backend:
   ```bash
   go build -o halooid-backend ./cmd/server
   ```

4. Install frontend dependencies:
   ```bash
   cd ../web
   npm install
   ```

5. Build the frontend:
   ```bash
   npm run build
   ```

6. Set up the database:
   ```bash
   cd ../scripts
   ./setup-database.sh
   ```

7. Configure the environment variables:
   ```bash
   cp .env.example .env
   # Edit the .env file with your preferred text editor
   ```

8. Start the backend:
   ```bash
   cd ../backend
   ./halooid-backend
   ```

9. Start the frontend:
   ```bash
   cd ../web
   npm run start
   ```

10. Access the platform:
    Open your web browser and navigate to `http://localhost:3000`.

11. Create an admin account:
    ```bash
    cd ../scripts
    ./create-admin.sh
    ```

### Configuration

The Halooid platform can be configured through environment variables. Here are some important configuration options:

- `HALOOID_DB_HOST`: PostgreSQL host
- `HALOOID_DB_PORT`: PostgreSQL port
- `HALOOID_DB_NAME`: PostgreSQL database name
- `HALOOID_DB_USER`: PostgreSQL username
- `HALOOID_DB_PASSWORD`: PostgreSQL password
- `HALOOID_REDIS_HOST`: Redis host
- `HALOOID_REDIS_PORT`: Redis port
- `HALOOID_SECRET_KEY`: Secret key for JWT tokens
- `HALOOID_ADMIN_EMAIL`: Admin email address
- `HALOOID_ADMIN_PASSWORD`: Admin password

For a complete list of configuration options, see the [Configuration Guide](../deployment/backend.md#configuration).

## Troubleshooting

### Common Issues

#### Database Connection Issues

If you encounter database connection issues, check the following:

1. Ensure that PostgreSQL is running and accessible.
2. Verify that the database credentials in the configuration are correct.
3. Check that the database has been created and initialized.

#### Redis Connection Issues

If you encounter Redis connection issues, check the following:

1. Ensure that Redis is running and accessible.
2. Verify that the Redis connection details in the configuration are correct.

#### Web Server Issues

If you encounter web server issues, check the following:

1. Ensure that the backend server is running.
2. Verify that the frontend has been built correctly.
3. Check that the web server port is not being used by another application.

### Getting Help

If you encounter issues that are not covered in this guide, you can:

1. Check the [FAQ](../faq.md) for common questions and answers.
2. Search for similar issues in the [GitHub Issues](https://github.com/yourusername/halooid/issues).
3. Ask for help in the [GitHub Discussions](https://github.com/yourusername/halooid/discussions).
4. Contact our support team at [support@halooid.com](mailto:support@halooid.com).

## Next Steps

After installing the Halooid platform, you can:

1. [Configure the platform](../deployment/backend.md#configuration) according to your needs.
2. [Set up your organization](quick-start.md#setting-up-your-organization) and invite team members.
3. [Explore the features](../products/index.md) of each product.
4. [Integrate with other systems](../api-reference/index.md) using the API.
