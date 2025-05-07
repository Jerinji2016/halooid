import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class ApiConfig {
  static const String baseUrl = 'https://api.halooid.com/v1';
  
  static Future<Map<String, String>> getHeaders() async {
    const storage = FlutterSecureStorage();
    final token = await storage.read(key: 'auth_token');
    
    return {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
      if (token != null) 'Authorization': 'Bearer $token',
    };
  }
}
