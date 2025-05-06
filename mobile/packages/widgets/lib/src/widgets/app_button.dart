import 'package:flutter/material.dart';
import '../theme/app_colors.dart';

/// Button variants
enum AppButtonVariant {
  /// Primary button
  primary,
  
  /// Secondary button
  secondary,
  
  /// Success button
  success,
  
  /// Danger button
  danger,
  
  /// Warning button
  warning,
  
  /// Info button
  info,
}

/// Button sizes
enum AppButtonSize {
  /// Small button
  small,
  
  /// Medium button
  medium,
  
  /// Large button
  large,
}

/// A customizable button widget
class AppButton extends StatelessWidget {
  /// Creates an AppButton
  const AppButton({
    Key? key,
    required this.onPressed,
    required this.child,
    this.variant = AppButtonVariant.primary,
    this.size = AppButtonSize.medium,
    this.outlined = false,
    this.rounded = false,
    this.fullWidth = false,
    this.loading = false,
    this.disabled = false,
    this.icon,
    this.iconPosition = IconPosition.left,
    this.elevation,
  }) : super(key: key);

  /// Button press callback
  final VoidCallback? onPressed;
  
  /// Button content
  final Widget child;
  
  /// Button variant
  final AppButtonVariant variant;
  
  /// Button size
  final AppButtonSize size;
  
  /// Whether the button is outlined
  final bool outlined;
  
  /// Whether the button is rounded
  final bool rounded;
  
  /// Whether the button is full width
  final bool fullWidth;
  
  /// Whether the button is in a loading state
  final bool loading;
  
  /// Whether the button is disabled
  final bool disabled;
  
  /// Optional icon to display
  final Widget? icon;
  
  /// Position of the icon
  final IconPosition iconPosition;
  
  /// Button elevation
  final double? elevation;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isDisabled = disabled || loading || onPressed == null;
    
    // Determine button colors based on variant
    final Color backgroundColor;
    final Color foregroundColor;
    final Color borderColor;
    
    switch (variant) {
      case AppButtonVariant.primary:
        backgroundColor = outlined ? Colors.transparent : AppColors.primaryMain;
        foregroundColor = outlined ? AppColors.primaryMain : AppColors.primaryContrastText;
        borderColor = AppColors.primaryMain;
        break;
      case AppButtonVariant.secondary:
        backgroundColor = outlined ? Colors.transparent : AppColors.secondaryMain;
        foregroundColor = outlined ? AppColors.secondaryMain : AppColors.secondaryContrastText;
        borderColor = AppColors.secondaryMain;
        break;
      case AppButtonVariant.success:
        backgroundColor = outlined ? Colors.transparent : AppColors.successMain;
        foregroundColor = outlined ? AppColors.successMain : AppColors.successContrastText;
        borderColor = AppColors.successMain;
        break;
      case AppButtonVariant.danger:
        backgroundColor = outlined ? Colors.transparent : AppColors.errorMain;
        foregroundColor = outlined ? AppColors.errorMain : AppColors.errorContrastText;
        borderColor = AppColors.errorMain;
        break;
      case AppButtonVariant.warning:
        backgroundColor = outlined ? Colors.transparent : AppColors.warningMain;
        foregroundColor = outlined ? AppColors.warningMain : AppColors.warningContrastText;
        borderColor = AppColors.warningMain;
        break;
      case AppButtonVariant.info:
        backgroundColor = outlined ? Colors.transparent : AppColors.infoMain;
        foregroundColor = outlined ? AppColors.infoMain : AppColors.infoContrastText;
        borderColor = AppColors.infoMain;
        break;
    }
    
    // Determine button padding and text style based on size
    final EdgeInsetsGeometry padding;
    final TextStyle textStyle;
    final double iconSize;
    final double loadingSize;
    
    switch (size) {
      case AppButtonSize.small:
        padding = const EdgeInsets.symmetric(horizontal: 12, vertical: 8);
        textStyle = theme.textTheme.labelMedium!;
        iconSize = 16;
        loadingSize = 16;
        break;
      case AppButtonSize.large:
        padding = const EdgeInsets.symmetric(horizontal: 24, vertical: 16);
        textStyle = theme.textTheme.titleMedium!;
        iconSize = 24;
        loadingSize = 24;
        break;
      case AppButtonSize.medium:
      default:
        padding = const EdgeInsets.symmetric(horizontal: 16, vertical: 12);
        textStyle = theme.textTheme.labelLarge!;
        iconSize = 20;
        loadingSize = 20;
        break;
    }
    
    // Create button style
    final ButtonStyle buttonStyle = ElevatedButton.styleFrom(
      backgroundColor: isDisabled
          ? AppColors.neutralLight
          : backgroundColor,
      foregroundColor: isDisabled
          ? AppColors.textDisabled
          : foregroundColor,
      disabledBackgroundColor: AppColors.neutralLight,
      disabledForegroundColor: AppColors.textDisabled,
      padding: padding,
      elevation: elevation,
      textStyle: textStyle,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(rounded ? 50 : 8),
        side: outlined
            ? BorderSide(color: isDisabled ? AppColors.neutralMain : borderColor)
            : BorderSide.none,
      ),
      minimumSize: fullWidth ? const Size(double.infinity, 0) : null,
    );
    
    // Create button content
    Widget buttonContent;
    
    if (loading) {
      buttonContent = SizedBox(
        width: loadingSize,
        height: loadingSize,
        child: CircularProgressIndicator(
          strokeWidth: 2,
          valueColor: AlwaysStoppedAnimation<Color>(foregroundColor),
        ),
      );
    } else if (icon != null) {
      final iconWidget = IconTheme(
        data: IconThemeData(
          size: iconSize,
          color: foregroundColor,
        ),
        child: icon!,
      );
      
      if (iconPosition == IconPosition.left) {
        buttonContent = Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            iconWidget,
            const SizedBox(width: 8),
            child,
          ],
        );
      } else {
        buttonContent = Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            child,
            const SizedBox(width: 8),
            iconWidget,
          ],
        );
      }
    } else {
      buttonContent = child;
    }
    
    return ElevatedButton(
      onPressed: isDisabled ? null : onPressed,
      style: buttonStyle,
      child: buttonContent,
    );
  }
}

/// Icon position in a button
enum IconPosition {
  /// Icon on the left
  left,
  
  /// Icon on the right
  right,
}
