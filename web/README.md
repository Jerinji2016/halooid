# Halooid Web Frontend

This directory contains the web frontend applications for the Halooid platform, built with Svelte and SvelteKit.

## Directory Structure

- `packages/`: Shared packages and components
- `apps/`: SvelteKit applications for each product

## Getting Started

### Prerequisites

- Node.js 18 or later
- npm or yarn

### Setup

1. Install dependencies:
   ```bash
   npm install
   ```

2. Start the development server:
   ```bash
   npm run dev
   ```

## Development

### Adding a New Component

1. Create a new component in the appropriate package in `packages/`
2. Export the component from the package's index.js
3. Import and use the component in your application

### Building for Production

Build the application with:
```bash
npm run build
```

## Storybook

We use Storybook to document and test our UI components.

Start Storybook with:
```bash
npm run storybook
```
