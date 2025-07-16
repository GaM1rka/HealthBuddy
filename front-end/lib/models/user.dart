import 'package:health_buddy_app/models/post.dart';

class UserProfile {
  final String userId;
  final String name;
  final String bio;
  final String avatarUrl;
  final DateTime createdAt;
  final List<Post> posts;

  UserProfile({
    required this.userId,
    required this.name,
    required this.bio,
    required this.avatarUrl,
    required this.createdAt,
    required this.posts,
  });

  factory UserProfile.fromJson(Map<String, dynamic> json) {
    var postList = json['posts'] as List? ?? [];
    List<Post> posts = postList.map((i) => Post.fromJson(i)).toList();

    return UserProfile(
      userId: json['user_id'] ?? '',
      name: json['name'] ?? '',
      bio: json['bio'] ?? '',
      avatarUrl: json['avatar'] ?? 'https://i.pravatar.cc/150?u=${json['user_id']}',
      createdAt: DateTime.parse(json['created_at']),
      posts: posts,
    );
  }
}

class User {
  final String id;
  final String username;
  final String email;
  final DateTime createdAt;

  User({
    required this.id,
    required this.username,
    required this.email,
    required this.createdAt,
  });

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'],
      username: json['username'],
      email: json['email'],
      createdAt: DateTime.parse(json['created_at']),
    );
  }
}

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
