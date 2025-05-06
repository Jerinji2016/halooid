import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import '../theme/app_colors.dart';

/// A customizable text field widget
class AppTextField extends StatelessWidget {
  /// Creates an AppTextField
  const AppTextField({
    Key? key,
    this.controller,
    this.focusNode,
    this.decoration,
    this.keyboardType,
    this.textInputAction,
    this.textCapitalization = TextCapitalization.none,
    this.style,
    this.strutStyle,
    this.textAlign = TextAlign.start,
    this.textAlignVertical,
    this.textDirection,
    this.readOnly = false,
    this.showCursor,
    this.autofocus = false,
    this.obscureText = false,
    this.autocorrect = true,
    this.enableSuggestions = true,
    this.maxLines = 1,
    this.minLines,
    this.expands = false,
    this.maxLength,
    this.maxLengthEnforcement,
    this.onChanged,
    this.onEditingComplete,
    this.onSubmitted,
    this.onTap,
    this.inputFormatters,
    this.enabled,
    this.cursorWidth = 2.0,
    this.cursorHeight,
    this.cursorRadius,
    this.cursorColor,
    this.keyboardAppearance,
    this.scrollPadding = const EdgeInsets.all(20.0),
    this.enableInteractiveSelection = true,
    this.selectionControls,
    this.buildCounter,
    this.scrollController,
    this.scrollPhysics,
    this.autofillHints,
    this.restorationId,
    this.label,
    this.hint,
    this.helperText,
    this.errorText,
    this.prefixIcon,
    this.suffixIcon,
    this.borderRadius = 8.0,
    this.filled = true,
    this.fillColor,
    this.borderColor,
    this.focusedBorderColor,
    this.errorBorderColor,
    this.contentPadding,
    this.dense = false,
  }) : super(key: key);

  /// Text editing controller
  final TextEditingController? controller;
  
  /// Focus node
  final FocusNode? focusNode;
  
  /// Input decoration
  final InputDecoration? decoration;
  
  /// Keyboard type
  final TextInputType? keyboardType;
  
  /// Text input action
  final TextInputAction? textInputAction;
  
  /// Text capitalization
  final TextCapitalization textCapitalization;
  
  /// Text style
  final TextStyle? style;
  
  /// Strut style
  final StrutStyle? strutStyle;
  
  /// Text alignment
  final TextAlign textAlign;
  
  /// Text alignment vertical
  final TextAlignVertical? textAlignVertical;
  
  /// Text direction
  final TextDirection? textDirection;
  
  /// Whether the text field is read-only
  final bool readOnly;
  
  /// Whether to show cursor
  final bool? showCursor;
  
  /// Whether to autofocus
  final bool autofocus;
  
  /// Whether to obscure text
  final bool obscureText;
  
  /// Whether to enable autocorrect
  final bool autocorrect;
  
  /// Whether to enable suggestions
  final bool enableSuggestions;
  
  /// Maximum number of lines
  final int? maxLines;
  
  /// Minimum number of lines
  final int? minLines;
  
  /// Whether the text field expands to fill its parent
  final bool expands;
  
  /// Maximum length of input
  final int? maxLength;
  
  /// Maximum length enforcement strategy
  final MaxLengthEnforcement? maxLengthEnforcement;
  
  /// Callback when text changes
  final ValueChanged<String>? onChanged;
  
  /// Callback when editing is complete
  final VoidCallback? onEditingComplete;
  
  /// Callback when text is submitted
  final ValueChanged<String>? onSubmitted;
  
  /// Callback when text field is tapped
  final GestureTapCallback? onTap;
  
  /// Input formatters
  final List<TextInputFormatter>? inputFormatters;
  
  /// Whether the text field is enabled
  final bool? enabled;
  
  /// Cursor width
  final double cursorWidth;
  
  /// Cursor height
  final double? cursorHeight;
  
  /// Cursor radius
  final Radius? cursorRadius;
  
  /// Cursor color
  final Color? cursorColor;
  
  /// Keyboard appearance
  final Brightness? keyboardAppearance;
  
  /// Scroll padding
  final EdgeInsets scrollPadding;
  
  /// Whether to enable interactive selection
  final bool enableInteractiveSelection;
  
  /// Selection controls
  final TextSelectionControls? selectionControls;
  
  /// Counter builder
  final InputCounterWidgetBuilder? buildCounter;
  
  /// Scroll controller
  final ScrollController? scrollController;
  
  /// Scroll physics
  final ScrollPhysics? scrollPhysics;
  
  /// Autofill hints
  final Iterable<String>? autofillHints;
  
  /// Restoration ID
  final String? restorationId;
  
  /// Field label
  final String? label;
  
  /// Field hint
  final String? hint;
  
  /// Helper text
  final String? helperText;
  
  /// Error text
  final String? errorText;
  
  /// Prefix icon
  final Widget? prefixIcon;
  
  /// Suffix icon
  final Widget? suffixIcon;
  
  /// Border radius
  final double borderRadius;
  
  /// Whether the field is filled
  final bool filled;
  
  /// Fill color
  final Color? fillColor;
  
  /// Border color
  final Color? borderColor;
  
  /// Focused border color
  final Color? focusedBorderColor;
  
  /// Error border color
  final Color? errorBorderColor;
  
  /// Content padding
  final EdgeInsetsGeometry? contentPadding;
  
  /// Whether the field is dense
  final bool dense;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    
    // Create input decoration
    final InputDecoration effectiveDecoration = decoration ??
        InputDecoration(
          labelText: label,
          hintText: hint,
          helperText: helperText,
          errorText: errorText,
          prefixIcon: prefixIcon,
          suffixIcon: suffixIcon,
          filled: filled,
          fillColor: fillColor ?? theme.inputDecorationTheme.fillColor,
          contentPadding: contentPadding ?? (dense
              ? const EdgeInsets.symmetric(horizontal: 12, vertical: 8)
              : const EdgeInsets.symmetric(horizontal: 16, vertical: 12)),
          isDense: dense,
          border: OutlineInputBorder(
            borderRadius: BorderRadius.circular(borderRadius),
            borderSide: BorderSide(
              color: borderColor ?? AppColors.neutralMain,
            ),
          ),
          enabledBorder: OutlineInputBorder(
            borderRadius: BorderRadius.circular(borderRadius),
            borderSide: BorderSide(
              color: borderColor ?? AppColors.neutralMain,
            ),
          ),
          focusedBorder: OutlineInputBorder(
            borderRadius: BorderRadius.circular(borderRadius),
            borderSide: BorderSide(
              color: focusedBorderColor ?? AppColors.primaryMain,
              width: 2,
            ),
          ),
          errorBorder: OutlineInputBorder(
            borderRadius: BorderRadius.circular(borderRadius),
            borderSide: BorderSide(
              color: errorBorderColor ?? AppColors.errorMain,
            ),
          ),
          focusedErrorBorder: OutlineInputBorder(
            borderRadius: BorderRadius.circular(borderRadius),
            borderSide: BorderSide(
              color: errorBorderColor ?? AppColors.errorMain,
              width: 2,
            ),
          ),
        );
    
    return TextField(
      controller: controller,
      focusNode: focusNode,
      decoration: effectiveDecoration,
      keyboardType: keyboardType,
      textInputAction: textInputAction,
      textCapitalization: textCapitalization,
      style: style,
      strutStyle: strutStyle,
      textAlign: textAlign,
      textAlignVertical: textAlignVertical,
      textDirection: textDirection,
      readOnly: readOnly,
      showCursor: showCursor,
      autofocus: autofocus,
      obscureText: obscureText,
      autocorrect: autocorrect,
      enableSuggestions: enableSuggestions,
      maxLines: maxLines,
      minLines: minLines,
      expands: expands,
      maxLength: maxLength,
      maxLengthEnforcement: maxLengthEnforcement,
      onChanged: onChanged,
      onEditingComplete: onEditingComplete,
      onSubmitted: onSubmitted,
      onTap: onTap,
      inputFormatters: inputFormatters,
      enabled: enabled,
      cursorWidth: cursorWidth,
      cursorHeight: cursorHeight,
      cursorRadius: cursorRadius,
      cursorColor: cursorColor ?? AppColors.primaryMain,
      keyboardAppearance: keyboardAppearance,
      scrollPadding: scrollPadding,
      enableInteractiveSelection: enableInteractiveSelection,
      selectionControls: selectionControls,
      buildCounter: buildCounter,
      scrollController: scrollController,
      scrollPhysics: scrollPhysics,
      autofillHints: autofillHints,
      restorationId: restorationId,
    );
  }
}
