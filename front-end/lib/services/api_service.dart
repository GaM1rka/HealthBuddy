import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';
import '../models/post.dart';
import '../models/user.dart';

class ApiService {
  final String _baseUrl = 'http://5.159.102.12:8080';

  // --- Token Management ---

  Future<void> _saveToken(String token) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('jwt_token', token);
  }

  Future<String?> getToken() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString('jwt_token');
  }

  Future<void> deleteToken() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove('jwt_token');
  }

  Future<Map<String, String>> _getHeaders() async {
    final token = await getToken();
    if (token != null) {
      return {
        'Content-Type': 'application/json; charset=UTF-8',
        'Authorization': 'Bearer $token',
      };
    }
    return {
      'Content-Type': 'application/json; charset=UTF-8',
    };
  }

  // --- Auth Service ---

  Future<String> login(String username, String password) async {
    final response = await http.post(
      Uri.parse('$_baseUrl/auth/login'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({'username': username, 'password': password}),
    );

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      final token = data['token'];
      await _saveToken(token);
      return token;
    } else {
      throw Exception('Failed to login');
    }
  }

  Future<String> register(String username, String email, String password) async {
    final response = await http.post(
      Uri.parse('$_baseUrl/auth/register'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({
        'username': username,
        'email': email,
        'password': password,
      }),
    );

    if (response.statusCode == 201) {
      final data = jsonDecode(response.body);
      final token = data['token'];
      await _saveToken(token);
      return token;
    } else {
      throw Exception('Failed to register');
    }
  }

  Future<User> getUserById(String id) async {
    final response = await http.get(
      Uri.parse('$_baseUrl/auth/users/$id'),
      headers: await _getHeaders(),
    );

    if (response.statusCode == 200) {
      return User.fromJson(jsonDecode(response.body));
    } else {
      throw Exception('Failed to get user');
    }
  }

  Future<void> deleteUser(String id) async {
    final response = await http.delete(
      Uri.parse('$_baseUrl/auth/users/$id'),
      headers: await _getHeaders(),
    );

    if (response.statusCode != 204) {
      throw Exception('Failed to delete user');
    }
  }

  // --- Feed Service ---

  Future<List<Post>> getAllPublications() async {
    final response = await http.get(
      Uri.parse('$_baseUrl/feed/publications'),
      headers: await _getHeaders(),
    );

    if (response.statusCode == 200) {
      final List<dynamic> data = jsonDecode(response.body);
      return data.map((json) => Post.fromJson(json)).toList();
    } else {
      throw Exception('Failed to load publications');
    }
  }

  Future<Post> createPublication(String title, String content) async {
    final response = await http.post(
      Uri.parse('$_baseUrl/feed/publications'),
      headers: await _getHeaders(),
      body: jsonEncode({'title': title, 'content': content}),
    );

    if (response.statusCode == 201) {
      return Post.fromJson(jsonDecode(response.body));
    } else {
      throw Exception('Failed to create publication');
    }
  }
  
  Future<Post> getPublicationById(String id) async {
    final response = await http.get(
      Uri.parse('$_baseUrl/feed/publications/$id'),
      headers: await _getHeaders(),
    );

    if (response.statusCode == 200) {
      return Post.fromJson(jsonDecode(response.body));
    } else {
      throw Exception('Failed to get publication');
    }
  }

  Future<Post> updatePublication(String id, String title, String content) async {
    final response = await http.put(
      Uri.parse('$_baseUrl/feed/publications/$id'),
      headers: await _getHeaders(),
      body: jsonEncode({'title': title, 'content': content}),
    );

    if (response.statusCode == 200) {
      return Post.fromJson(jsonDecode(response.body));
    } else {
      throw Exception('Failed to update publication');
    }
  }

  Future<void> deletePublication(String id) async {
    final response = await http.delete(
      Uri.parse('$_baseUrl/feed/publications/$id'),
      headers: await _getHeaders(),
    );

    if (response.statusCode != 204) {
      throw Exception('Failed to delete publication');
    }
  }

  Future<List<Post>> getUserPublications(String userId) async {
    final response = await http.get(
      Uri.parse('$_baseUrl/feed/users/$userId/publications'),
      headers: await _getHeaders(),
    );

    if (response.statusCode == 200) {
      final List<dynamic> data = jsonDecode(response.body);
      return data.map((json) => Post.fromJson(json)).toList();
    } else {
      throw Exception('Failed to load user publications');
    }
  }

  Future<Comment> createComment(String postId, String content) async {
    final response = await http.post(
      Uri.parse('$_baseUrl/feed/comments'),
      headers: await _getHeaders(),
      body: jsonEncode({'post_id': postId, 'content': content}),
    );

    if (response.statusCode == 201) {
      return Comment.fromJson(jsonDecode(response.body));
    } else {
      throw Exception('Failed to create comment');
    }
  }

  Future<List<Comment>> getCommentsForPost(String postId) async {
    final response = await http.get(
      Uri.parse('$_baseUrl/feed/comments?post_id=$postId'),
      headers: await _getHeaders(),
    );

    if (response.statusCode == 200) {
      final List<dynamic> data = jsonDecode(response.body);
      return data.map((json) => Comment.fromJson(json)).toList();
    } else {
      throw Exception('Failed to load comments');
    }
  }

  // --- Profile Service ---

  Future<UserProfile> getProfile() async {
    final response = await http.get(
      Uri.parse('$_baseUrl/profile'),
      headers: await _getHeaders(),
    );

    if (response.statusCode == 200) {
      return UserProfile.fromJson(jsonDecode(response.body));
    } else {
      throw Exception('Failed to load user profile');
    }
  }

  Future<UserProfile> createProfile(String name, {String? bio, String? avatarUrl}) async {
    final response = await http.post(
      Uri.parse('$_baseUrl/profile'),
      headers: await _getHeaders(),
      body: jsonEncode({
        'name': name,
        'bio': bio,
        'avatar_url': avatarUrl,
      }),
    );

    if (response.statusCode == 201) {
      return UserProfile.fromJson(jsonDecode(response.body));
    } else {
      throw Exception('Failed to create profile');
    }
  }

  Future<UserProfile> updateProfile({String? name, String? bio, String? avatarUrl}) async {
    final body = <String, String>{};
    if (name != null) body['name'] = name;
    if (bio != null) body['bio'] = bio;
    if (avatarUrl != null) body['avatar_url'] = avatarUrl;

    final response = await http.put(
      Uri.parse('$_baseUrl/profile'),
      headers: await _getHeaders(),
      body: jsonEncode(body),
    );

    if (response.statusCode == 200) {
      return UserProfile.fromJson(jsonDecode(response.body));
    } else {
      throw Exception('Failed to update profile');
    }
  }

  Future<void> deleteProfile() async {
    final response = await http.delete(
      Uri.parse('$_baseUrl/profile'),
      headers: await _getHeaders(),
    );

    if (response.statusCode != 204) {
      throw Exception('Failed to delete profile');
    }
  }
}
