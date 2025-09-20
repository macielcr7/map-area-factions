import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:map_factions_app/core/theme/app_theme.dart';
import 'package:map_factions_app/core/router/app_router.dart';
import 'package:map_factions_app/core/providers/auth_provider.dart';

void main() {
  runApp(
    const ProviderScope(
      child: MapFactionsApp(),
    ),
  );
}

class MapFactionsApp extends ConsumerWidget {
  const MapFactionsApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final router = ref.watch(routerProvider);
    
    return MaterialApp.router(
      title: 'Map Area Factions',
      theme: AppTheme.light,
      darkTheme: AppTheme.dark,
      themeMode: ThemeMode.system,
      routerConfig: router,
      debugShowCheckedModeBanner: false,
    );
  }
}