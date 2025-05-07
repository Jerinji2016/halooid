import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:taskake/services/task_service.dart';
import 'package:taskake/screens/task_list_screen.dart';
import 'package:taskake/screens/task_detail_screen.dart';
import 'package:taskake/screens/task_create_screen.dart';
import 'package:taskake/screens/task_edit_screen.dart';
import 'package:taskake/models/task.dart';
import 'package:taskake/theme/app_theme.dart';

void main() {
  runApp(
    MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (_) => TaskService()),
      ],
      child: const TaskakeApp(),
    ),
  );
}

class TaskakeApp extends StatelessWidget {
  const TaskakeApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Taskake',
      theme: AppTheme.lightTheme,
      darkTheme: AppTheme.darkTheme,
      themeMode: ThemeMode.system,
      initialRoute: '/',
      routes: {
        '/': (context) => const TaskListScreen(),
        '/task/create': (context) => const TaskCreateScreen(),
      },
      onGenerateRoute: (settings) {
        if (settings.name == '/task/detail') {
          final Task task = settings.arguments as Task;
          return MaterialPageRoute(
            builder: (context) => TaskDetailScreen(task: task),
          );
        } else if (settings.name == '/task/edit') {
          final Task task = settings.arguments as Task;
          return MaterialPageRoute(
            builder: (context) => TaskEditScreen(task: task),
          );
        }
        return null;
      },
    );
  }
}
