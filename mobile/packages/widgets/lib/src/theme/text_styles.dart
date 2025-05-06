import 'package:flutter/material.dart';
import 'app_colors.dart';

/// Text styles for the app
class TextStyles {
  /// Light theme text styles
  static const TextTheme textTheme = TextTheme(
    displayLarge: TextStyle(
      fontFamily: 'Inter',
      fontSize: 57,
      fontWeight: FontWeight.w400,
      letterSpacing: -0.25,
      color: AppColors.textPrimary,
    ),
    displayMedium: TextStyle(
      fontFamily: 'Inter',
      fontSize: 45,
      fontWeight: FontWeight.w400,
      color: AppColors.textPrimary,
    ),
    displaySmall: TextStyle(
      fontFamily: 'Inter',
      fontSize: 36,
      fontWeight: FontWeight.w400,
      color: AppColors.textPrimary,
    ),
    headlineLarge: TextStyle(
      fontFamily: 'Inter',
      fontSize: 32,
      fontWeight: FontWeight.w700,
      color: AppColors.textPrimary,
    ),
    headlineMedium: TextStyle(
      fontFamily: 'Inter',
      fontSize: 28,
      fontWeight: FontWeight.w700,
      color: AppColors.textPrimary,
    ),
    headlineSmall: TextStyle(
      fontFamily: 'Inter',
      fontSize: 24,
      fontWeight: FontWeight.w700,
      color: AppColors.textPrimary,
    ),
    titleLarge: TextStyle(
      fontFamily: 'Inter',
      fontSize: 22,
      fontWeight: FontWeight.w600,
      color: AppColors.textPrimary,
    ),
    titleMedium: TextStyle(
      fontFamily: 'Inter',
      fontSize: 16,
      fontWeight: FontWeight.w600,
      letterSpacing: 0.15,
      color: AppColors.textPrimary,
    ),
    titleSmall: TextStyle(
      fontFamily: 'Inter',
      fontSize: 14,
      fontWeight: FontWeight.w600,
      letterSpacing: 0.1,
      color: AppColors.textPrimary,
    ),
    bodyLarge: TextStyle(
      fontFamily: 'Inter',
      fontSize: 16,
      fontWeight: FontWeight.w400,
      letterSpacing: 0.15,
      color: AppColors.textPrimary,
    ),
    bodyMedium: TextStyle(
      fontFamily: 'Inter',
      fontSize: 14,
      fontWeight: FontWeight.w400,
      letterSpacing: 0.25,
      color: AppColors.textPrimary,
    ),
    bodySmall: TextStyle(
      fontFamily: 'Inter',
      fontSize: 12,
      fontWeight: FontWeight.w400,
      letterSpacing: 0.4,
      color: AppColors.textSecondary,
    ),
    labelLarge: TextStyle(
      fontFamily: 'Inter',
      fontSize: 14,
      fontWeight: FontWeight.w500,
      letterSpacing: 0.1,
      color: AppColors.textPrimary,
    ),
    labelMedium: TextStyle(
      fontFamily: 'Inter',
      fontSize: 12,
      fontWeight: FontWeight.w500,
      letterSpacing: 0.5,
      color: AppColors.textPrimary,
    ),
    labelSmall: TextStyle(
      fontFamily: 'Inter',
      fontSize: 11,
      fontWeight: FontWeight.w500,
      letterSpacing: 0.5,
      color: AppColors.textSecondary,
    ),
  );

  /// Dark theme text styles
  static final TextTheme darkTextTheme = textTheme.copyWith(
    displayLarge: textTheme.displayLarge?.copyWith(color: Colors.white),
    displayMedium: textTheme.displayMedium?.copyWith(color: Colors.white),
    displaySmall: textTheme.displaySmall?.copyWith(color: Colors.white),
    headlineLarge: textTheme.headlineLarge?.copyWith(color: Colors.white),
    headlineMedium: textTheme.headlineMedium?.copyWith(color: Colors.white),
    headlineSmall: textTheme.headlineSmall?.copyWith(color: Colors.white),
    titleLarge: textTheme.titleLarge?.copyWith(color: Colors.white),
    titleMedium: textTheme.titleMedium?.copyWith(color: Colors.white),
    titleSmall: textTheme.titleSmall?.copyWith(color: Colors.white),
    bodyLarge: textTheme.bodyLarge?.copyWith(color: Colors.white),
    bodyMedium: textTheme.bodyMedium?.copyWith(color: Colors.white),
    bodySmall: textTheme.bodySmall?.copyWith(color: Colors.white70),
    labelLarge: textTheme.labelLarge?.copyWith(color: Colors.white),
    labelMedium: textTheme.labelMedium?.copyWith(color: Colors.white),
    labelSmall: textTheme.labelSmall?.copyWith(color: Colors.white70),
  );
}
