class Post {
  final String postId;
  final String userId;
  final String authorName;
  final String content;
  final DateTime createdAt;

  Post({
    required this.postId,
    required this.userId,
    required this.authorName,
    required this.content,
    required this.createdAt,
  });

  factory Post.fromJson(Map<String, dynamic> json) {
    return Post(
      postId: json['id'],
      userId: json['userId'],
      authorName: json['authorName'],
      content: json['content'],
      createdAt: DateTime.parse(json['createdAt']),
    );
  }
}