import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:health_buddy_app/models/post.dart';
import 'package:health_buddy_app/screens/profile_screen.dart';
import 'package:intl/intl.dart';

class FeedScreen extends StatefulWidget {
  const FeedScreen({super.key});

  @override
  _FeedScreenState createState() => _FeedScreenState();
}

class _FeedScreenState extends State<FeedScreen> {
  // Mock data for now
  final List<Post> _posts = [
    Post(
      postId: '1',
      userId: 'user1',
      authorName: 'John Doe',
      content: 'Just finished a 5k run! Feeling great. #health #running',
      createdAt: DateTime.now().subtract(const Duration(hours: 1)),
    ),
    Post(
      postId: '2',
      userId: 'user2',
      authorName: 'Jane Smith',
      content: 'My new favorite healthy recipe: Quinoa salad with avocado and chickpeas. So delicious and nutritious!',
      createdAt: DateTime.now().subtract(const Duration(days: 1)),
    ),
    Post(
      postId: '3',
      userId: 'user1',
      authorName: 'John Doe',
      content: 'Starting a new yoga challenge tomorrow. Wish me luck!',
      createdAt: DateTime.now().subtract(const Duration(days: 2)),
    ),
  ];

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
                  'My Profile', // Replace with actual user name
                  style: GoogleFonts.roboto(color: fern, fontSize: 16),
                ),
                const SizedBox(width: 8),
                const CircleAvatar(
                  // Replace with actual user avatar
                  backgroundImage: NetworkImage('https://picsum.photos/seed/my-profile/200'),
                  radius: 20,
                ),
                const SizedBox(width: 16),
              ],
            ),
          ),
        ],
      ),
      body: ListView.builder(
        padding: const EdgeInsets.all(16.0),
        itemCount: _posts.length,
        itemBuilder: (context, index) {
          final post = _posts[index];
          return _buildPostCard(post, fern);
        },
      ),
    );
  }

  Widget _buildPostCard(Post post, Color borderColor) {
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
            Row(
              children: [
                CircleAvatar(
                  backgroundImage: NetworkImage('https://picsum.photos/seed/${post.userId}/200'),
                  radius: 20,
                ),
                const SizedBox(width: 12),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      post.authorName,
                      style: GoogleFonts.roboto(fontWeight: FontWeight.bold, fontSize: 16),
                    ),
                    Text(
                      DateFormat.yMMMd().add_jm().format(post.createdAt),
                      style: GoogleFonts.roboto(color: Colors.grey[600], fontSize: 12),
                    ),
                  ],
                ),
              ],
            ),
            const SizedBox(height: 16),
            Text(
              post.content,
              style: GoogleFonts.roboto(fontSize: 15, height: 1.4),
            ),
            const SizedBox(height: 16),
            Row(
              mainAxisAlignment: MainAxisAlignment.end,
              children: [
                TextButton.icon(
                  onPressed: () {
                    // TODO: Implement comment functionality
                  },
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
}
