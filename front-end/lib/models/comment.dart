class Comment {
  final String commentId;
  final String userId;
  final String content;
  final DateTime createdAt;

  Comment({
    required this.commentId,
    required this.userId,
    required this.content,
    required this.createdAt,
  });

  factory Comment.fromJson(Map<String, dynamic> json) {
    return Comment(
      commentId: json['comment_id'],
      userId: json['user_id'],
      content: json['content'],
      createdAt: DateTime.parse(json['created_at']),
    );
  }
}
