import 'package:flutter/material.dart';
import 'package:taskake/models/task.dart';
import 'package:taskake/theme/app_theme.dart';

class StatusSelector extends StatelessWidget {
  final TaskStatus initialStatus;
  final Function(TaskStatus) onChanged;

  const StatusSelector({
    super.key,
    required this.initialStatus,
    required this.onChanged,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        _buildStatusButton(
          context,
          TaskStatus.todo,
          'To Do',
          AppTheme.todoStatusColor,
          Icons.check_box_outline_blank,
        ),
        const SizedBox(width: 8),
        _buildStatusButton(
          context,
          TaskStatus.inProgress,
          'In Progress',
          AppTheme.inProgressStatusColor,
          Icons.play_circle_outline,
        ),
        const SizedBox(width: 8),
        _buildStatusButton(
          context,
          TaskStatus.completed,
          'Completed',
          AppTheme.completedStatusColor,
          Icons.check_circle_outline,
        ),
        const SizedBox(width: 8),
        _buildStatusButton(
          context,
          TaskStatus.blocked,
          'Blocked',
          AppTheme.blockedStatusColor,
          Icons.block,
        ),
      ],
    );
  }

  Widget _buildStatusButton(
    BuildContext context,
    TaskStatus status,
    String label,
    Color color,
    IconData icon,
  ) {
    final isSelected = initialStatus == status;

    return Expanded(
      child: InkWell(
        onTap: () => onChanged(status),
        borderRadius: BorderRadius.circular(8),
        child: Container(
          padding: const EdgeInsets.symmetric(vertical: 12),
          decoration: BoxDecoration(
            color: isSelected
                ? color.withAlpha(51)
                : Colors.transparent, // 0.2 * 255 = 51
            border: Border.all(
              color: isSelected
                  ? color
                  : Colors.grey.withAlpha(128), // 0.5 * 255 = 128
              width: isSelected ? 2 : 1,
            ),
            borderRadius: BorderRadius.circular(8),
          ),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(
                icon,
                color: isSelected ? color : Colors.grey,
              ),
              const SizedBox(height: 4),
              Text(
                label,
                style: TextStyle(
                  color: isSelected ? color : Colors.grey,
                  fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
                ),
                textAlign: TextAlign.center,
              ),
            ],
          ),
        ),
      ),
    );
  }
}
