import 'package:flutter/material.dart';
import 'package:taskake/models/task.dart';
import 'package:taskake/theme/app_theme.dart';

class TaskFilterDialog extends StatefulWidget {
  final TaskStatus? statusFilter;
  final TaskPriority? priorityFilter;
  final Function(TaskStatus?, TaskPriority?) onApplyFilters;

  const TaskFilterDialog({
    super.key,
    this.statusFilter,
    this.priorityFilter,
    required this.onApplyFilters,
  });

  @override
  State<TaskFilterDialog> createState() => _TaskFilterDialogState();
}

class _TaskFilterDialogState extends State<TaskFilterDialog> {
  late TaskStatus? _selectedStatus;
  late TaskPriority? _selectedPriority;

  @override
  void initState() {
    super.initState();
    _selectedStatus = widget.statusFilter;
    _selectedPriority = widget.priorityFilter;
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: const Text('Filter Tasks'),
      content: SingleChildScrollView(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Status filter
            Text(
              'Status',
              style: Theme.of(context).textTheme.titleMedium,
            ),
            const SizedBox(height: 8),
            Wrap(
              spacing: 8,
              runSpacing: 8,
              children: [
                _buildStatusChip(null, 'All'),
                _buildStatusChip(TaskStatus.todo, 'To Do'),
                _buildStatusChip(TaskStatus.inProgress, 'In Progress'),
                _buildStatusChip(TaskStatus.completed, 'Completed'),
                _buildStatusChip(TaskStatus.blocked, 'Blocked'),
              ],
            ),

            const SizedBox(height: 16),

            // Priority filter
            Text(
              'Priority',
              style: Theme.of(context).textTheme.titleMedium,
            ),
            const SizedBox(height: 8),
            Wrap(
              spacing: 8,
              runSpacing: 8,
              children: [
                _buildPriorityChip(null, 'All'),
                _buildPriorityChip(TaskPriority.low, 'Low'),
                _buildPriorityChip(TaskPriority.medium, 'Medium'),
                _buildPriorityChip(TaskPriority.high, 'High'),
                _buildPriorityChip(TaskPriority.urgent, 'Urgent'),
              ],
            ),
          ],
        ),
      ),
      actions: [
        TextButton(
          onPressed: () {
            Navigator.pop(context);
          },
          child: const Text('CANCEL'),
        ),
        TextButton(
          onPressed: () {
            widget.onApplyFilters(_selectedStatus, _selectedPriority);
            Navigator.pop(context);
          },
          child: const Text('APPLY'),
        ),
      ],
    );
  }

  Widget _buildStatusChip(TaskStatus? status, String label) {
    final isSelected = _selectedStatus == status;

    Color getChipColor() {
      if (status == null) return Colors.grey;

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

    return FilterChip(
      label: Text(label),
      selected: isSelected,
      onSelected: (selected) {
        setState(() {
          _selectedStatus = selected ? status : null;
        });
      },
      backgroundColor: getChipColor().withAlpha(51), // 0.2 * 255 = 51
      selectedColor: getChipColor().withAlpha(153), // 0.6 * 255 = 153
    );
  }

  Widget _buildPriorityChip(TaskPriority? priority, String label) {
    final isSelected = _selectedPriority == priority;

    Color getChipColor() {
      if (priority == null) return Colors.grey;

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

    return FilterChip(
      label: Text(label),
      selected: isSelected,
      onSelected: (selected) {
        setState(() {
          _selectedPriority = selected ? priority : null;
        });
      },
      backgroundColor: getChipColor().withAlpha(51), // 0.2 * 255 = 51
      selectedColor: getChipColor().withAlpha(153), // 0.6 * 255 = 153
    );
  }
}
