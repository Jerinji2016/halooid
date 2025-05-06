import 'package:flutter/material.dart';
import '../theme/app_colors.dart';

/// A customizable card widget
class AppCard extends StatelessWidget {
  /// Creates an AppCard
  const AppCard({
    Key? key,
    required this.child,
    this.elevation,
    this.color,
    this.shadowColor,
    this.borderRadius = 12.0,
    this.borderColor,
    this.margin,
    this.padding,
    this.clipBehavior = Clip.antiAlias,
    this.onTap,
    this.onLongPress,
    this.hoverable = false,
  }) : super(key: key);

  /// Card content
  final Widget child;
  
  /// Card elevation
  final double? elevation;
  
  /// Card color
  final Color? color;
  
  /// Shadow color
  final Color? shadowColor;
  
  /// Border radius
  final double borderRadius;
  
  /// Border color
  final Color? borderColor;
  
  /// Card margin
  final EdgeInsetsGeometry? margin;
  
  /// Card padding
  final EdgeInsetsGeometry? padding;
  
  /// Clip behavior
  final Clip clipBehavior;
  
  /// Callback when card is tapped
  final VoidCallback? onTap;
  
  /// Callback when card is long pressed
  final VoidCallback? onLongPress;
  
  /// Whether the card is hoverable
  final bool hoverable;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    
    // Create card shape
    final shape = RoundedRectangleBorder(
      borderRadius: BorderRadius.circular(borderRadius),
      side: borderColor != null
          ? BorderSide(color: borderColor!)
          : BorderSide.none,
    );
    
    // Create card content
    Widget cardContent = child;
    
    // Add padding if specified
    if (padding != null) {
      cardContent = Padding(
        padding: padding!,
        child: cardContent,
      );
    }
    
    // Create card
    Widget card = Card(
      elevation: elevation,
      color: color ?? theme.cardTheme.color,
      shadowColor: shadowColor,
      shape: shape,
      margin: margin ?? theme.cardTheme.margin,
      clipBehavior: clipBehavior,
      child: cardContent,
    );
    
    // Add tap behavior if specified
    if (onTap != null || onLongPress != null) {
      card = InkWell(
        onTap: onTap,
        onLongPress: onLongPress,
        borderRadius: BorderRadius.circular(borderRadius),
        child: card,
      );
    }
    
    // Add hover effect if specified
    if (hoverable) {
      card = MouseRegion(
        cursor: onTap != null ? SystemMouseCursors.click : SystemMouseCursors.basic,
        child: AnimatedContainer(
          duration: const Duration(milliseconds: 200),
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(borderRadius),
            boxShadow: [
              BoxShadow(
                color: AppColors.neutralDark.withOpacity(0.1),
                blurRadius: 8,
                offset: const Offset(0, 4),
              ),
            ],
          ),
          child: card,
        ),
      );
    }
    
    return card;
  }
}

/// A card header widget
class AppCardHeader extends StatelessWidget {
  /// Creates an AppCardHeader
  const AppCardHeader({
    Key? key,
    this.title,
    this.subtitle,
    this.leading,
    this.trailing,
    this.padding,
    this.borderBottom = true,
    this.borderColor,
  }) : super(key: key);

  /// Header title
  final Widget? title;
  
  /// Header subtitle
  final Widget? subtitle;
  
  /// Leading widget
  final Widget? leading;
  
  /// Trailing widget
  final Widget? trailing;
  
  /// Header padding
  final EdgeInsetsGeometry? padding;
  
  /// Whether to show a border at the bottom
  final bool borderBottom;
  
  /// Border color
  final Color? borderColor;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    
    return Container(
      padding: padding ?? const EdgeInsets.all(16),
      decoration: borderBottom
          ? BoxDecoration(
              border: Border(
                bottom: BorderSide(
                  color: borderColor ?? AppColors.neutralLight,
                  width: 1,
                ),
              ),
            )
          : null,
      child: Row(
        children: [
          if (leading != null) ...[
            leading!,
            const SizedBox(width: 16),
          ],
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisSize: MainAxisSize.min,
              children: [
                if (title != null)
                  DefaultTextStyle(
                    style: theme.textTheme.titleMedium!,
                    child: title!,
                  ),
                if (title != null && subtitle != null)
                  const SizedBox(height: 4),
                if (subtitle != null)
                  DefaultTextStyle(
                    style: theme.textTheme.bodyMedium!.copyWith(
                      color: AppColors.textSecondary,
                    ),
                    child: subtitle!,
                  ),
              ],
            ),
          ),
          if (trailing != null) ...[
            const SizedBox(width: 16),
            trailing!,
          ],
        ],
      ),
    );
  }
}

/// A card content widget
class AppCardContent extends StatelessWidget {
  /// Creates an AppCardContent
  const AppCardContent({
    Key? key,
    required this.child,
    this.padding,
  }) : super(key: key);

  /// Content child
  final Widget child;
  
  /// Content padding
  final EdgeInsetsGeometry? padding;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: padding ?? const EdgeInsets.all(16),
      child: child,
    );
  }
}

/// A card footer widget
class AppCardFooter extends StatelessWidget {
  /// Creates an AppCardFooter
  const AppCardFooter({
    Key? key,
    required this.child,
    this.padding,
    this.borderTop = true,
    this.borderColor,
    this.alignment = MainAxisAlignment.end,
  }) : super(key: key);

  /// Footer child
  final Widget child;
  
  /// Footer padding
  final EdgeInsetsGeometry? padding;
  
  /// Whether to show a border at the top
  final bool borderTop;
  
  /// Border color
  final Color? borderColor;
  
  /// Footer alignment
  final MainAxisAlignment alignment;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: padding ?? const EdgeInsets.all(16),
      decoration: borderTop
          ? BoxDecoration(
              border: Border(
                top: BorderSide(
                  color: borderColor ?? AppColors.neutralLight,
                  width: 1,
                ),
              ),
            )
          : null,
      child: DefaultTextStyle(
        style: Theme.of(context).textTheme.bodyMedium!,
        child: Row(
          mainAxisAlignment: alignment,
          children: [child],
        ),
      ),
    );
  }
}
