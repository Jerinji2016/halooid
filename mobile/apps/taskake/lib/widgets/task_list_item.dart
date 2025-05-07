import 'package:flutter/material.dart';
import 'package:taskake/models/task.dart';
import 'package:taskake/theme/app_theme.dart';
import 'package:taskake/utils/date_utils.dart' as app_date_utils;

class TaskListItem extends StatelessWidget {
  final Task task;
  final VoidCallback onTap;

  const TaskListItem({
    super.key,
    required this.task,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.only(bottom: 12.0),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Title and Priority
              Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Expanded(
                    child: Text(
                      task.title,
                      style: Theme.of(context).textTheme.titleMedium?.copyWith(
                            decoration: task.status == TaskStatus.completed
                                ? TextDecoration.lineThrough
                                : null,
                          ),
                    ),
                  ),
                  const SizedBox(width: 8),
                  _buildPriorityIndicator(),
                ],
              ),

              // Description (if available)
              if (task.description != null && task.description!.isNotEmpty) ...[
                const SizedBox(height: 8),
                Text(
                  task.description!,
                  style: Theme.of(context).textTheme.bodyMedium,
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                ),
              ],

              const SizedBox(height: 12),

              // Bottom row with metadata
              Row(
                children: [
                  // Status indicator
                  Container(
                    padding:
                        const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                    decoration: BoxDecoration(
                      color: _getStatusColor().withAlpha(51), // 0.2 * 255 = 51
                      borderRadius: BorderRadius.circular(4),
                    ),
                    child: Text(
                      _getStatusText(),
                      style: Theme.of(context).textTheme.bodySmall?.copyWith(
                            color: _getStatusColor(),
                          ),
                    ),
                  ),

                  const SizedBox(width: 8),

                  // Due date (if available)
                  if (task.dueDate != null) ...[
                    Icon(
                      Icons.calendar_today,
                      size: 14,
                      color: _getDueDateColor(),
                    ),
                    const SizedBox(width: 4),
                    Text(
                      app_date_utils.DateUtils.formatDate(task.dueDate!),
                      style: Theme.of(context).textTheme.bodySmall?.copyWith(
                            color: _getDueDateColor(),
                          ),
                    ),
                    const SizedBox(width: 8),
                  ],

                  // Assignee (if available)
                  if (task.assigneeName != null) ...[
                    const Icon(
                      Icons.person,
                      size: 14,
                      color: Colors.grey,
                    ),
                    const SizedBox(width: 4),
                    Expanded(
                      child: Text(
                        task.assigneeName!,
                        style: Theme.of(context).textTheme.bodySmall,
                        overflow: TextOverflow.ellipsis,
                      ),
                    ),
                  ],
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildPriorityIndicator() {
    Color color;

    switch (task.priority) {
      case TaskPriority.low:
        color = AppTheme.lowPriorityColor;
        break;
      case TaskPriority.medium:
        color = AppTheme.mediumPriorityColor;
        break;
      case TaskPriority.high:
        color = AppTheme.highPriorityColor;
        break;
      case TaskPriority.urgent:
        color = AppTheme.urgentPriorityColor;
        break;
    }

    return Container(
      width: 16,
      height: 16,
      decoration: BoxDecoration(
        color: color,
        shape: BoxShape.circle,
      ),
    );
  }

  Color _getStatusColor() {
    switch (task.status) {
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
    switch (task.status) {
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

  Color _getDueDateColor() {
    if (task.dueDate == null) return Colors.grey;

    final now = DateTime.now();
    final difference = task.dueDate!.difference(now).inDays;

    if (difference < 0) {
      return Colors.red;
    } else if (difference == 0) {
      return Colors.orange;
    } else if (difference <= 2) {
      return Colors.amber;
    } else {
      return Colors.green;
    }
  }
}
