class ApiConstants {
  static const String baseUrl = 'http://localhost:8080';
  static const String version = 'v1';
  static const String apiPath = '/api/$version';
  
  // Auth endpoints
  static const String login = '$apiPath/auth/login';
  static const String refresh = '$apiPath/auth/refresh';
  static const String me = '$apiPath/auth/me';
  
  // Faction endpoints
  static const String factions = '$apiPath/factions';
  
  // Geometry endpoints
  static const String geometries = '$apiPath/geometries';
  
  // Search endpoints
  static const String search = '$apiPath/search';
  
  // Health endpoints
  static const String health = '/health';
}