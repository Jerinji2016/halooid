# CI/CD Pipeline

## Overview

The Halooid platform uses a comprehensive CI/CD pipeline to automate the building, testing, and deployment of all components. This document describes the CI/CD pipeline setup and workflow.

## Technology Stack

- **GitHub Actions**: CI/CD platform
- **Docker**: Containerization
- **Kubernetes**: Container orchestration
- **AWS EKS**: Managed Kubernetes service
- **GitHub Container Registry**: Container registry

## Pipeline Components

### Backend CI/CD

The backend CI/CD pipeline is responsible for building, testing, and deploying the Go microservices.

#### CI Workflow (`backend-ci.yml`)

This workflow runs on every push to the `main` branch and on pull requests targeting the `main` branch that affect backend code.

**Steps**:
1. **Lint**: Runs `golangci-lint` to check code quality
2. **Test**: Runs unit and integration tests with coverage reporting
3. **Build**: Builds the backend services

#### CD Workflow (`backend-cd.yml`)

This workflow runs on every push to the `main` branch and when tags are pushed.

**Steps**:
1. **Build and Push Docker Images**: Builds Docker images for each service and pushes them to GitHub Container Registry
2. **Deploy to Kubernetes**: When a tag is pushed, deploys the services to the Kubernetes cluster

### Web Frontend CI/CD

The web frontend CI/CD pipeline is responsible for building, testing, and deploying the web application.

#### CI Workflow (`web-ci.yml`)

This workflow runs on every push to the `main` branch and on pull requests targeting the `main` branch that affect web code.

**Steps**:
1. **Lint**: Runs ESLint to check code quality
2. **Test**: Runs unit and integration tests with coverage reporting
3. **Build**: Builds the web application

#### CD Workflow (`web-cd.yml`)

This workflow runs on every push to the `main` branch and when tags are pushed.

**Steps**:
1. **Build and Push Docker Image**: Builds a Docker image for the web frontend and pushes it to GitHub Container Registry
2. **Deploy to Kubernetes**: When a tag is pushed, deploys the web frontend to the Kubernetes cluster

### Mobile CI/CD

The mobile CI/CD pipeline is responsible for building, testing, and releasing the mobile application.

#### CI Workflow (`mobile-ci.yml`)

This workflow runs on every push to the `main` branch and on pull requests targeting the `main` branch that affect mobile code.

**Steps**:
1. **Lint**: Runs Flutter analyze to check code quality
2. **Test**: Runs unit and widget tests with coverage reporting
3. **Build Android**: Builds the Android APK
4. **Build iOS**: Builds the iOS app

#### CD Workflow (`mobile-cd.yml`)

This workflow runs when tags are pushed.

**Steps**:
1. **Build Android**: Builds the Android APK and App Bundle
2. **Release to Google Play**: Uploads the App Bundle to Google Play Store (internal track)
3. **Build iOS**: Builds the iOS app
4. **Release to App Store**: Uploads the iOS app to App Store Connect (TestFlight)

## Workflow Triggers

- **Push to Main**: Triggers CI workflows and CD workflows (build and push only)
- **Pull Request**: Triggers CI workflows
- **Tag Push**: Triggers CD workflows (build, push, and deploy)

## Environments

### Development

- Local development environment
- Docker Compose for local services
- Mock services for external dependencies

### Testing

- Kubernetes namespace: `halooid-test`
- Automated testing through CI pipeline
- Test data and configurations

### Staging

- Kubernetes namespace: `halooid-staging`
- Mirrors production environment
- Used for final testing before production deployment

### Production

- Kubernetes namespace: `halooid`
- High availability configuration
- Auto-scaling based on demand
- Regular backups

## Kubernetes Deployment

The CD workflows deploy the services to a Kubernetes cluster on AWS EKS.

### Kubernetes Manifests

- **Backend**: `kubernetes/backend/`
  - API Gateway
  - Authentication Service
  - RBAC Service
  - Database Services (PostgreSQL, Redis)

- **Web Frontend**: `kubernetes/web/`
  - Web Frontend

### Deployment Process

1. Update the Kubernetes manifests with the new image tag
2. Apply the Kubernetes manifests to the cluster
3. Wait for the deployment to complete
4. Verify the deployment

## Required Secrets

The following secrets need to be configured in the GitHub repository settings:

### AWS Secrets (for Kubernetes deployment)

- `AWS_ACCESS_KEY_ID`: AWS access key ID
- `AWS_SECRET_ACCESS_KEY`: AWS secret access key
- `AWS_REGION`: AWS region

### Android Secrets

- `KEYSTORE_BASE64`: Base64-encoded Android keystore file
- `KEYSTORE_PASSWORD`: Password for the keystore
- `KEY_PASSWORD`: Password for the key
- `KEY_ALIAS`: Alias of the key
- `GOOGLE_PLAY_SERVICE_ACCOUNT_JSON`: Google Play service account JSON

### iOS Secrets

- `IOS_P12_CERTIFICATE_BASE64`: Base64-encoded iOS P12 certificate
- `IOS_P12_CERTIFICATE_PASSWORD`: Password for the P12 certificate
- `KEYCHAIN_PASSWORD`: Password for the keychain
- `IOS_PROVISIONING_PROFILE_BASE64`: Base64-encoded iOS provisioning profile
- `APPLE_ID`: Apple ID
- `APPLE_APP_SPECIFIC_PASSWORD`: App-specific password for the Apple ID
- `APPLE_TEAM_ID`: Apple team ID

## Monitoring and Alerting

The CI/CD pipeline includes monitoring and alerting for deployment failures:

1. **GitHub Actions Notifications**: Email notifications for workflow failures
2. **Slack Notifications**: Notifications in the development Slack channel
3. **Deployment Monitoring**: Monitoring of deployment status in Kubernetes

## Best Practices

1. **Trunk-Based Development**: Use short-lived feature branches and merge frequently
2. **Automated Testing**: Ensure comprehensive test coverage
3. **Infrastructure as Code**: Manage all infrastructure through code
4. **Immutable Deployments**: Use immutable containers and avoid in-place updates
5. **Rollback Plan**: Always have a plan for rolling back deployments
6. **Security Scanning**: Include security scanning in the pipeline
7. **Artifact Versioning**: Use semantic versioning for all artifacts

## Troubleshooting

### Common Issues

1. **Failed Tests**: Check the test logs for details on which tests failed and why
2. **Build Failures**: Check the build logs for compilation errors or dependency issues
3. **Deployment Failures**: Check the Kubernetes events and logs for deployment issues
4. **Permission Issues**: Ensure the service account has the necessary permissions

### Debugging Steps

1. **Check Workflow Logs**: Review the GitHub Actions workflow logs
2. **Check Container Logs**: Review the container logs in Kubernetes
3. **Check Kubernetes Events**: Review the Kubernetes events for deployment issues
4. **Check Service Status**: Verify the service status in Kubernetes
5. **Check Network Connectivity**: Verify network connectivity between services
