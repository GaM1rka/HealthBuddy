import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';
import 'package:health_buddy_app/models/post.dart';
import 'package:health_buddy_app/models/user.dart';

class ApiService {
  final String _baseUrl = 'https://f9428f71-7371-4774-8fa6-af08d9664808.mock.pstmn.io'; // Replace with your actual backend URL

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

  Future<String> login(String email, String password) async {
    final response = await http.post(
      Uri.parse('$_baseUrl/api/login'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({'email': email, 'password': password}),
    );

    print('Response status: ${response.statusCode}');
    print('Response body: ${response.body}');

    if (response.statusCode == 200) {
      try {
        final responseData = jsonDecode(response.body) as Map<String, dynamic>;
        
        // Извлекаем токен из тела ответа
        final token = responseData['token'] as String?;
        
        if (token != null && token.isNotEmpty) {
          await _saveToken(token);
          return token;
        } else {
          throw Exception('Token not found in response body');
        }
      } catch (e) {
        throw Exception('Failed to parse response: $e');
      }
    } else {
      // Пытаемся извлечь сообщение об ошибке из тела ответа
      try {
        final errorData = jsonDecode(response.body) as Map<String, dynamic>;
        throw Exception(errorData['message'] ?? 'Login failed with status ${response.statusCode}');
      } catch (_) {
        throw Exception('Login failed with status ${response.statusCode}');
      }
    }
  }


  Future<String> register(String username, String email, String password) async {
    final response = await http.post(
      Uri.parse('$_baseUrl/api/register'),
      headers: {'Content-Type': 'application/json'}, // Не используем _getHeaders() для регистрации
      body: jsonEncode({
        'username': username,
        'email': email,
        'password': password,
      }),
    );

    print('Registration response status: ${response.statusCode}');
    print('Registration response body: ${response.body}');

    if (response.statusCode == 200 || response.statusCode == 201) {
      try {
        final responseData = jsonDecode(response.body) as Map<String, dynamic>;
        
        // Извлекаем токен из тела ответа (аналогично логину)
        final token = responseData['token'] as String?;
        
        if (token != null && token.isNotEmpty) {
          await _saveToken(token);
          return token;
        } else {
          throw Exception('Token not found in response body');
        }
      } catch (e) {
        throw Exception('Failed to parse registration response: $e');
      }
    } else {
      try {
        final errorData = jsonDecode(response.body) as Map<String, dynamic>;
        throw Exception(errorData['message'] ?? 'Registration failed with status ${response.statusCode}');
      } catch (_) {
        throw Exception('Registration failed with status ${response.statusCode}');
      }
    }
  }

  Future<List<Post>> getPosts() async {
    final response = await http.get(
      Uri.parse('$_baseUrl/posts'),
      headers: await _getHeaders(),
    );

    if (response.statusCode == 200) {
      final List<dynamic> data = jsonDecode(response.body);
      return data.map((json) => Post.fromJson(json)).toList();
    } else {
      throw Exception('Failed to load posts');
    }
  }

  Future<void> createPost(String content) async {
    final response = await http.post(
      Uri.parse('$_baseUrl/post/create'),
      headers: await _getHeaders(),
      body: jsonEncode({'content': content}),
    );

    if (response.statusCode != 201) {
      throw Exception('Failed to create post');
    }
  }

  Future<UserProfile> getUserProfile() async {
    final response = await http.get(
      Uri.parse('$_baseUrl/profile'), // Assuming this is the endpoint
      headers: await _getHeaders(),
    );

    if (response.statusCode == 200) {
      //TODO: Update this once the UserProfile model has a fromJson factory
      final Map<String, dynamic> data = jsonDecode(response.body);
      return UserProfile.fromJson(data);
    } else {
      throw Exception('Failed to load user profile');
    }
  }
  
    Future<List<Post>> getUserPosts() async {
    final response = await http.get(
      Uri.parse('$_baseUrl/user/posts'), // Assuming this is the endpoint
      headers: await _getHeaders(),
    );

    if (response.statusCode == 200) {
      final List<dynamic> data = jsonDecode(response.body);
      return data.map((json) => Post.fromJson(json)).toList();
    } else {
      throw Exception('Failed to load user posts');
    }
  }
}