import 'package:flutter/material.dart';
import 'package:taskake/models/task.dart';
import 'package:taskake/theme/app_theme.dart';

class PriorityBadge extends StatelessWidget {
  final TaskPriority priority;

  const PriorityBadge({super.key, required this.priority});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: _getPriorityColor().withAlpha(51), // 0.2 * 255 = 51
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: _getPriorityColor()),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(
            _getPriorityIcon(),
            size: 16,
            color: _getPriorityColor(),
          ),
          const SizedBox(width: 4),
          Text(
            _getPriorityText(),
            style: TextStyle(
              color: _getPriorityColor(),
              fontWeight: FontWeight.bold,
            ),
          ),
        ],
      ),
    );
  }

  Color _getPriorityColor() {
    switch (priority) {
      case TaskPriority.low:
        return AppTheme.lowPriorityColor;
      case TaskPriority.medium:
        return AppTheme.mediumPriorityColor;
      case TaskPriority.high:
        return AppTheme.highPriorityColor;
      case TaskPriority.urgent:
        return AppTheme.urgentPriorityColor;
    }
  }

  String _getPriorityText() {
    switch (priority) {
      case TaskPriority.low:
        return 'Low';
      case TaskPriority.medium:
        return 'Medium';
      case TaskPriority.high:
        return 'High';
      case TaskPriority.urgent:
        return 'Urgent';
    }
  }

  IconData _getPriorityIcon() {
    switch (priority) {
      case TaskPriority.low:
        return Icons.arrow_downward;
      case TaskPriority.medium:
        return Icons.remove;
      case TaskPriority.high:
        return Icons.arrow_upward;
      case TaskPriority.urgent:
        return Icons.priority_high;
    }
  }
}
