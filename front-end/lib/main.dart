import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:health_buddy_app/screens/feed_screen.dart';
import 'package:health_buddy_app/screens/home_screen.dart';
import 'package:health_buddy_app/screens/login_screen.dart';
import 'package:health_buddy_app/screens/profile_screen.dart';
import 'package:health_buddy_app/screens/registration_screen.dart';
import 'package:health_buddy_app/services/api_service.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    final textTheme = Theme.of(context).textTheme;

    return MaterialApp(
      title: 'HealthBuddy',
      theme: ThemeData(
        primaryColor: const Color(0xFF66BB6A),
        scaffoldBackgroundColor: const Color(0xAFFFFBEF),
        textTheme: GoogleFonts.robotoTextTheme(textTheme),
        appBarTheme: const AppBarTheme(
          backgroundColor: Colors.transparent,
          elevation: 0,
          iconTheme: IconThemeData(color: Color(0xFF66BB6A)),
        ),
      ),
      debugShowCheckedModeBanner: false,
      home: const AuthCheck(),
      routes: {
        '/login': (context) => const LoginScreen(),
        '/register': (context) => const RegistrationScreen(),
        '/feed': (context) => const FeedScreen(),
        '/profile': (context) => const ProfileScreen(),
        '/home': (context) => const HomeScreen(),
      },
    );
  }
}

class AuthCheck extends StatefulWidget {
  const AuthCheck({super.key});

  @override
  _AuthCheckState createState() => _AuthCheckState();
}

class _AuthCheckState extends State<AuthCheck> {
  final ApiService _apiService = ApiService();
  late Future<String?> _tokenFuture;

  @override
  void initState() {
    super.initState();
    _tokenFuture = _apiService.getToken();
  }

  @override
  Widget build(BuildContext context) {
    return FutureBuilder<String?>(
      future: _tokenFuture,
      builder: (context, snapshot) {
        if (snapshot.connectionState == ConnectionState.waiting) {
          return const Scaffold(
            body: Center(child: CircularProgressIndicator()),
          );
        }

        if (snapshot.hasData && snapshot.data != null) {
          // TODO: Add token validation logic
          return const FeedScreen();
        } else {
          return const HomeScreen();
        }
      },
    );
  }
}