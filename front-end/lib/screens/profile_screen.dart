import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:health_buddy_app/models/post.dart';
import 'package:health_buddy_app/models/user.dart';
import 'package:health_buddy_app/screens/post_creation_screen.dart';
import 'package:health_buddy_app/services/api_service.dart';
import 'package:intl/intl.dart';

class ProfileScreen extends StatefulWidget {
  const ProfileScreen({super.key});

  @override
  _ProfileScreenState createState() => _ProfileScreenState();
}

class _ProfileScreenState extends State<ProfileScreen> {
  late Future<UserProfile> _userProfileFuture;
  final ApiService _apiService = ApiService();

  @override
  Future<void> initState() async {
    super.initState();
    _userProfileFuture = (await _apiService.getProfile()) as Future<UserProfile>;
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
        iconTheme: const IconThemeData(color: fern),
      ),
      body: FutureBuilder<UserProfile>(
        future: _userProfileFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          } else if (snapshot.hasError) {
            return Center(child: Text('Error: ${snapshot.error}'));
          } else if (snapshot.hasData) {
            final userProfile = snapshot.data!;
            return NestedScrollView(
              headerSliverBuilder: (context, innerBoxIsScrolled) {
                return [
                  SliverToBoxAdapter(child: _buildProfileHeader(userProfile, fern)),
                ];
              },
              body: _buildPostsList(userProfile.posts, fern),
            );
          } else {
            return const Center(child: Text('No profile data found.'));
          }
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          Navigator.push(
            context,
            MaterialPageRoute(builder: (context) => const PostCreationScreen()),
          );
        },
        backgroundColor: fern,
        child: const Icon(Icons.add, color: Colors.white),
      ),
    );
  }

  Widget _buildProfileHeader(UserProfile user, Color fernColor) {
    return Container(
      padding: const EdgeInsets.all(24.0),
      child: Column(
        children: [
          CircleAvatar(
            radius: 50,
            backgroundImage: NetworkImage('https://picsum.photos/seed/${user.userId}/200'),
          ),
          const SizedBox(height: 16),
          Text(
            user.name,
            style: GoogleFonts.roboto(fontSize: 24, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 4),
          Text(
            '@${user.userId}',
            style: GoogleFonts.roboto(fontSize: 16, color: Colors.grey[600]),
          ),
          const SizedBox(height: 12),
          Text(
            user.bio,
            textAlign: TextAlign.center,
            style: GoogleFonts.roboto(fontSize: 15, height: 1.4),
          ),
        ],
      ),
    );
  }

  Widget _buildPostsList(List<Post> posts, Color fernColor) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 16.0),
          child: Text(
            'My Posts',
            style: GoogleFonts.roboto(
              fontSize: 20,
              fontWeight: FontWeight.bold,
              color: fernColor,
            ),
          ),
        ),
        const SizedBox(height: 16),
        Expanded(
          child: ListView.builder(
            padding: const EdgeInsets.symmetric(horizontal: 16.0),
            itemCount: posts.length,
            itemBuilder: (context, index) {
              return _buildPostCard(posts[index], fernColor);
            },
          ),
        ),
      ],
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
            Text(
              DateFormat.yMMMd().add_jm().format(post.createdAt),
              style: GoogleFonts.roboto(color: Colors.grey[600], fontSize: 12),
            ),
            const SizedBox(height: 8),
            Text(
              post.content,
              style: GoogleFonts.roboto(fontSize: 15, height: 1.4),
            ),
          ],
        ),
      ),
    );
  }
}
