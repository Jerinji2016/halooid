import 'package:flutter/material.dart';

class CommentList extends StatefulWidget {
  final String taskId;

  const CommentList({super.key, required this.taskId});

  @override
  State<CommentList> createState() => _CommentListState();
}

class _CommentListState extends State<CommentList> {
  final TextEditingController _commentController = TextEditingController();
  final List<Comment> _comments = [];
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    _loadComments();
  }

  @override
  void dispose() {
    _commentController.dispose();
    super.dispose();
  }

  Future<void> _loadComments() async {
    setState(() {
      _isLoading = true;
    });

    // This is a placeholder for actual API call
    // In a real implementation, you would fetch comments from the API
    await Future.delayed(const Duration(milliseconds: 500));
    
    setState(() {
      // Sample comments for UI demonstration
      _comments.addAll([
        Comment(
          id: '1',
          taskId: widget.taskId,
          text: 'This task is progressing well.',
          authorName: 'John Doe',
          createdAt: DateTime.now().subtract(const Duration(days: 2)),
        ),
        Comment(
          id: '2',
          taskId: widget.taskId,
          text: 'I\'ve added some additional requirements to the description.',
          authorName: 'Jane Smith',
          createdAt: DateTime.now().subtract(const Duration(hours: 5)),
        ),
      ]);
      _isLoading = false;
    });
  }

  Future<void> _addComment() async {
    if (_commentController.text.trim().isEmpty) return;
    
    final newComment = Comment(
      id: DateTime.now().millisecondsSinceEpoch.toString(),
      taskId: widget.taskId,
      text: _commentController.text.trim(),
      authorName: 'You', // In a real app, this would be the current user's name
      createdAt: DateTime.now(),
    );
    
    setState(() {
      _comments.insert(0, newComment);
      _commentController.clear();
    });
    
    // In a real implementation, you would send the comment to the API
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        // Comment input
        Row(
          children: [
            Expanded(
              child: TextField(
                controller: _commentController,
                decoration: const InputDecoration(
                  hintText: 'Add a comment...',
                  border: OutlineInputBorder(),
                ),
                maxLines: 3,
                minLines: 1,
              ),
            ),
            const SizedBox(width: 8),
            IconButton(
              icon: const Icon(Icons.send),
              onPressed: _addComment,
              color: Theme.of(context).colorScheme.primary,
            ),
          ],
        ),
        const SizedBox(height: 16),
        
        // Comments list
        if (_isLoading)
          const Center(child: CircularProgressIndicator())
        else if (_comments.isEmpty)
          Center(
            child: Text(
              'No comments yet',
              style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                color: Colors.grey,
              ),
            ),
          )
        else
          ListView.separated(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            itemCount: _comments.length,
            separatorBuilder: (context, index) => const Divider(),
            itemBuilder: (context, index) {
              final comment = _comments[index];
              return _CommentItem(comment: comment);
            },
          ),
      ],
    );
  }
}

class _CommentItem extends StatelessWidget {
  final Comment comment;

  const _CommentItem({required this.comment});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              CircleAvatar(
                radius: 16,
                backgroundColor: Theme.of(context).colorScheme.primary,
                child: Text(
                  comment.authorName.substring(0, 1).toUpperCase(),
                  style: const TextStyle(
                    color: Colors.white,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ),
              const SizedBox(width: 8),
              Text(
                comment.authorName,
                style: Theme.of(context).textTheme.titleSmall?.copyWith(
                  fontWeight: FontWeight.bold,
                ),
              ),
              const Spacer(),
              Text(
                _formatDate(comment.createdAt),
                style: Theme.of(context).textTheme.bodySmall,
              ),
            ],
          ),
          const SizedBox(height: 8),
          Text(
            comment.text,
            style: Theme.of(context).textTheme.bodyMedium,
          ),
        ],
      ),
    );
  }

  String _formatDate(DateTime date) {
    final now = DateTime.now();
    final difference = now.difference(date);
    
    if (difference.inDays > 0) {
      return '${difference.inDays}d ago';
    } else if (difference.inHours > 0) {
      return '${difference.inHours}h ago';
    } else if (difference.inMinutes > 0) {
      return '${difference.inMinutes}m ago';
    } else {
      return 'Just now';
    }
  }
}

class Comment {
  final String id;
  final String taskId;
  final String text;
  final String authorName;
  final DateTime createdAt;

  Comment({
    required this.id,
    required this.taskId,
    required this.text,
    required this.authorName,
    required this.createdAt,
  });
}
