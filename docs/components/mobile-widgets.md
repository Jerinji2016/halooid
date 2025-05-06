# Mobile Widget Library

## Overview

The Halooid platform includes a shared widget library for mobile applications. This library provides a consistent look and feel across all mobile interfaces in the platform.

## Technology Stack

- **Flutter**: UI framework
- **Dart**: Programming language
- **Provider**: State management
- **Material Design**: Design system

## Installation

Add the package to your `pubspec.yaml`:

```yaml
dependencies:
  halooid_widgets:
    path: ../packages/widgets
```

Then run:

```bash
flutter pub get
```

## Widgets

### AppButton

A customizable button widget with various styles and states.

```dart
import 'package:halooid_widgets/halooid_widgets.dart';

// Primary button
AppButton(
  onPressed: () => print('Button pressed'),
  child: Text('Click Me'),
);

// Secondary button
AppButton(
  onPressed: () => print('Button pressed'),
  variant: AppButtonVariant.secondary,
  child: Text('Secondary Button'),
);

// Outlined button
AppButton(
  onPressed: () => print('Button pressed'),
  outlined: true,
  child: Text('Outlined Button'),
);

// Disabled button
AppButton(
  onPressed: null,
  child: Text('Disabled Button'),
);

// Loading button
AppButton(
  onPressed: () => print('Button pressed'),
  loading: true,
  child: Text('Loading...'),
);

// Button with icon
AppButton(
  onPressed: () => print('Button pressed'),
  icon: Icon(Icons.add),
  child: Text('Add Item'),
);
```

#### Properties

- `onPressed`: Button press callback
- `child`: Button content
- `variant`: Button style variant (`primary`, `secondary`, `success`, `danger`, `warning`, `info`)
- `size`: Button size (`small`, `medium`, `large`)
- `outlined`: Whether the button is outlined
- `rounded`: Whether the button is rounded
- `disabled`: Whether the button is disabled
- `fullWidth`: Whether the button is full width
- `loading`: Whether the button is in a loading state
- `icon`: Optional icon to display
- `iconPosition`: Position of the icon (`left`, `right`)

### AppTextField

A text input widget with various styles and states.

```dart
import 'package:halooid_widgets/halooid_widgets.dart';

// Basic text field
AppTextField(
  label: 'Email',
  hint: 'Enter your email',
);

// Text field with error
AppTextField(
  label: 'Email',
  hint: 'Enter your email',
  errorText: 'Please enter a valid email address',
);

// Text field with helper text
AppTextField(
  label: 'Email',
  hint: 'Enter your email',
  helperText: 'We\'ll never share your email with anyone else',
);

// Text field with icon
AppTextField(
  label: 'Email',
  hint: 'Enter your email',
  prefixIcon: Icon(Icons.email),
);
```

#### Properties

- `controller`: Text editing controller
- `focusNode`: Focus node
- `decoration`: Input decoration
- `label`: Field label
- `hint`: Field hint
- `helperText`: Helper text
- `errorText`: Error text
- `prefixIcon`: Icon to display at the start of the input
- `suffixIcon`: Icon to display at the end of the input
- `obscureText`: Whether to obscure text (for passwords)
- `keyboardType`: Keyboard type
- `textInputAction`: Text input action
- `onChanged`: Change handler function
- `onSubmitted`: Submit handler function
- `enabled`: Whether the text field is enabled
- `borderRadius`: Border radius
- `filled`: Whether the field is filled
- `fillColor`: Fill color
- `borderColor`: Border color
- `focusedBorderColor`: Focused border color
- `errorBorderColor`: Error border color

### AppCard

A container widget for displaying content in a contained format.

```dart
import 'package:halooid_widgets/halooid_widgets.dart';

AppCard(
  elevation: 2,
  borderRadius: 12,
  onTap: () => print('Card tapped'),
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
);
```

#### Properties

- `child`: Card content
- `elevation`: Card elevation
- `color`: Card color
- `shadowColor`: Shadow color
- `borderRadius`: Border radius
- `borderColor`: Border color
- `margin`: Card margin
- `padding`: Card padding
- `clipBehavior`: Clip behavior
- `onTap`: Callback when card is tapped
- `onLongPress`: Callback when card is long pressed
- `hoverable`: Whether the card is hoverable

## Theme

The widget library includes a theme system that provides consistent colors, typography, and styling across all widgets.

```dart
import 'package:flutter/material.dart';
import 'package:halooid_widgets/halooid_widgets.dart';

void main() {
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Halooid App',
      theme: AppTheme.lightTheme(),
      darkTheme: AppTheme.darkTheme(),
      themeMode: ThemeMode.system,
      home: MyHomePage(),
    );
  }
}
```

### Theme Properties

- `colorScheme`: Color scheme
- `textTheme`: Typography settings
- `appBarTheme`: App bar styling
- `cardTheme`: Card styling
- `elevatedButtonTheme`: Elevated button styling
- `outlinedButtonTheme`: Outlined button styling
- `textButtonTheme`: Text button styling
- `inputDecorationTheme`: Input decoration styling

## Utilities

### ScreenUtils

Utility functions for responsive design.

```dart
import 'package:halooid_widgets/halooid_widgets.dart';

// Check if the device is a mobile
if (ScreenUtils.isMobile(context)) {
  // Show mobile layout
}

// Check if the device is a tablet
if (ScreenUtils.isTablet(context)) {
  // Show tablet layout
}

// Get responsive value based on screen size
double fontSize = ScreenUtils.responsiveValue(
  context,
  mobile: 14,
  tablet: 16,
  desktop: 18,
);

// Get screen width
double width = ScreenUtils.screenWidth(context);

// Get screen height
double height = ScreenUtils.screenHeight(context);

// Check orientation
if (ScreenUtils.isPortrait(context)) {
  // Portrait layout
} else {
  // Landscape layout
}
```

### DateTimeUtils

Utility functions for date and time operations.

```dart
import 'package:halooid_widgets/halooid_widgets.dart';

// Format a date
String formattedDate = DateTimeUtils.formatDate(
  DateTime.now(),
  format: 'MMM d, yyyy',
);

// Format a time
String formattedTime = DateTimeUtils.formatTime(
  DateTime.now(),
  format: 'h:mm a',
);

// Format a date and time
String formattedDateTime = DateTimeUtils.formatDateTime(
  DateTime.now(),
  format: 'MMM d, yyyy h:mm a',
);

// Get relative time
String relativeTime = DateTimeUtils.getRelativeTime(
  DateTime.now().subtract(Duration(hours: 2)),
);

// Check if a date is today
bool isToday = DateTimeUtils.isToday(DateTime.now());

// Check if a date is yesterday
bool isYesterday = DateTimeUtils.isYesterday(
  DateTime.now().subtract(Duration(days: 1)),
);
```

## Development

### Adding a New Widget

1. Create a new file in `lib/src/widgets` with the widget name
2. Implement the widget
3. Add tests in the `test` directory
4. Export the widget in `lib/src/widgets/index.dart`
5. Update the main library file if necessary

### Running Tests

```bash
# Run tests
flutter test

# Run tests with coverage
flutter test --coverage
```

### Building the Library

```bash
# Build the library
flutter pub run build_runner build
```

## Best Practices

1. **Consistency**: Follow the established patterns and styles
2. **Accessibility**: Ensure widgets are accessible
3. **Performance**: Keep widgets lightweight and efficient
4. **Documentation**: Document all properties and usage examples
5. **Testing**: Write tests for all widgets
