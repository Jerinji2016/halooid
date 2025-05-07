import 'package:json_annotation/json_annotation.dart';

part 'task.g.dart';

enum TaskStatus {
  @JsonValue('todo')
  todo,
  @JsonValue('in_progress')
  inProgress,
  @JsonValue('completed')
  completed,
  @JsonValue('blocked')
  blocked,
}

enum TaskPriority {
  @JsonValue('low')
  low,
  @JsonValue('medium')
  medium,
  @JsonValue('high')
  high,
  @JsonValue('urgent')
  urgent,
}

@JsonSerializable()
class Task {
  final String id;
  final String title;
  final String? description;
  final DateTime createdAt;
  final DateTime? updatedAt;
  final DateTime? dueDate;
  final TaskStatus status;
  final TaskPriority priority;
  final String? assigneeId;
  final String? assigneeName;
  final String? projectId;
  final String? projectName;
  final List<String>? tags;
  final int? estimatedMinutes;
  final int? actualMinutes;

  Task({
    required this.id,
    required this.title,
    this.description,
    required this.createdAt,
    this.updatedAt,
    this.dueDate,
    required this.status,
    required this.priority,
    this.assigneeId,
    this.assigneeName,
    this.projectId,
    this.projectName,
    this.tags,
    this.estimatedMinutes,
    this.actualMinutes,
  });

  factory Task.fromJson(Map<String, dynamic> json) => _$TaskFromJson(json);

  Map<String, dynamic> toJson() => _$TaskToJson(this);

  Task copyWith({
    String? id,
    String? title,
    String? description,
    DateTime? createdAt,
    DateTime? updatedAt,
    DateTime? dueDate,
    TaskStatus? status,
    TaskPriority? priority,
    String? assigneeId,
    String? assigneeName,
    String? projectId,
    String? projectName,
    List<String>? tags,
    int? estimatedMinutes,
    int? actualMinutes,
  }) {
    return Task(
      id: id ?? this.id,
      title: title ?? this.title,
      description: description ?? this.description,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
      dueDate: dueDate ?? this.dueDate,
      status: status ?? this.status,
      priority: priority ?? this.priority,
      assigneeId: assigneeId ?? this.assigneeId,
      assigneeName: assigneeName ?? this.assigneeName,
      projectId: projectId ?? this.projectId,
      projectName: projectName ?? this.projectName,
      tags: tags ?? this.tags,
      estimatedMinutes: estimatedMinutes ?? this.estimatedMinutes,
      actualMinutes: actualMinutes ?? this.actualMinutes,
    );
  }
}
