import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:health_buddy_app/models/post.dart';
import 'package:health_buddy_app/screens/profile_screen.dart';
import 'package:health_buddy_app/services/api_service.dart';
import 'package:intl/intl.dart';

class FeedScreen extends StatefulWidget {
  const FeedScreen({super.key});

  @override
  _FeedScreenState createState() => _FeedScreenState();
}

class _FeedScreenState extends State<FeedScreen> {
  final ApiService _apiService = ApiService();
  late Future<List<Post>> _postsFuture;

  @override
  void initState() {
    super.initState();
    _postsFuture = _apiService.getAllPublications();
  }

  @override
  Widget build(BuildContext context) {
    const Color fern = Color(0xFF66BB6A);

    return Scaffold(
      backgroundColor: const Color(0xAFFFFBEF),
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        title: Text(
          'HealthBuddy',
          style: GoogleFonts.roboto(color: fern, fontWeight: FontWeight.bold),
        ),
        actions: [
          GestureDetector(
            onTap: () {
              Navigator.push(
                context,
                MaterialPageRoute(builder: (context) => const ProfileScreen()),
              );
            },
            child: Row(
              children: [
                Text(
                  'My Profile',
                  style: GoogleFonts.roboto(color: fern, fontSize: 16),
                ),
                const SizedBox(width: 8),
                const CircleAvatar(
                  backgroundImage: NetworkImage('https://picsum.photos/seed/my-profile/200'),
                  radius: 20,
                ),
                const SizedBox(width: 16),
              ],
            ),
          ),
        ],
      ),
      body: FutureBuilder<List<Post>>(
        future: _postsFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          } else if (snapshot.hasError) {
            return Center(child: Text('Failed to load posts: ${snapshot.error}'));
          } else if (snapshot.hasData && snapshot.data!.isNotEmpty) {
            final posts = snapshot.data!;
            return ListView.builder(
              padding: const EdgeInsets.all(16.0),
              itemCount: posts.length,
              itemBuilder: (context, index) {
                final post = posts[index];
                return _buildPostCard(post, fern);
              },
            );
          } else {
            return const Center(child: Text('No posts to show.'));
          }
        },
      ),
    );
  }

  Widget _buildPostCard(Post post, Color borderColor) {
    final authorName = post.authorName.trim().isNotEmpty ? post.authorName : 'Unknown Author';

    return Card(
      margin: const EdgeInsets.only(bottom: 16.0),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
        side: BorderSide(color: borderColor, width: 1.5),
      ),
      elevation: 2,
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            _buildAuthorInfo(authorName, post.createdAt, userId: post.userId),
            const SizedBox(height: 16),
            if (post.title.isNotEmpty)
              Padding(
                padding: const EdgeInsets.only(bottom: 8.0),
                child: Text(
                  post.title,
                  style: GoogleFonts.roboto(fontSize: 18, fontWeight: FontWeight.bold),
                ),
              ),
            Text(
              post.content,
              style: GoogleFonts.roboto(fontSize: 15, height: 1.4),
            ),
            const SizedBox(height: 16),
            Row(
              mainAxisAlignment: MainAxisAlignment.end,
              children: [
                TextButton.icon(
                  onPressed: () {},
                  icon: const Icon(Icons.comment_outlined, color: Colors.grey),
                  label: Text('Comment', style: GoogleFonts.roboto(color: Colors.grey)),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildAuthorInfo(String name, DateTime postDate, {String? userId}) {
    final avatar = userId != null
        ? NetworkImage('https://picsum.photos/seed/$userId/200')
        : const AssetImage('assets/default_avatar.png') as ImageProvider;

    return Row(
      children: [
        CircleAvatar(
          backgroundImage: avatar,
          radius: 20,
        ),
        const SizedBox(width: 12),
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              name,
              style: GoogleFonts.roboto(fontWeight: FontWeight.bold, fontSize: 16),
            ),
            Text(
              DateFormat.yMMMd().add_jm().format(postDate),
              style: GoogleFonts.roboto(color: Colors.grey[600], fontSize: 12),
            ),
          ],
        ),
      ],
    );
  }
}
