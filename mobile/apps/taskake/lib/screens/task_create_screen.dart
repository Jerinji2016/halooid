import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:taskake/models/task.dart';
import 'package:taskake/services/task_service.dart';
import 'package:taskake/widgets/priority_selector.dart';
import 'package:taskake/widgets/date_time_picker.dart';
import 'package:taskake/widgets/tag_input.dart';

class TaskCreateScreen extends StatefulWidget {
  const TaskCreateScreen({super.key});

  @override
  State<TaskCreateScreen> createState() => _TaskCreateScreenState();
}

class _TaskCreateScreenState extends State<TaskCreateScreen> {
  final _formKey = GlobalKey<FormState>();
  final _titleController = TextEditingController();
  final _descriptionController = TextEditingController();
  TaskPriority _priority = TaskPriority.medium;
  DateTime? _dueDate;
  List<String> _tags = [];
  int? _estimatedHours;
  int? _estimatedMinutes;
  String? _assigneeId;
  String? _assigneeName;
  String? _projectId;
  String? _projectName;
  bool _isSubmitting = false;

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
      
      // Create a new task
      final newTask = Task(
        id: DateTime.now().millisecondsSinceEpoch.toString(), // Temporary ID, will be replaced by server
        title: _titleController.text,
        description: _descriptionController.text,
        createdAt: DateTime.now(),
        status: TaskStatus.todo,
        priority: _priority,
        dueDate: _dueDate,
        assigneeId: _assigneeId,
        assigneeName: _assigneeName,
        projectId: _projectId,
        projectName: _projectName,
        tags: _tags.isNotEmpty ? _tags : null,
        estimatedMinutes: totalEstimatedMinutes,
      );
      
      final result = await taskService.createTask(newTask);
      
      setState(() {
        _isSubmitting = false;
      });
      
      if (result != null && mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Task created successfully')),
        );
        Navigator.pop(context);
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Create Task'),
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
                      : const Text('CREATE TASK'),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
