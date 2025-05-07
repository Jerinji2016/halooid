import 'package:flutter/material.dart';

class AppTheme {
  // Light theme colors
  static const Color _lightPrimaryColor = Color(0xFF3F51B5);
  static const Color _lightAccentColor = Color(0xFF536DFE);
  static const Color _lightBackgroundColor = Color(0xFFF5F5F5);
  static const Color _lightCardColor = Colors.white;
  static const Color _lightTextColor = Color(0xFF212121);
  static const Color _lightSecondaryTextColor = Color(0xFF757575);

  // Dark theme colors
  static const Color _darkPrimaryColor = Color(0xFF3F51B5);
  static const Color _darkAccentColor = Color(0xFF536DFE);
  static const Color _darkBackgroundColor = Color(0xFF121212);
  static const Color _darkCardColor = Color(0xFF1E1E1E);
  static const Color _darkTextColor = Color(0xFFEEEEEE);
  static const Color _darkSecondaryTextColor = Color(0xFFAAAAAA);

  // Task priority colors
  static const Color lowPriorityColor = Color(0xFF4CAF50);
  static const Color mediumPriorityColor = Color(0xFFFFC107);
  static const Color highPriorityColor = Color(0xFFFF9800);
  static const Color urgentPriorityColor = Color(0xFFF44336);

  // Task status colors
  static const Color todoStatusColor = Color(0xFF9E9E9E);
  static const Color inProgressStatusColor = Color(0xFF2196F3);
  static const Color completedStatusColor = Color(0xFF4CAF50);
  static const Color blockedStatusColor = Color(0xFFF44336);

  // Light theme
  static final ThemeData lightTheme = ThemeData(
    useMaterial3: true,
    colorScheme: ColorScheme.light(
      primary: _lightPrimaryColor,
      secondary: _lightAccentColor,
      background: _lightBackgroundColor,
      surface: _lightCardColor,
      onPrimary: Colors.white,
      onSecondary: Colors.white,
      onBackground: _lightTextColor,
      onSurface: _lightTextColor,
    ),
    scaffoldBackgroundColor: _lightBackgroundColor,
    cardTheme: CardTheme(
      color: _lightCardColor,
      elevation: 2,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
      ),
    ),
    appBarTheme: AppBarTheme(
      backgroundColor: _lightPrimaryColor,
      foregroundColor: Colors.white,
      elevation: 0,
    ),
    textTheme: TextTheme(
      displayLarge: TextStyle(color: _lightTextColor),
      displayMedium: TextStyle(color: _lightTextColor),
      displaySmall: TextStyle(color: _lightTextColor),
      headlineMedium: TextStyle(color: _lightTextColor),
      headlineSmall: TextStyle(color: _lightTextColor),
      titleLarge: TextStyle(color: _lightTextColor),
      titleMedium: TextStyle(color: _lightTextColor),
      titleSmall: TextStyle(color: _lightTextColor),
      bodyLarge: TextStyle(color: _lightTextColor),
      bodyMedium: TextStyle(color: _lightTextColor),
      bodySmall: TextStyle(color: _lightSecondaryTextColor),
      labelLarge: TextStyle(color: _lightTextColor),
    ),
    floatingActionButtonTheme: FloatingActionButtonThemeData(
      backgroundColor: _lightAccentColor,
      foregroundColor: Colors.white,
    ),
    inputDecorationTheme: InputDecorationTheme(
      border: OutlineInputBorder(
        borderRadius: BorderRadius.circular(8),
      ),
      focusedBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(8),
        borderSide: BorderSide(color: _lightPrimaryColor, width: 2),
      ),
      filled: true,
      fillColor: Colors.white,
    ),
  );

  // Dark theme
  static final ThemeData darkTheme = ThemeData(
    useMaterial3: true,
    colorScheme: ColorScheme.dark(
      primary: _darkPrimaryColor,
      secondary: _darkAccentColor,
      background: _darkBackgroundColor,
      surface: _darkCardColor,
      onPrimary: Colors.white,
      onSecondary: Colors.white,
      onBackground: _darkTextColor,
      onSurface: _darkTextColor,
    ),
    scaffoldBackgroundColor: _darkBackgroundColor,
    cardTheme: CardTheme(
      color: _darkCardColor,
      elevation: 2,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
      ),
    ),
    appBarTheme: AppBarTheme(
      backgroundColor: _darkCardColor,
      foregroundColor: _darkTextColor,
      elevation: 0,
    ),
    textTheme: TextTheme(
      displayLarge: TextStyle(color: _darkTextColor),
      displayMedium: TextStyle(color: _darkTextColor),
      displaySmall: TextStyle(color: _darkTextColor),
      headlineMedium: TextStyle(color: _darkTextColor),
      headlineSmall: TextStyle(color: _darkTextColor),
      titleLarge: TextStyle(color: _darkTextColor),
      titleMedium: TextStyle(color: _darkTextColor),
      titleSmall: TextStyle(color: _darkTextColor),
      bodyLarge: TextStyle(color: _darkTextColor),
      bodyMedium: TextStyle(color: _darkTextColor),
      bodySmall: TextStyle(color: _darkSecondaryTextColor),
      labelLarge: TextStyle(color: _darkTextColor),
    ),
    floatingActionButtonTheme: FloatingActionButtonThemeData(
      backgroundColor: _darkAccentColor,
      foregroundColor: Colors.white,
    ),
    inputDecorationTheme: InputDecorationTheme(
      border: OutlineInputBorder(
        borderRadius: BorderRadius.circular(8),
      ),
      focusedBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(8),
        borderSide: BorderSide(color: _darkPrimaryColor, width: 2),
      ),
      filled: true,
      fillColor: Color(0xFF2C2C2C),
    ),
  );
}
