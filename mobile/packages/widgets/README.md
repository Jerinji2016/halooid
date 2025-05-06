# Halooid Widgets

A shared widget library for Halooid mobile applications.

## Installation

Add this package to your `pubspec.yaml`:

```yaml
dependencies:
  halooid_widgets:
    path: ../packages/widgets
```

## Usage

```dart
import 'package:halooid_widgets/halooid_widgets.dart';

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Halooid App',
      theme: AppTheme.lightTheme(),
      darkTheme: AppTheme.darkTheme(),
      home: MyHomePage(),
    );
  }
}

class MyHomePage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Halooid App'),
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            AppCard(
              padding: EdgeInsets.all(16),
              child: Column(
                children: [
                  Text('Welcome to Halooid', style: Theme.of(context).textTheme.headlineMedium),
                  SizedBox(height: 16),
                  AppTextField(
                    label: 'Email',
                    hint: 'Enter your email',
                    prefixIcon: Icon(Icons.email),
                  ),
                  SizedBox(height: 16),
                  AppTextField(
                    label: 'Password',
                    hint: 'Enter your password',
                    prefixIcon: Icon(Icons.lock),
                    obscureText: true,
                  ),
                  SizedBox(height: 24),
                  AppButton(
                    onPressed: () {},
                    child: Text('Login'),
                    fullWidth: true,
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
```

## Components

### AppButton

A customizable button widget.

```dart
AppButton(
  onPressed: () {
    print('Button pressed');
  },
  variant: AppButtonVariant.primary,
  size: AppButtonSize.medium,
  child: Text('Click Me'),
)
```

### AppTextField

A text input widget with various styles and states.

```dart
AppTextField(
  controller: _emailController,
  label: 'Email',
  hint: 'Enter your email',
  prefixIcon: Icon(Icons.email),
  keyboardType: TextInputType.emailAddress,
  onChanged: (value) {
    print('Email changed: $value');
  },
)
```

### AppCard

A container widget for displaying content in a contained format.

```dart
AppCard(
  elevation: 2,
  borderRadius: 12,
  onTap: () {
    print('Card tapped');
  },
  child: Column(
    children: [
      AppCardHeader(
        title: Text('Card Title'),
        subtitle: Text('Card Subtitle'),
        leading: Icon(Icons.info),
        trailing: Icon(Icons.more_vert),
      ),
      AppCardContent(
        child: Text('This is the card content.'),
      ),
      AppCardFooter(
        child: Row(
          children: [
            AppButton(
              onPressed: () {},
              variant: AppButtonVariant.primary,
              size: AppButtonSize.small,
              child: Text('Action'),
            ),
            SizedBox(width: 8),
            AppButton(
              onPressed: () {},
              variant: AppButtonVariant.secondary,
              size: AppButtonSize.small,
              outlined: true,
              child: Text('Cancel'),
            ),
          ],
        ),
      ),
    ],
  ),
)
```

## Theme

The widget library comes with a default theme that can be customized.

```dart
// Use the default light theme
theme: AppTheme.lightTheme(),

// Use the default dark theme
darkTheme: AppTheme.darkTheme(),

// Customize the theme
theme: AppTheme.lightTheme().copyWith(
  colorScheme: AppTheme.lightTheme().colorScheme.copyWith(
    primary: Colors.red,
  ),
),
```

## Utilities

### ScreenUtils

Utility functions for responsive design.

```dart
// Check if the device is a mobile
if (ScreenUtils.isMobile(context)) {
  // Show mobile layout
}

// Get responsive value based on screen size
double fontSize = ScreenUtils.responsiveValue(
  context,
  mobile: 14,
  tablet: 16,
  desktop: 18,
);
```

### DateTimeUtils

Utility functions for date and time operations.

```dart
// Format a date
String formattedDate = DateTimeUtils.formatDate(
  DateTime.now(),
  format: 'MMM d, yyyy',
);

// Get relative time
String relativeTime = DateTimeUtils.getRelativeTime(
  DateTime.now().subtract(Duration(hours: 2)),
);
```
