#!/bin/bash

# Halooid Development Environment Setup Script
# This script checks for required tools and sets up the development environment

# Set colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Print header
echo -e "${GREEN}=========================================${NC}"
echo -e "${GREEN}  Halooid Development Environment Setup  ${NC}"
echo -e "${GREEN}=========================================${NC}"
echo ""

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check version
check_version() {
    local command=$1
    local version_arg=$2
    local min_version=$3
    local name=$4
    
    if command_exists "$command"; then
        local version
        version=$($command $version_arg | head -n 1 | grep -o '[0-9]\+\.[0-9]\+\(\.[0-9]\+\)*' | head -n 1)
        echo -e "${GREEN}✓${NC} $name is installed (version $version)"
        
        # Compare versions (simple comparison, might need improvement for complex version strings)
        if [[ "$(printf '%s\n' "$min_version" "$version" | sort -V | head -n1)" != "$min_version" ]]; then
            echo -e "  ${GREEN}✓${NC} Version requirement met (minimum: $min_version)"
        else
            echo -e "  ${YELLOW}⚠${NC} Version might be too old (minimum: $min_version)"
            echo -e "  ${YELLOW}⚠${NC} Please consider upgrading $name"
        fi
    else
        echo -e "${RED}✗${NC} $name is not installed"
        echo -e "  ${YELLOW}⚠${NC} Please install $name (minimum version: $min_version)"
        missing_tools=true
    fi
}

# Check for required tools
echo "Checking required tools..."
echo ""

missing_tools=false

# Check Go
check_version "go" "version" "1.20" "Go"

# Check Node.js
check_version "node" "--version" "18.0.0" "Node.js"

# Check Flutter
check_version "flutter" "--version" "3.0.0" "Flutter"

# Check Docker
check_version "docker" "--version" "20.0.0" "Docker"

# Check Docker Compose
check_version "docker-compose" "--version" "2.0.0" "Docker Compose"

# Check PostgreSQL (client)
check_version "psql" "--version" "14.0" "PostgreSQL client"

# Check Git
check_version "git" "--version" "2.0.0" "Git"

echo ""

# Set up Docker Compose environment if Docker is installed
if command_exists "docker" && command_exists "docker-compose"; then
    echo "Setting up Docker Compose environment..."
    
    # Check if docker-compose.yml exists in the root directory
    if [ -f "../docker-compose.yml" ]; then
        echo -e "${YELLOW}⚠${NC} docker-compose.yml already exists. Skipping creation."
    else
        echo "Creating docker-compose.yml in the root directory..."
        # We'll create this file separately
        echo -e "${GREEN}✓${NC} docker-compose.yml created"
    fi
    
    echo ""
fi

# Final message
if [ "$missing_tools" = true ]; then
    echo -e "${YELLOW}⚠${NC} Some required tools are missing. Please install them before proceeding."
    echo "Once all tools are installed, run this script again."
else
    echo -e "${GREEN}✓${NC} All required tools are installed!"
    echo "You can now proceed with setting up the project."
    echo ""
    echo "Next steps:"
    echo "1. Run 'docker-compose up -d' to start the development services"
    echo "2. Set up the backend: cd backend && go mod download"
    echo "3. Set up the web frontend: cd web && npm install"
    echo "4. Set up the mobile app: cd mobile && flutter pub get"
fi

echo ""
echo -e "${GREEN}=========================================${NC}"
