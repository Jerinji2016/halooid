import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:taskake/models/task.dart';
import 'package:taskake/services/task_service.dart';
import 'package:taskake/theme/app_theme.dart';
import 'package:taskake/utils/date_utils.dart';
import 'package:taskake/widgets/status_badge.dart';
import 'package:taskake/widgets/priority_badge.dart';
import 'package:taskake/widgets/task_timer.dart';
import 'package:taskake/widgets/comment_list.dart';

class TaskDetailScreen extends StatelessWidget {
  final Task task;

  const TaskDetailScreen({super.key, required this.task});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Task Details'),
        actions: [
          IconButton(
            icon: const Icon(Icons.edit),
            onPressed: () {
              Navigator.pushNamed(
                context,
                '/task/edit',
                arguments: task,
              );
            },
            tooltip: 'Edit Task',
          ),
          IconButton(
            icon: const Icon(Icons.delete),
            onPressed: () {
              _showDeleteConfirmation(context);
            },
            tooltip: 'Delete Task',
          ),
        ],
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Title
            Text(
              task.title,
              style: Theme.of(context).textTheme.headlineMedium,
            ),
            const SizedBox(height: 16),
            
            // Status and Priority
            Row(
              children: [
                StatusBadge(status: task.status),
                const SizedBox(width: 8),
                PriorityBadge(priority: task.priority),
              ],
            ),
            const SizedBox(height: 16),
            
            // Due Date
            if (task.dueDate != null) ...[
              _buildInfoRow(
                context,
                'Due Date',
                DateUtils.formatDate(task.dueDate!),
                Icons.calendar_today,
                _getDueDateColor(task.dueDate!),
              ),
              const SizedBox(height: 12),
            ],
            
            // Assignee
            if (task.assigneeName != null) ...[
              _buildInfoRow(
                context,
                'Assigned to',
                task.assigneeName!,
                Icons.person,
                null,
              ),
              const SizedBox(height: 12),
            ],
            
            // Project
            if (task.projectName != null) ...[
              _buildInfoRow(
                context,
                'Project',
                task.projectName!,
                Icons.folder,
                null,
              ),
              const SizedBox(height: 12),
            ],
            
            // Time Tracking
            if (task.estimatedMinutes != null || task.actualMinutes != null) ...[
              _buildTimeTracking(context),
              const SizedBox(height: 12),
            ],
            
            // Tags
            if (task.tags != null && task.tags!.isNotEmpty) ...[
              Text(
                'Tags',
                style: Theme.of(context).textTheme.titleMedium,
              ),
              const SizedBox(height: 8),
              Wrap(
                spacing: 8,
                runSpacing: 8,
                children: task.tags!.map((tag) => Chip(
                  label: Text(tag),
                  backgroundColor: Theme.of(context).colorScheme.surface,
                )).toList(),
              ),
              const SizedBox(height: 16),
            ],
            
            // Description
            Text(
              'Description',
              style: Theme.of(context).textTheme.titleMedium,
            ),
            const SizedBox(height: 8),
            Container(
              width: double.infinity,
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                color: Theme.of(context).colorScheme.surface,
                borderRadius: BorderRadius.circular(8),
              ),
              child: Text(
                task.description ?? 'No description provided',
                style: Theme.of(context).textTheme.bodyMedium,
              ),
            ),
            const SizedBox(height: 24),
            
            // Time Tracking Widget
            if (task.status == TaskStatus.inProgress) ...[
              TaskTimer(taskId: task.id),
              const SizedBox(height: 24),
            ],
            
            // Comments Section
            Text(
              'Comments',
              style: Theme.of(context).textTheme.titleMedium,
            ),
            const SizedBox(height: 8),
            CommentList(taskId: task.id),
            
            // Created/Updated Info
            const SizedBox(height: 24),
            Text(
              'Created: ${DateUtils.formatDateTime(task.createdAt)}',
              style: Theme.of(context).textTheme.bodySmall,
            ),
            if (task.updatedAt != null)
              Text(
                'Last updated: ${DateUtils.formatDateTime(task.updatedAt!)}',
                style: Theme.of(context).textTheme.bodySmall,
              ),
          ],
        ),
      ),
      bottomNavigationBar: BottomAppBar(
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 8.0),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceEvenly,
            children: [
              _buildStatusButton(context, TaskStatus.todo),
              _buildStatusButton(context, TaskStatus.inProgress),
              _buildStatusButton(context, TaskStatus.completed),
              _buildStatusButton(context, TaskStatus.blocked),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildInfoRow(BuildContext context, String label, String value, IconData icon, Color? color) {
    return Row(
      children: [
        Icon(
          icon,
          size: 20,
          color: color ?? Theme.of(context).colorScheme.primary,
        ),
        const SizedBox(width: 8),
        Text(
          '$label: ',
          style: Theme.of(context).textTheme.titleSmall,
        ),
        Text(
          value,
          style: Theme.of(context).textTheme.bodyMedium?.copyWith(
            color: color,
          ),
        ),
      ],
    );
  }

  Widget _buildTimeTracking(BuildContext context) {
    final estimatedMinutes = task.estimatedMinutes ?? 0;
    final actualMinutes = task.actualMinutes ?? 0;
    final progress = estimatedMinutes > 0 ? (actualMinutes / estimatedMinutes).clamp(0.0, 1.0) : 0.0;
    
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Time Tracking',
          style: Theme.of(context).textTheme.titleMedium,
        ),
        const SizedBox(height: 8),
        Row(
          children: [
            const Icon(Icons.access_time, size: 16),
            const SizedBox(width: 8),
            Expanded(
              child: LinearProgressIndicator(
                value: progress,
                backgroundColor: Theme.of(context).colorScheme.surface,
              ),
            ),
            const SizedBox(width: 8),
            Text(
              '${actualMinutes ~/ 60}h ${actualMinutes % 60}m / ${estimatedMinutes ~/ 60}h ${estimatedMinutes % 60}m',
              style: Theme.of(context).textTheme.bodySmall,
            ),
          ],
        ),
      ],
    );
  }

  Widget _buildStatusButton(BuildContext context, TaskStatus status) {
    final isCurrentStatus = task.status == status;
    
    Color getStatusColor() {
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

    String getStatusText() {
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

    return ElevatedButton(
      onPressed: isCurrentStatus ? null : () => _updateTaskStatus(context, status),
      style: ElevatedButton.styleFrom(
        backgroundColor: isCurrentStatus ? getStatusColor() : null,
        foregroundColor: isCurrentStatus ? Colors.white : null,
        padding: const EdgeInsets.symmetric(horizontal: 8),
      ),
      child: Text(getStatusText()),
    );
  }

  void _updateTaskStatus(BuildContext context, TaskStatus newStatus) {
    final taskService = Provider.of<TaskService>(context, listen: false);
    final updatedTask = task.copyWith(
      status: newStatus,
      updatedAt: DateTime.now(),
    );
    
    taskService.updateTask(updatedTask).then((result) {
      if (result != null) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Task status updated to ${newStatus.toString().split('.').last}')),
        );
        Navigator.pop(context);
      }
    });
  }

  void _showDeleteConfirmation(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Delete Task'),
        content: const Text('Are you sure you want to delete this task? This action cannot be undone.'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('CANCEL'),
          ),
          TextButton(
            onPressed: () {
              Navigator.pop(context);
              _deleteTask(context);
            },
            child: const Text('DELETE'),
          ),
        ],
      ),
    );
  }

  void _deleteTask(BuildContext context) {
    final taskService = Provider.of<TaskService>(context, listen: false);
    
    taskService.deleteTask(task.id).then((success) {
      if (success) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Task deleted successfully')),
        );
        Navigator.pop(context);
      }
    });
  }

  Color _getDueDateColor(DateTime dueDate) {
    final now = DateTime.now();
    final difference = dueDate.difference(now).inDays;
    
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
