import 'package:intl/intl.dart';

/// Utility functions for date and time operations
class DateTimeUtils {
  /// Formats a date as a string using the specified format
  /// 
  /// Example:
  /// ```dart
  /// String formattedDate = DateTimeUtils.formatDate(
  ///   DateTime.now(),
  ///   format: 'MMM d, yyyy',
  /// );
  /// ```
  static String formatDate(DateTime date, {String format = 'yyyy-MM-dd'}) {
    return DateFormat(format).format(date);
  }

  /// Formats a time as a string using the specified format
  /// 
  /// Example:
  /// ```dart
  /// String formattedTime = DateTimeUtils.formatTime(
  ///   DateTime.now(),
  ///   format: 'h:mm a',
  /// );
  /// ```
  static String formatTime(DateTime time, {String format = 'HH:mm'}) {
    return DateFormat(format).format(time);
  }

  /// Formats a date and time as a string using the specified format
  /// 
  /// Example:
  /// ```dart
  /// String formattedDateTime = DateTimeUtils.formatDateTime(
  ///   DateTime.now(),
  ///   format: 'MMM d, yyyy h:mm a',
  /// );
  /// ```
  static String formatDateTime(DateTime dateTime, {String format = 'yyyy-MM-dd HH:mm'}) {
    return DateFormat(format).format(dateTime);
  }

  /// Returns a relative time string (e.g., "2 hours ago")
  /// 
  /// Example:
  /// ```dart
  /// String relativeTime = DateTimeUtils.getRelativeTime(
  ///   DateTime.now().subtract(Duration(hours: 2)),
  /// );
  /// ```
  static String getRelativeTime(DateTime dateTime) {
    final now = DateTime.now();
    final difference = now.difference(dateTime);

    if (difference.inSeconds < 60) {
      return 'Just now';
    } else if (difference.inMinutes < 60) {
      return '${difference.inMinutes} ${difference.inMinutes == 1 ? 'minute' : 'minutes'} ago';
    } else if (difference.inHours < 24) {
      return '${difference.inHours} ${difference.inHours == 1 ? 'hour' : 'hours'} ago';
    } else if (difference.inDays < 7) {
      return '${difference.inDays} ${difference.inDays == 1 ? 'day' : 'days'} ago';
    } else if (difference.inDays < 30) {
      final weeks = (difference.inDays / 7).floor();
      return '$weeks ${weeks == 1 ? 'week' : 'weeks'} ago';
    } else if (difference.inDays < 365) {
      final months = (difference.inDays / 30).floor();
      return '$months ${months == 1 ? 'month' : 'months'} ago';
    } else {
      final years = (difference.inDays / 365).floor();
      return '$years ${years == 1 ? 'year' : 'years'} ago';
    }
  }

  /// Returns the start of the day for the given date
  /// 
  /// Example:
  /// ```dart
  /// DateTime startOfDay = DateTimeUtils.startOfDay(DateTime.now());
  /// ```
  static DateTime startOfDay(DateTime date) {
    return DateTime(date.year, date.month, date.day);
  }

  /// Returns the end of the day for the given date
  /// 
  /// Example:
  /// ```dart
  /// DateTime endOfDay = DateTimeUtils.endOfDay(DateTime.now());
  /// ```
  static DateTime endOfDay(DateTime date) {
    return DateTime(date.year, date.month, date.day, 23, 59, 59, 999);
  }

  /// Returns the start of the week for the given date
  /// 
  /// Example:
  /// ```dart
  /// DateTime startOfWeek = DateTimeUtils.startOfWeek(
  ///   DateTime.now(),
  ///   startOfWeek: DateTime.monday,
  /// );
  /// ```
  static DateTime startOfWeek(DateTime date, {int startOfWeek = DateTime.monday}) {
    final day = date.weekday;
    return startOfDay(date.subtract(Duration(days: (day - startOfWeek) % 7)));
  }

  /// Returns the end of the week for the given date
  /// 
  /// Example:
  /// ```dart
  /// DateTime endOfWeek = DateTimeUtils.endOfWeek(
  ///   DateTime.now(),
  ///   startOfWeek: DateTime.monday,
  /// );
  /// ```
  static DateTime endOfWeek(DateTime date, {int startOfWeek = DateTime.monday}) {
    final startDay = startOfWeek(date, startOfWeek: startOfWeek);
    return endOfDay(startDay.add(const Duration(days: 6)));
  }

  /// Returns the start of the month for the given date
  /// 
  /// Example:
  /// ```dart
  /// DateTime startOfMonth = DateTimeUtils.startOfMonth(DateTime.now());
  /// ```
  static DateTime startOfMonth(DateTime date) {
    return DateTime(date.year, date.month, 1);
  }

  /// Returns the end of the month for the given date
  /// 
  /// Example:
  /// ```dart
  /// DateTime endOfMonth = DateTimeUtils.endOfMonth(DateTime.now());
  /// ```
  static DateTime endOfMonth(DateTime date) {
    return endOfDay(DateTime(date.year, date.month + 1, 0));
  }

  /// Returns the start of the year for the given date
  /// 
  /// Example:
  /// ```dart
  /// DateTime startOfYear = DateTimeUtils.startOfYear(DateTime.now());
  /// ```
  static DateTime startOfYear(DateTime date) {
    return DateTime(date.year, 1, 1);
  }

  /// Returns the end of the year for the given date
  /// 
  /// Example:
  /// ```dart
  /// DateTime endOfYear = DateTimeUtils.endOfYear(DateTime.now());
  /// ```
  static DateTime endOfYear(DateTime date) {
    return endOfDay(DateTime(date.year, 12, 31));
  }

  /// Checks if two dates are on the same day
  /// 
  /// Example:
  /// ```dart
  /// bool isSameDay = DateTimeUtils.isSameDay(
  ///   DateTime.now(),
  ///   DateTime.now().add(Duration(hours: 2)),
  /// );
  /// ```
  static bool isSameDay(DateTime date1, DateTime date2) {
    return date1.year == date2.year && date1.month == date2.month && date1.day == date2.day;
  }

  /// Checks if a date is today
  /// 
  /// Example:
  /// ```dart
  /// bool isToday = DateTimeUtils.isToday(DateTime.now());
  /// ```
  static bool isToday(DateTime date) {
    return isSameDay(date, DateTime.now());
  }

  /// Checks if a date is yesterday
  /// 
  /// Example:
  /// ```dart
  /// bool isYesterday = DateTimeUtils.isYesterday(
  ///   DateTime.now().subtract(Duration(days: 1)),
  /// );
  /// ```
  static bool isYesterday(DateTime date) {
    return isSameDay(date, DateTime.now().subtract(const Duration(days: 1)));
  }

  /// Checks if a date is tomorrow
  /// 
  /// Example:
  /// ```dart
  /// bool isTomorrow = DateTimeUtils.isTomorrow(
  ///   DateTime.now().add(Duration(days: 1)),
  /// );
  /// ```
  static bool isTomorrow(DateTime date) {
    return isSameDay(date, DateTime.now().add(const Duration(days: 1)));
  }
}
