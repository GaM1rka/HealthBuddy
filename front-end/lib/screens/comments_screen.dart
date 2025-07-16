import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:health_buddy_app/models/comment.dart';
import 'package:health_buddy_app/services/api_service.dart';
import 'package:intl/intl.dart';

class CommentsScreen extends StatefulWidget {
  final String postId;

  const CommentsScreen({super.key, required this.postId});

  @override
  _CommentsScreenState createState() => _CommentsScreenState();
}

class _CommentsScreenState extends State<CommentsScreen> {
  final ApiService _apiService = ApiService();
  final TextEditingController _commentController = TextEditingController();
  late Future<List<Comment>> _commentsFuture;

  @override
  void initState() {
    super.initState();
    _commentsFuture = _apiService.getCommentsForPost(widget.postId);
  }

  void _postComment() async {
    if (_commentController.text.isEmpty) {
      return;
    }

    try {
      await _apiService.createComment(widget.postId, _commentController.text);
      _commentController.clear();
      setState(() {
        _commentsFuture = _apiService.getCommentsForPost(widget.postId);
      });
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Failed to post comment: $e')),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Comments', style: GoogleFonts.roboto()),
      ),
      body: Column(
        children: [
          Expanded(
            child: FutureBuilder<List<Comment>>(
              future: _commentsFuture,
              builder: (context, snapshot) {
                if (snapshot.connectionState == ConnectionState.waiting) {
                  return const Center(child: CircularProgressIndicator());
                } else if (snapshot.hasError) {
                  return Center(child: Text('Failed to load comments: ${snapshot.error}'));
                } else if (snapshot.hasData && snapshot.data!.isNotEmpty) {
                  final comments = snapshot.data!;
                  return ListView.builder(
                    itemCount: comments.length,
                    itemBuilder: (context, index) {
                      final comment = comments[index];
                      return ListTile(
                        title: Text(comment.content),
                        subtitle: Text('by ${comment.userId} on ${DateFormat.yMMMd().format(comment.createdAt)}'),
                      );
                    },
                  );
                } else {
                  return const Center(child: Text('No comments yet.'));
                }
              },
            ),
          ),
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: Row(
              children: [
                Expanded(
                  child: TextField(
                    controller: _commentController,
                    decoration: const InputDecoration(
                      hintText: 'Write a comment...',
                      border: OutlineInputBorder(),
                    ),
                  ),
                ),
                IconButton(
                  icon: const Icon(Icons.send),
                  onPressed: _postComment,
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
