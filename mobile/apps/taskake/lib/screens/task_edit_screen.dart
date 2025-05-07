import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:taskake/models/task.dart';
import 'package:taskake/services/task_service.dart';
import 'package:taskake/widgets/priority_selector.dart';
import 'package:taskake/widgets/date_time_picker.dart';
import 'package:taskake/widgets/tag_input.dart';
import 'package:taskake/widgets/status_selector.dart';

class TaskEditScreen extends StatefulWidget {
  final Task task;

  const TaskEditScreen({super.key, required this.task});

  @override
  State<TaskEditScreen> createState() => _TaskEditScreenState();
}

class _TaskEditScreenState extends State<TaskEditScreen> {
  final _formKey = GlobalKey<FormState>();
  late TextEditingController _titleController;
  late TextEditingController _descriptionController;
  late TaskStatus _status;
  late TaskPriority _priority;
  late DateTime? _dueDate;
  late List<String> _tags;
  late int? _estimatedHours;
  late int? _estimatedMinutes;
  late String? _assigneeId;
  late String? _assigneeName;
  late String? _projectId;
  late String? _projectName;
  bool _isSubmitting = false;

  @override
  void initState() {
    super.initState();
    _titleController = TextEditingController(text: widget.task.title);
    _descriptionController = TextEditingController(text: widget.task.description);
    _status = widget.task.status;
    _priority = widget.task.priority;
    _dueDate = widget.task.dueDate;
    _tags = widget.task.tags ?? [];
    
    // Split estimated minutes into hours and minutes
    if (widget.task.estimatedMinutes != null) {
      _estimatedHours = widget.task.estimatedMinutes! ~/ 60;
      _estimatedMinutes = widget.task.estimatedMinutes! % 60;
    } else {
      _estimatedHours = null;
      _estimatedMinutes = null;
    }
    
    _assigneeId = widget.task.assigneeId;
    _assigneeName = widget.task.assigneeName;
    _projectId = widget.task.projectId;
    _projectName = widget.task.projectName;
  }

  @override
  void dispose() {
    _titleController.dispose();
    _descriptionController.dispose();
    super.dispose();
  }

  Future<void> _submitForm() async {
    if (_formKey.currentState!.validate()) {
      setState(() {
        _isSubmitting = true;
      });

      final taskService = Provider.of<TaskService>(context, listen: false);
      
      // Calculate estimated minutes
      int? totalEstimatedMinutes;
      if (_estimatedHours != null || _estimatedMinutes != null) {
        totalEstimatedMinutes = (_estimatedHours ?? 0) * 60 + (_estimatedMinutes ?? 0);
      }
      
      // Create updated task
      final updatedTask = widget.task.copyWith(
        title: _titleController.text,
        description: _descriptionController.text,
        updatedAt: DateTime.now(),
        status: _status,
        priority: _priority,
        dueDate: _dueDate,
        assigneeId: _assigneeId,
        assigneeName: _assigneeName,
        projectId: _projectId,
        projectName: _projectName,
        tags: _tags.isNotEmpty ? _tags : null,
        estimatedMinutes: totalEstimatedMinutes,
      );
      
      final result = await taskService.updateTask(updatedTask);
      
      setState(() {
        _isSubmitting = false;
      });
      
      if (result != null && mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Task updated successfully')),
        );
        Navigator.pop(context);
        // Also pop the detail screen to go back to the list
        Navigator.pop(context);
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Edit Task'),
      ),
      body: Form(
        key: _formKey,
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Title
              TextFormField(
                controller: _titleController,
                decoration: const InputDecoration(
                  labelText: 'Title',
                  hintText: 'Enter task title',
                ),
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return 'Please enter a title';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 16),
              
              // Description
              TextFormField(
                controller: _descriptionController,
                decoration: const InputDecoration(
                  labelText: 'Description',
                  hintText: 'Enter task description',
                  alignLabelWithHint: true,
                ),
                maxLines: 5,
              ),
              const SizedBox(height: 16),
              
              // Status
              Text(
                'Status',
                style: Theme.of(context).textTheme.titleMedium,
              ),
              const SizedBox(height: 8),
              StatusSelector(
                initialStatus: _status,
                onChanged: (status) {
                  setState(() {
                    _status = status;
                  });
                },
              ),
              const SizedBox(height: 16),
              
              // Priority
              Text(
                'Priority',
                style: Theme.of(context).textTheme.titleMedium,
              ),
              const SizedBox(height: 8),
              PrioritySelector(
                initialPriority: _priority,
                onChanged: (priority) {
                  setState(() {
                    _priority = priority;
                  });
                },
              ),
              const SizedBox(height: 16),
              
              // Due Date
              Text(
                'Due Date',
                style: Theme.of(context).textTheme.titleMedium,
              ),
              const SizedBox(height: 8),
              DateTimePicker(
                selectedDate: _dueDate,
                onDateSelected: (date) {
                  setState(() {
                    _dueDate = date;
                  });
                },
              ),
              const SizedBox(height: 16),
              
              // Estimated Time
              Text(
                'Estimated Time',
                style: Theme.of(context).textTheme.titleMedium,
              ),
              const SizedBox(height: 8),
              Row(
                children: [
                  Expanded(
                    child: TextFormField(
                      decoration: const InputDecoration(
                        labelText: 'Hours',
                        suffixText: 'h',
                      ),
                      keyboardType: TextInputType.number,
                      initialValue: _estimatedHours?.toString(),
                      onChanged: (value) {
                        setState(() {
                          _estimatedHours = int.tryParse(value);
                        });
                      },
                    ),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: TextFormField(
                      decoration: const InputDecoration(
                        labelText: 'Minutes',
                        suffixText: 'm',
                      ),
                      keyboardType: TextInputType.number,
                      initialValue: _estimatedMinutes?.toString(),
                      onChanged: (value) {
                        setState(() {
                          _estimatedMinutes = int.tryParse(value);
                        });
                      },
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 16),
              
              // Tags
              Text(
                'Tags',
                style: Theme.of(context).textTheme.titleMedium,
              ),
              const SizedBox(height: 8),
              TagInput(
                initialTags: _tags,
                onTagsChanged: (tags) {
                  setState(() {
                    _tags = tags;
                  });
                },
              ),
              const SizedBox(height: 16),
              
              // Assignee (simplified for now)
              TextFormField(
                decoration: const InputDecoration(
                  labelText: 'Assignee',
                  hintText: 'Enter assignee name',
                ),
                initialValue: _assigneeName,
                onChanged: (value) {
                  setState(() {
                    _assigneeName = value.isNotEmpty ? value : null;
                  });
                },
              ),
              const SizedBox(height: 16),
              
              // Project (simplified for now)
              TextFormField(
                decoration: const InputDecoration(
                  labelText: 'Project',
                  hintText: 'Enter project name',
                ),
                initialValue: _projectName,
                onChanged: (value) {
                  setState(() {
                    _projectName = value.isNotEmpty ? value : null;
                  });
                },
              ),
              const SizedBox(height: 24),
              
              // Submit Button
              SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  onPressed: _isSubmitting ? null : _submitForm,
                  style: ElevatedButton.styleFrom(
                    padding: const EdgeInsets.symmetric(vertical: 16),
                  ),
                  child: _isSubmitting
                      ? const CircularProgressIndicator()
                      : const Text('UPDATE TASK'),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
