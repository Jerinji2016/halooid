import 'dart:convert';
import 'package:flutter/foundation.dart';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';
import 'package:taskake/models/task.dart';
import 'package:taskake/utils/api_config.dart';
import 'package:taskake/utils/connectivity_helper.dart';

class TaskService extends ChangeNotifier {
  List<Task> _tasks = [];
  bool _isLoading = false;
  String? _error;
  
  List<Task> get tasks => _tasks;
  bool get isLoading => _isLoading;
  String? get error => _error;

  // Fetch all tasks
  Future<void> fetchTasks() async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      // Check if device is online
      final isOnline = await ConnectivityHelper.isOnline();
      
      if (isOnline) {
        // Fetch from API
        final response = await http.get(
          Uri.parse('${ApiConfig.baseUrl}/tasks'),
          headers: await ApiConfig.getHeaders(),
        );

        if (response.statusCode == 200) {
          final List<dynamic> data = json.decode(response.body);
          _tasks = data.map((json) => Task.fromJson(json)).toList();
          
          // Save to local storage
          await _saveTasks(_tasks);
        } else {
          throw Exception('Failed to load tasks: ${response.statusCode}');
        }
      } else {
        // Load from local storage
        _tasks = await _loadTasks();
      }
      
      _isLoading = false;
      notifyListeners();
    } catch (e) {
      _error = e.toString();
      _isLoading = false;
      
      // Try to load from local storage if API call fails
      try {
        _tasks = await _loadTasks();
      } catch (e) {
        // If local storage also fails, keep the error
      }
      
      notifyListeners();
    }
  }

  // Get task by ID
  Future<Task?> getTask(String id) async {
    try {
      final isOnline = await ConnectivityHelper.isOnline();
      
      if (isOnline) {
        final response = await http.get(
          Uri.parse('${ApiConfig.baseUrl}/tasks/$id'),
          headers: await ApiConfig.getHeaders(),
        );

        if (response.statusCode == 200) {
          return Task.fromJson(json.decode(response.body));
        } else {
          throw Exception('Failed to load task: ${response.statusCode}');
        }
      } else {
        // Find in local cache
        final tasks = await _loadTasks();
        return tasks.firstWhere((task) => task.id == id);
      }
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return null;
    }
  }

  // Create a new task
  Future<Task?> createTask(Task task) async {
    _isLoading = true;
    notifyListeners();

    try {
      final isOnline = await ConnectivityHelper.isOnline();
      
      if (isOnline) {
        final response = await http.post(
          Uri.parse('${ApiConfig.baseUrl}/tasks'),
          headers: await ApiConfig.getHeaders(),
          body: json.encode(task.toJson()),
        );

        if (response.statusCode == 201) {
          final newTask = Task.fromJson(json.decode(response.body));
          _tasks.add(newTask);
          await _saveTasks(_tasks);
          _isLoading = false;
          notifyListeners();
          return newTask;
        } else {
          throw Exception('Failed to create task: ${response.statusCode}');
        }
      } else {
        // Store locally and mark for sync later
        final offlineTask = task;
        _tasks.add(offlineTask);
        await _saveTasks(_tasks);
        await _addToSyncQueue('create', offlineTask);
        _isLoading = false;
        notifyListeners();
        return offlineTask;
      }
    } catch (e) {
      _error = e.toString();
      _isLoading = false;
      notifyListeners();
      return null;
    }
  }

  // Update an existing task
  Future<Task?> updateTask(Task task) async {
    _isLoading = true;
    notifyListeners();

    try {
      final isOnline = await ConnectivityHelper.isOnline();
      
      if (isOnline) {
        final response = await http.put(
          Uri.parse('${ApiConfig.baseUrl}/tasks/${task.id}'),
          headers: await ApiConfig.getHeaders(),
          body: json.encode(task.toJson()),
        );

        if (response.statusCode == 200) {
          final updatedTask = Task.fromJson(json.decode(response.body));
          final index = _tasks.indexWhere((t) => t.id == task.id);
          if (index != -1) {
            _tasks[index] = updatedTask;
          }
          await _saveTasks(_tasks);
          _isLoading = false;
          notifyListeners();
          return updatedTask;
        } else {
          throw Exception('Failed to update task: ${response.statusCode}');
        }
      } else {
        // Update locally and mark for sync later
        final index = _tasks.indexWhere((t) => t.id == task.id);
        if (index != -1) {
          _tasks[index] = task;
        }
        await _saveTasks(_tasks);
        await _addToSyncQueue('update', task);
        _isLoading = false;
        notifyListeners();
        return task;
      }
    } catch (e) {
      _error = e.toString();
      _isLoading = false;
      notifyListeners();
      return null;
    }
  }

  // Delete a task
  Future<bool> deleteTask(String id) async {
    _isLoading = true;
    notifyListeners();

    try {
      final isOnline = await ConnectivityHelper.isOnline();
      
      if (isOnline) {
        final response = await http.delete(
          Uri.parse('${ApiConfig.baseUrl}/tasks/$id'),
          headers: await ApiConfig.getHeaders(),
        );

        if (response.statusCode == 204) {
          _tasks.removeWhere((task) => task.id == id);
          await _saveTasks(_tasks);
          _isLoading = false;
          notifyListeners();
          return true;
        } else {
          throw Exception('Failed to delete task: ${response.statusCode}');
        }
      } else {
        // Delete locally and mark for sync later
        final taskToDelete = _tasks.firstWhere((task) => task.id == id);
        _tasks.removeWhere((task) => task.id == id);
        await _saveTasks(_tasks);
        await _addToSyncQueue('delete', taskToDelete);
        _isLoading = false;
        notifyListeners();
        return true;
      }
    } catch (e) {
      _error = e.toString();
      _isLoading = false;
      notifyListeners();
      return false;
    }
  }

  // Sync pending changes when back online
  Future<void> syncPendingChanges() async {
    final prefs = await SharedPreferences.getInstance();
    final syncQueueJson = prefs.getString('task_sync_queue');
    
    if (syncQueueJson == null || syncQueueJson.isEmpty) {
      return;
    }
    
    final List<dynamic> syncQueue = json.decode(syncQueueJson);
    
    for (var item in syncQueue) {
      final action = item['action'];
      final task = Task.fromJson(item['task']);
      
      switch (action) {
        case 'create':
          await createTask(task);
          break;
        case 'update':
          await updateTask(task);
          break;
        case 'delete':
          await deleteTask(task.id);
          break;
      }
    }
    
    // Clear sync queue after processing
    await prefs.setString('task_sync_queue', '[]');
  }

  // Save tasks to local storage
  Future<void> _saveTasks(List<Task> tasks) async {
    final prefs = await SharedPreferences.getInstance();
    final tasksJson = json.encode(tasks.map((task) => task.toJson()).toList());
    await prefs.setString('tasks', tasksJson);
  }

  // Load tasks from local storage
  Future<List<Task>> _loadTasks() async {
    final prefs = await SharedPreferences.getInstance();
    final tasksJson = prefs.getString('tasks');
    
    if (tasksJson == null || tasksJson.isEmpty) {
      return [];
    }
    
    final List<dynamic> tasksList = json.decode(tasksJson);
    return tasksList.map((json) => Task.fromJson(json)).toList();
  }

  // Add an operation to the sync queue
  Future<void> _addToSyncQueue(String action, Task task) async {
    final prefs = await SharedPreferences.getInstance();
    final syncQueueJson = prefs.getString('task_sync_queue') ?? '[]';
    final List<dynamic> syncQueue = json.decode(syncQueueJson);
    
    syncQueue.add({
      'action': action,
      'task': task.toJson(),
      'timestamp': DateTime.now().toIso8601String(),
    });
    
    await prefs.setString('task_sync_queue', json.encode(syncQueue));
  }
}
