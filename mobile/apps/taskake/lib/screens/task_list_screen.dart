import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:taskake/models/task.dart';
import 'package:taskake/services/task_service.dart';
import 'package:taskake/widgets/task_list_item.dart';
import 'package:taskake/widgets/task_filter_dialog.dart';
import 'package:taskake/theme/app_theme.dart';

class TaskListScreen extends StatefulWidget {
  const TaskListScreen({super.key});

  @override
  State<TaskListScreen> createState() => _TaskListScreenState();
}

class _TaskListScreenState extends State<TaskListScreen> {
  TaskStatus? _statusFilter;
  TaskPriority? _priorityFilter;
  String? _searchQuery;
  bool _isAscending = true;
  String _sortBy = 'dueDate';

  @override
  void initState() {
    super.initState();
    // Fetch tasks when the screen is first loaded
    WidgetsBinding.instance.addPostFrameCallback((_) {
      Provider.of<TaskService>(context, listen: false).fetchTasks();
    });
  }

  // Filter tasks based on current filters
  List<Task> _getFilteredTasks(List<Task> tasks) {
    return tasks.where((task) {
      // Apply status filter
      if (_statusFilter != null && task.status != _statusFilter) {
        return false;
      }
      
      // Apply priority filter
      if (_priorityFilter != null && task.priority != _priorityFilter) {
        return false;
      }
      
      // Apply search query
      if (_searchQuery != null && _searchQuery!.isNotEmpty) {
        final query = _searchQuery!.toLowerCase();
        return task.title.toLowerCase().contains(query) || 
               (task.description?.toLowerCase().contains(query) ?? false);
      }
      
      return true;
    }).toList();
  }

  // Sort tasks based on current sort criteria
  List<Task> _getSortedTasks(List<Task> tasks) {
    switch (_sortBy) {
      case 'title':
        tasks.sort((a, b) => _isAscending 
          ? a.title.compareTo(b.title) 
          : b.title.compareTo(a.title));
        break;
      case 'priority':
        tasks.sort((a, b) => _isAscending 
          ? a.priority.index.compareTo(b.priority.index) 
          : b.priority.index.compareTo(a.priority.index));
        break;
      case 'status':
        tasks.sort((a, b) => _isAscending 
          ? a.status.index.compareTo(b.status.index) 
          : b.status.index.compareTo(a.status.index));
        break;
      case 'dueDate':
      default:
        tasks.sort((a, b) {
          if (a.dueDate == null && b.dueDate == null) return 0;
          if (a.dueDate == null) return _isAscending ? 1 : -1;
          if (b.dueDate == null) return _isAscending ? -1 : 1;
          return _isAscending 
            ? a.dueDate!.compareTo(b.dueDate!) 
            : b.dueDate!.compareTo(a.dueDate!);
        });
    }
    return tasks;
  }

  void _showFilterDialog() {
    showDialog(
      context: context,
      builder: (context) => TaskFilterDialog(
        statusFilter: _statusFilter,
        priorityFilter: _priorityFilter,
        onApplyFilters: (status, priority) {
          setState(() {
            _statusFilter = status;
            _priorityFilter = priority;
          });
        },
      ),
    );
  }

  void _showSortMenu(BuildContext context) {
    final RenderBox button = context.findRenderObject() as RenderBox;
    final RenderBox overlay = Overlay.of(context).context.findRenderObject() as RenderBox;
    final RelativeRect position = RelativeRect.fromRect(
      Rect.fromPoints(
        button.localToGlobal(Offset.zero, ancestor: overlay),
        button.localToGlobal(button.size.bottomRight(Offset.zero), ancestor: overlay),
      ),
      Offset.zero & overlay.size,
    );

    showMenu<String>(
      context: context,
      position: position,
      items: [
        PopupMenuItem<String>(
          value: 'dueDate',
          child: Text('Due Date ${_sortBy == 'dueDate' ? (_isAscending ? '↑' : '↓') : ''}'),
        ),
        PopupMenuItem<String>(
          value: 'title',
          child: Text('Title ${_sortBy == 'title' ? (_isAscending ? '↑' : '↓') : ''}'),
        ),
        PopupMenuItem<String>(
          value: 'priority',
          child: Text('Priority ${_sortBy == 'priority' ? (_isAscending ? '↑' : '↓') : ''}'),
        ),
        PopupMenuItem<String>(
          value: 'status',
          child: Text('Status ${_sortBy == 'status' ? (_isAscending ? '↑' : '↓') : ''}'),
        ),
      ],
    ).then((value) {
      if (value != null) {
        setState(() {
          if (_sortBy == value) {
            // Toggle direction if same sort criteria
            _isAscending = !_isAscending;
          } else {
            // New sort criteria, default to ascending
            _sortBy = value;
            _isAscending = true;
          }
        });
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Tasks'),
        actions: [
          IconButton(
            icon: const Icon(Icons.filter_list),
            onPressed: _showFilterDialog,
            tooltip: 'Filter Tasks',
          ),
          Builder(
            builder: (context) => IconButton(
              icon: const Icon(Icons.sort),
              onPressed: () => _showSortMenu(context),
              tooltip: 'Sort Tasks',
            ),
          ),
        ],
      ),
      body: Column(
        children: [
          // Search bar
          Padding(
            padding: const EdgeInsets.all(16.0),
            child: TextField(
              decoration: const InputDecoration(
                hintText: 'Search tasks...',
                prefixIcon: Icon(Icons.search),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.all(Radius.circular(25.0)),
                ),
              ),
              onChanged: (value) {
                setState(() {
                  _searchQuery = value;
                });
              },
            ),
          ),
          
          // Filter chips
          if (_statusFilter != null || _priorityFilter != null)
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16.0),
              child: Wrap(
                spacing: 8.0,
                children: [
                  if (_statusFilter != null)
                    Chip(
                      label: Text('Status: ${_statusFilter.toString().split('.').last}'),
                      onDeleted: () {
                        setState(() {
                          _statusFilter = null;
                        });
                      },
                    ),
                  if (_priorityFilter != null)
                    Chip(
                      label: Text('Priority: ${_priorityFilter.toString().split('.').last}'),
                      onDeleted: () {
                        setState(() {
                          _priorityFilter = null;
                        });
                      },
                    ),
                ],
              ),
            ),
          
          // Task list
          Expanded(
            child: Consumer<TaskService>(
              builder: (context, taskService, child) {
                if (taskService.isLoading) {
                  return const Center(child: CircularProgressIndicator());
                }
                
                if (taskService.error != null) {
                  return Center(
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Text(
                          'Error loading tasks',
                          style: Theme.of(context).textTheme.titleLarge,
                        ),
                        const SizedBox(height: 8),
                        Text(
                          taskService.error!,
                          style: Theme.of(context).textTheme.bodyMedium,
                        ),
                        const SizedBox(height: 16),
                        ElevatedButton(
                          onPressed: () => taskService.fetchTasks(),
                          child: const Text('Retry'),
                        ),
                      ],
                    ),
                  );
                }
                
                final filteredTasks = _getFilteredTasks(taskService.tasks);
                final sortedTasks = _getSortedTasks(filteredTasks);
                
                if (sortedTasks.isEmpty) {
                  return Center(
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Icon(
                          Icons.task_alt,
                          size: 64,
                          color: Theme.of(context).disabledColor,
                        ),
                        const SizedBox(height: 16),
                        Text(
                          'No tasks found',
                          style: Theme.of(context).textTheme.titleLarge,
                        ),
                        if (_statusFilter != null || _priorityFilter != null || _searchQuery != null)
                          const SizedBox(height: 8),
                        if (_statusFilter != null || _priorityFilter != null || _searchQuery != null)
                          Text(
                            'Try changing your filters',
                            style: Theme.of(context).textTheme.bodyMedium,
                          ),
                      ],
                    ),
                  );
                }
                
                return RefreshIndicator(
                  onRefresh: () => taskService.fetchTasks(),
                  child: ListView.builder(
                    padding: const EdgeInsets.all(16.0),
                    itemCount: sortedTasks.length,
                    itemBuilder: (context, index) {
                      final task = sortedTasks[index];
                      return TaskListItem(
                        task: task,
                        onTap: () {
                          Navigator.pushNamed(
                            context,
                            '/task/detail',
                            arguments: task,
                          );
                        },
                      );
                    },
                  ),
                );
              },
            ),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          Navigator.pushNamed(context, '/task/create');
        },
        tooltip: 'Add Task',
        child: const Icon(Icons.add),
      ),
    );
  }
}
