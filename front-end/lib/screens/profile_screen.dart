import 'dart:io';
import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:health_buddy_app/models/post.dart';
import 'package:health_buddy_app/models/user.dart';
import 'package:health_buddy_app/screens/login_screen.dart';
import 'package:health_buddy_app/screens/post_creation_screen.dart';
import 'package:health_buddy_app/services/api_service.dart';
import 'package:image_picker/image_picker.dart';
import 'package:intl/intl.dart';

class ProfileScreen extends StatefulWidget {
  const ProfileScreen({super.key});

  @override
  _ProfileScreenState createState() => _ProfileScreenState();
}

class _ProfileScreenState extends State<ProfileScreen> {
  late Future<UserProfile> _userProfileFuture;
  final ApiService _apiService = ApiService();
  final ImagePicker _picker = ImagePicker();
  bool _isUpdatingAvatar = false;

  @override
  void initState() {
    super.initState();
    _loadProfile();
  }

  Future<void> _loadProfile() async {
    setState(() {
      _userProfileFuture = _apiService.getProfile();
    });
  }

  Future<void> _updateAvatar() async {
    try {
      final XFile? pickedFile = await _picker.pickImage(
        source: ImageSource.gallery,
        maxWidth: 800,
        maxHeight: 800,
        imageQuality: 85,
      );

      if (pickedFile != null) {
        setState(() => _isUpdatingAvatar = true);
        final updatedProfile = await _apiService.updateProfile(
          avatarUrl: pickedFile.path,
        );
        setState(() {
          _userProfileFuture = Future.value(updatedProfile);
        });
      }
    } catch (e) {
      debugPrint('Error picking image: $e');
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error: ${e.toString()}')),
        );
      }
    } finally {
      setState(() => _isUpdatingAvatar = false);
    }
  }

  Future<void> _showEditBioDialog(UserProfile currentUser) async {
    final bioController = TextEditingController(text: currentUser.bio);
    return showDialog<void>(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          title: const Text('Edit Bio'),
          content: TextField(
            controller: bioController,
            decoration: const InputDecoration(hintText: "Enter your bio here"),
            autofocus: true,
            maxLines: 3,
          ),
          actions: <Widget>[
            TextButton(
              child: const Text('Cancel'),
              onPressed: () {
                Navigator.of(context).pop();
              },
            ),
            TextButton(
              child: const Text('Save'),
              onPressed: () {
                Navigator.of(context).pop();
                _updateBio(bioController.text);
              },
            ),
          ],
        );
      },
    );
  }

  Future<void> _updateBio(String newBio) async {
    try {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Saving bio...')),
      );

      final updatedProfile = await _apiService.updateProfile(bio: newBio);
      setState(() {
        _userProfileFuture = Future.value(updatedProfile);
      });
      if (mounted) {
        ScaffoldMessenger.of(context).hideCurrentSnackBar();
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Bio updated successfully! ✅')),
        );
      }
    } catch (e) {
      debugPrint('Failed to update bio: $e');
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to update bio: ${e.toString()}')),
        );
      }
    }
  }

  void _logout() {
    Navigator.pushAndRemoveUntil(
      context,
      MaterialPageRoute(builder: (_) => const LoginScreen()),
      (route) => false,
    );
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
        actions: [
          IconButton(
            icon: const Icon(Icons.logout),
            color: fern,
            onPressed: _logout,
            tooltip: 'Log Out',
          ),
        ],
      ),
      body: FutureBuilder<UserProfile>(
        future: _userProfileFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          } else if (snapshot.hasError) {
            return Center(child: Text('Error: \${snapshot.error}'));
          } else if (snapshot.hasData) {
            final userProfile = snapshot.data!;
            return RefreshIndicator(
              onRefresh: _loadProfile,
              child: NestedScrollView(
                headerSliverBuilder: (context, innerBoxIsScrolled) {
                  return [
                    SliverToBoxAdapter(child: _buildProfileHeader(userProfile, fern)),
                  ];
                },
                body: _buildPostsList(userProfile.posts, fern),
              ),
            );
          } else {
            return const Center(child: Text('No profile data found.'));
          }
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () async {
          final result = await Navigator.push(
            context,
            MaterialPageRoute(builder: (context) => const PostCreationScreen()),
          );
          if (result == true) {
            _loadProfile();
          }
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
          Stack(
            alignment: Alignment.center,
            children: [
              if (_isUpdatingAvatar)
                const CircularProgressIndicator()
              else
                GestureDetector(
                  onTap: _updateAvatar,
                  child: CircleAvatar(
                    radius: 50,
                    backgroundImage: _getAvatarImageProvider(user.avatarUrl),
                  ),
                ),
              if (!_isUpdatingAvatar)
                Positioned(
                  bottom: 0,
                  right: 0,
                  child: Container(
                    padding: const EdgeInsets.all(4),
                    decoration: BoxDecoration(
                      color: fernColor,
                      shape: BoxShape.circle,
                    ),
                    child: const Icon(
                      Icons.camera_alt,
                      size: 20,
                      color: Colors.white,
                    ),
                  ),
                ),
            ],
          ),
          const SizedBox(height: 16),
          Text(
            user.name,
            style: GoogleFonts.roboto(fontSize: 24, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 4),
          Text(
            '@\${user.userId}',
            style: GoogleFonts.roboto(fontSize: 16, color: Colors.grey[600]),
          ),
          const SizedBox(height: 12),
          GestureDetector(
            onTap: () => _showEditBioDialog(user),
            child: Text(
              user.bio.isNotEmpty ? user.bio : 'Enter your bio',
              textAlign: TextAlign.center,
              style: user.bio.isNotEmpty
                  ? GoogleFonts.roboto(fontSize: 15, height: 1.4)
                  : GoogleFonts.roboto(
                      fontSize: 15,
                      height: 1.4,
                      fontStyle: FontStyle.italic,
                      color: Colors.grey[500],
                    ),
            ),
          ),
        ],
      ),
    );
  }

  ImageProvider _getAvatarImageProvider(String avatarUrl) {
    try {
      if (avatarUrl.startsWith('http')) {
        return NetworkImage(avatarUrl);
      } else if (avatarUrl.isNotEmpty) {
        return FileImage(File(avatarUrl));
      }
    } catch (e) {
      debugPrint('Error while loading photo: \$e');
    }
    return const AssetImage('assets/default_avatar.png');
  }

  Widget _buildPostsList(List<Post> posts, Color fernColor) {
    if (posts.isEmpty) {
      return Center(
        child: Text(
          'You have no posts yet.\nTap the + button to create one!',
          textAlign: TextAlign.center,
          style: GoogleFonts.roboto(color: Colors.grey[600], fontSize: 16),
        ),
      );
    }
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
            if (post.title.isNotEmpty)
              Padding(
                padding: const EdgeInsets.only(bottom: 6.0),
                child: Text(
                  post.title,
                  style: GoogleFonts.roboto(fontSize: 18, fontWeight: FontWeight.bold),
                ),
              ),
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
