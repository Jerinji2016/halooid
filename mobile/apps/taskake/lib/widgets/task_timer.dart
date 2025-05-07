import 'dart:async';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:taskake/services/task_service.dart';

class TaskTimer extends StatefulWidget {
  final String taskId;

  const TaskTimer({super.key, required this.taskId});

  @override
  State<TaskTimer> createState() => _TaskTimerState();
}

class _TaskTimerState extends State<TaskTimer> {
  bool _isRunning = false;
  int _elapsedSeconds = 0;
  Timer? _timer;
  final ValueNotifier<String> _timeDisplay = ValueNotifier<String>('00:00:00');

  @override
  void dispose() {
    _timer?.cancel();
    _timeDisplay.dispose();
    super.dispose();
  }

  void _startTimer() {
    setState(() {
      _isRunning = true;
    });
    
    _timer = Timer.periodic(const Duration(seconds: 1), (timer) {
      setState(() {
        _elapsedSeconds++;
        _updateTimeDisplay();
      });
    });
  }

  void _stopTimer() {
    _timer?.cancel();
    setState(() {
      _isRunning = false;
    });
    
    // Save the time to the task
    _saveTimeToTask();
  }

  void _resetTimer() {
    _timer?.cancel();
    setState(() {
      _isRunning = false;
      _elapsedSeconds = 0;
      _updateTimeDisplay();
    });
  }

  void _updateTimeDisplay() {
    final hours = (_elapsedSeconds ~/ 3600).toString().padLeft(2, '0');
    final minutes = ((_elapsedSeconds % 3600) ~/ 60).toString().padLeft(2, '0');
    final seconds = (_elapsedSeconds % 60).toString().padLeft(2, '0');
    _timeDisplay.value = '$hours:$minutes:$seconds';
  }

  Future<void> _saveTimeToTask() async {
    final taskService = Provider.of<TaskService>(context, listen: false);
    final task = await taskService.getTask(widget.taskId);
    
    if (task != null) {
      final currentMinutes = task.actualMinutes ?? 0;
      final additionalMinutes = _elapsedSeconds ~/ 60;
      
      final updatedTask = task.copyWith(
        actualMinutes: currentMinutes + additionalMinutes,
        updatedAt: DateTime.now(),
      );
      
      await taskService.updateTask(updatedTask);
      
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Added ${additionalMinutes}m to task time'),
            duration: const Duration(seconds: 2),
          ),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Time Tracking',
              style: Theme.of(context).textTheme.titleMedium,
            ),
            const SizedBox(height: 16),
            
            // Timer display
            Center(
              child: ValueListenableBuilder<String>(
                valueListenable: _timeDisplay,
                builder: (context, timeString, child) {
                  return Text(
                    timeString,
                    style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                      fontFamily: 'monospace',
                      fontWeight: FontWeight.bold,
                    ),
                  );
                },
              ),
            ),
            const SizedBox(height: 16),
            
            // Timer controls
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceEvenly,
              children: [
                if (!_isRunning)
                  ElevatedButton.icon(
                    onPressed: _startTimer,
                    icon: const Icon(Icons.play_arrow),
                    label: const Text('Start'),
                    style: ElevatedButton.styleFrom(
                      backgroundColor: Colors.green,
                      foregroundColor: Colors.white,
                    ),
                  )
                else
                  ElevatedButton.icon(
                    onPressed: _stopTimer,
                    icon: const Icon(Icons.pause),
                    label: const Text('Stop'),
                    style: ElevatedButton.styleFrom(
                      backgroundColor: Colors.orange,
                      foregroundColor: Colors.white,
                    ),
                  ),
                
                ElevatedButton.icon(
                  onPressed: _resetTimer,
                  icon: const Icon(Icons.refresh),
                  label: const Text('Reset'),
                  style: ElevatedButton.styleFrom(
                    backgroundColor: Colors.red,
                    foregroundColor: Colors.white,
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
