# Halooid CI/CD Pipeline

This directory contains the GitHub Actions workflows for the Halooid platform's CI/CD pipeline.

## Workflows

### Backend

- **Backend CI (`backend-ci.yml`)**: Runs linting, tests, and builds for the backend services.
- **Backend CD (`backend-cd.yml`)**: Builds and pushes Docker images for the backend services, and deploys them to Kubernetes when a tag is pushed.

### Web Frontend

- **Web CI (`web-ci.yml`)**: Runs linting, tests, and builds for the web frontend.
- **Web CD (`web-cd.yml`)**: Builds and pushes Docker images for the web frontend, and deploys them to Kubernetes when a tag is pushed.

### Mobile

- **Mobile CI (`mobile-ci.yml`)**: Runs linting, tests, and builds for the mobile app.
- **Mobile CD (`mobile-cd.yml`)**: Builds and releases the mobile app to the Google Play Store and Apple App Store when a tag is pushed.

## CI/CD Pipeline Flow

1. **Development**:
   - Developers work on feature branches
   - Push changes to GitHub
   - CI workflows run automatically to validate changes

2. **Pull Request**:
   - Create a pull request to merge changes into the main branch
   - CI workflows run to validate the changes
   - Code review is performed
   - PR is merged when approved and all checks pass

3. **Release**:
   - Create and push a tag (e.g., `v1.0.0`) to trigger the CD workflows
   - CD workflows build and push Docker images with the tag
   - CD workflows deploy the new version to Kubernetes
   - CD workflows release the mobile app to app stores

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

## Kubernetes Deployment

The CD workflows deploy the services to a Kubernetes cluster. The Kubernetes manifests are located in the `kubernetes` directory:

- `kubernetes/backend/`: Manifests for the backend services
- `kubernetes/web/`: Manifests for the web frontend

## Adding a New Service

To add a new service to the CI/CD pipeline:

1. Create a Dockerfile for the service
2. Add the service to the appropriate CI workflow
3. Add the service to the appropriate CD workflow
4. Create Kubernetes manifests for the service
5. Update the deployment steps in the CD workflow
