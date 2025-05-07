import 'package:flutter/material.dart';
import 'package:taskake/models/task.dart';
import 'package:taskake/theme/app_theme.dart';

class PrioritySelector extends StatelessWidget {
  final TaskPriority initialPriority;
  final Function(TaskPriority) onChanged;

  const PrioritySelector({
    super.key,
    required this.initialPriority,
    required this.onChanged,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        _buildPriorityButton(
          context,
          TaskPriority.low,
          'Low',
          AppTheme.lowPriorityColor,
        ),
        const SizedBox(width: 8),
        _buildPriorityButton(
          context,
          TaskPriority.medium,
          'Medium',
          AppTheme.mediumPriorityColor,
        ),
        const SizedBox(width: 8),
        _buildPriorityButton(
          context,
          TaskPriority.high,
          'High',
          AppTheme.highPriorityColor,
        ),
        const SizedBox(width: 8),
        _buildPriorityButton(
          context,
          TaskPriority.urgent,
          'Urgent',
          AppTheme.urgentPriorityColor,
        ),
      ],
    );
  }

  Widget _buildPriorityButton(
    BuildContext context,
    TaskPriority priority,
    String label,
    Color color,
  ) {
    final isSelected = initialPriority == priority;
    
    return Expanded(
      child: InkWell(
        onTap: () => onChanged(priority),
        borderRadius: BorderRadius.circular(8),
        child: Container(
          padding: const EdgeInsets.symmetric(vertical: 12),
          decoration: BoxDecoration(
            color: isSelected ? color.withOpacity(0.2) : Colors.transparent,
            border: Border.all(
              color: isSelected ? color : Colors.grey.withOpacity(0.5),
              width: isSelected ? 2 : 1,
            ),
            borderRadius: BorderRadius.circular(8),
          ),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(
                _getPriorityIcon(priority),
                color: isSelected ? color : Colors.grey,
              ),
              const SizedBox(height: 4),
              Text(
                label,
                style: TextStyle(
                  color: isSelected ? color : Colors.grey,
                  fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  IconData _getPriorityIcon(TaskPriority priority) {
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
