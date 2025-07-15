import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:health_buddy_app/screens/registration_screen.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    const Color fern = Color(0xFF66BB6A);

    return Scaffold(
      backgroundColor: const Color(0xAFFFFBEF),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(
              'Welcome!',
              style: GoogleFonts.roboto(
                fontSize: 48,
                fontWeight: FontWeight.bold,
                color: fern,
              ),
            ),
            const SizedBox(height: 16),
            Text(
              'Share your journey. Inspire the change.',
              style: GoogleFonts.roboto(
                fontSize: 18,
                color: const Color(0x73000000),
              ),
            ),
            const SizedBox(height: 48),
            ElevatedButton(
              onPressed: () {
                // Logic: If logged in -> Profile, Else -> Registration
                // For now, let's just navigate to registration
                Navigator.push(
                  context,
                  MaterialPageRoute(builder: (context) => const RegistrationScreen()),
                );
              },
              style: ElevatedButton.styleFrom(
                backgroundColor: fern,
                padding: const EdgeInsets.symmetric(horizontal: 40, vertical: 16),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(30),
                ),
              ),
              child: Text(
                'Get Started',
                style: GoogleFonts.roboto(
                  fontSize: 18,
                  color: Colors.white,
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
