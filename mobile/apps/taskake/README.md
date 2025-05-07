# Taskake Mobile App

Task management mobile application for the Halooid platform.

## Features

- View, create, edit, and delete tasks
- Filter and sort tasks by various criteria
- Track time spent on tasks
- Add comments to tasks
- Offline support for managing tasks without internet connection
- Synchronization of offline changes when back online

## Getting Started

### Prerequisites

- Flutter SDK (version 3.0.0 or higher)
- Dart SDK (version 3.0.0 or higher)
- Android Studio or VS Code with Flutter extensions

### Installation

1. Clone the repository:
   ```
   git clone https://github.com/Jerinji2016/halooid.git
   ```

2. Navigate to the Taskake app directory:
   ```
   cd halooid/mobile/apps/taskake
   ```

3. Install dependencies:
   ```
   flutter pub get
   ```

4. Generate the necessary files:
   ```
   flutter pub run build_runner build --delete-conflicting-outputs
   ```

5. Run the app:
   ```
   flutter run
   ```

## Project Structure

- `lib/main.dart` - Entry point of the application
- `lib/models/` - Data models
- `lib/screens/` - UI screens
- `lib/widgets/` - Reusable UI components
- `lib/services/` - Business logic and API communication
- `lib/utils/` - Utility functions and helpers
- `lib/theme/` - App theme and styling

## Dependencies

- `provider` - State management
- `http` - API communication
- `shared_preferences` - Local storage for offline support
- `flutter_secure_storage` - Secure storage for authentication tokens
- `json_annotation` - JSON serialization
- `intl` - Internationalization and date formatting

## Testing

Run the tests with:
```
flutter test
```

## Building for Production

Build an APK:
```
flutter build apk
```

Build an App Bundle:
```
flutter build appbundle
```

## Contributing

1. Create a feature branch
2. Commit your changes
3. Push to the branch
4. Create a new Pull Request
