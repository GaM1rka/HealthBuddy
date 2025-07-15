class Post {
  final String postId;
  final String userId;
  final String authorName;
  final String title;
  final String content;
  final DateTime createdAt;

  Post({
    required this.postId,
    required this.userId,
    required this.authorName,
    required this.title,
    required this.content,
    required this.createdAt,
  });

  factory Post.fromJson(Map<String, dynamic> json) {
    return Post(
      postId: json['post_id'],
      userId: json['user_id'],
      authorName: json['authorName'] ?? 'Unknown Author', // Set default value
      title: json['title'],
      content: json['content'],
      createdAt: DateTime.parse(json['created_at']),
    );
  }
}
