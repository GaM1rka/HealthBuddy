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
      // MODIFIED: Changed 'authorName' to 'author_name' to match the likely API response.
      authorName: json['name'] ?? 'Unknown Author', 
      title: json['title'],
      content: json['content'],
      createdAt: DateTime.parse(json['created_at']),
    );
  }
}