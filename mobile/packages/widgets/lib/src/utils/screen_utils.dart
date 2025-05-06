import 'package:flutter/material.dart';

/// Utility functions for responsive design
class ScreenUtils {
  /// Returns true if the screen width is less than 600dp
  static bool isMobile(BuildContext context) {
    return MediaQuery.of(context).size.width < 600;
  }

  /// Returns true if the screen width is between 600dp and 840dp
  static bool isTablet(BuildContext context) {
    final width = MediaQuery.of(context).size.width;
    return width >= 600 && width < 840;
  }

  /// Returns true if the screen width is greater than or equal to 840dp
  static bool isDesktop(BuildContext context) {
    return MediaQuery.of(context).size.width >= 840;
  }

  /// Returns the screen width
  static double screenWidth(BuildContext context) {
    return MediaQuery.of(context).size.width;
  }

  /// Returns the screen height
  static double screenHeight(BuildContext context) {
    return MediaQuery.of(context).size.height;
  }

  /// Returns the screen orientation
  static Orientation orientation(BuildContext context) {
    return MediaQuery.of(context).orientation;
  }

  /// Returns true if the screen is in portrait orientation
  static bool isPortrait(BuildContext context) {
    return orientation(context) == Orientation.portrait;
  }

  /// Returns true if the screen is in landscape orientation
  static bool isLandscape(BuildContext context) {
    return orientation(context) == Orientation.landscape;
  }

  /// Returns the status bar height
  static double statusBarHeight(BuildContext context) {
    return MediaQuery.of(context).padding.top;
  }

  /// Returns the bottom inset (e.g., for keyboard)
  static double bottomInset(BuildContext context) {
    return MediaQuery.of(context).viewInsets.bottom;
  }

  /// Returns the safe area padding
  static EdgeInsets safeAreaPadding(BuildContext context) {
    return MediaQuery.of(context).padding;
  }

  /// Returns a responsive value based on screen width
  /// 
  /// Example:
  /// ```dart
  /// double fontSize = ScreenUtils.responsiveValue(
  ///   context,
  ///   mobile: 14,
  ///   tablet: 16,
  ///   desktop: 18,
  /// );
  /// ```
  static T responsiveValue<T>({
    required BuildContext context,
    required T mobile,
    required T tablet,
    required T desktop,
  }) {
    if (isDesktop(context)) {
      return desktop;
    } else if (isTablet(context)) {
      return tablet;
    } else {
      return mobile;
    }
  }
}
