import 'package:flutter/material.dart';
import 'package:taskake/models/task.dart';
import 'package:taskake/theme/app_theme.dart';

class StatusBadge extends StatelessWidget {
  final TaskStatus status;

  const StatusBadge({super.key, required this.status});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: _getStatusColor().withOpacity(0.2),
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: _getStatusColor()),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(
            _getStatusIcon(),
            size: 16,
            color: _getStatusColor(),
          ),
          const SizedBox(width: 4),
          Text(
            _getStatusText(),
            style: TextStyle(
              color: _getStatusColor(),
              fontWeight: FontWeight.bold,
            ),
          ),
        ],
      ),
    );
  }

  Color _getStatusColor() {
    switch (status) {
      case TaskStatus.todo:
        return AppTheme.todoStatusColor;
      case TaskStatus.inProgress:
        return AppTheme.inProgressStatusColor;
      case TaskStatus.completed:
        return AppTheme.completedStatusColor;
      case TaskStatus.blocked:
        return AppTheme.blockedStatusColor;
    }
  }

  String _getStatusText() {
    switch (status) {
      case TaskStatus.todo:
        return 'To Do';
      case TaskStatus.inProgress:
        return 'In Progress';
      case TaskStatus.completed:
        return 'Completed';
      case TaskStatus.blocked:
        return 'Blocked';
    }
  }

  IconData _getStatusIcon() {
    switch (status) {
      case TaskStatus.todo:
        return Icons.check_box_outline_blank;
      case TaskStatus.inProgress:
        return Icons.play_circle_outline;
      case TaskStatus.completed:
        return Icons.check_circle_outline;
      case TaskStatus.blocked:
        return Icons.block;
    }
  }
}
