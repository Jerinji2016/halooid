# Halooid Mobile App

This directory contains the mobile applications for the Halooid platform, built with Flutter.

## Directory Structure

- `packages/`: Shared packages and widgets
- `apps/`: Flutter applications for each product

## Getting Started

### Prerequisites

- Flutter 3.0 or later
- Android Studio or Xcode for device emulation

### Setup

1. Install dependencies:
   ```bash
   flutter pub get
   ```

2. Run the app:
   ```bash
   flutter run
   ```

## Development

### Adding a New Widget

1. Create a new widget in the appropriate package in `packages/`
2. Export the widget from the package's library
3. Import and use the widget in your application

### Building for Production

Build the Android app with:
```bash
flutter build apk
```

Build the iOS app with:
```bash
flutter build ios
```

## Testing

Run tests with:
```bash
flutter test
```
