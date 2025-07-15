class UserProfile {
  final String id;
  final String username;
  final String fullName;
  final String bio;
  final String avatarUrl;

  UserProfile({
    required this.id,
    required this.username,
    required this.fullName,
    required this.bio,
    required this.avatarUrl,
  });

  factory UserProfile.fromJson(Map<String, dynamic> json) {
    return UserProfile(
      id: json['id'] ?? '',
      username: json['username'] ?? '',
      fullName: json['fullName'] ?? '',
      bio: json['bio'] ?? '',
      avatarUrl: json['avatarUrl'] ?? 'https://i.pravatar.cc/150?u=${json['id']}',
    );
  }
}