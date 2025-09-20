# Map Factions App - Flutter

Aplicativo multiplataforma (Android, iOS, Web, Desktop) para visualização de áreas dominadas por facções.

## 🏗️ Arquitetura

### Stack Tecnológico
- **Framework**: Flutter 3.16+
- **Mapas**: Mapbox Maps Flutter / Google Maps
- **Estado**: Riverpod
- **HTTP**: Dio
- **WebSocket**: Socket.io
- **Cache**: Hive
- **Navegação**: Go Router
- **UI**: Material 3

### Plataformas Suportadas
- ✅ Android (API 21+)
- ✅ iOS (iOS 12+)
- ✅ Web (Chrome, Firefox, Safari, Edge)
- ✅ Windows (Windows 10+)
- ✅ macOS (macOS 10.14+)
- ✅ Linux (Ubuntu 18.04+)

## 📁 Estrutura do Projeto

```
app/
├── lib/
│   ├── src/
│   │   ├── models/           # Modelos de dados
│   │   ├── services/         # Serviços (API, Cache, etc)
│   │   ├── providers/        # Providers (Riverpod)
│   │   ├── screens/          # Telas da aplicação
│   │   ├── widgets/          # Widgets reutilizáveis
│   │   └── utils/            # Utilitários e helpers
│   ├── main.dart             # Entry point
│   └── app.dart              # Widget principal
├── android/                  # Configurações Android
├── ios/                      # Configurações iOS
├── web/                      # Configurações Web
├── macos/                    # Configurações macOS
├── linux/                    # Configurações Linux
├── windows/                  # Configurações Windows
├── assets/                   # Assets (images, fonts, etc)
├── test/                     # Testes
└── docs/                     # Documentação específica
```

## 🚀 Funcionalidades

### 🗺️ Mapa Interativo
- Visualização de polígonos por facção
- Cores e rótulos configuráveis
- Zoom e navegação suaves
- Busca por localização
- "Perto de mim" com GPS

### 🔍 Busca e Filtros
- Busca por bairro/cidade/CEP
- Geocodificação automática
- Filtros por facção e risco
- Histórico de buscas
- Sugestões inteligentes

### 💳 Sistema de Assinaturas
- Planos: Mensal, Trimestral, Vitalício
- Integração Mercado Pago + Pix
- Paywall inteligente
- Gestão de recibos
- Renovação automática

### 📱 Recursos Móveis
- Notificações push
- Cache offline por cidade
- Sincronização em background
- Touch gestures otimizados
- Modo escuro/claro

### 🖥️ Recursos Desktop
- Layouts adaptáveis
- Keyboard shortcuts
- Menu bar nativo
- Window management
- Sistema de tray

## 🎯 Telas Principais

### OnboardingScreen
- Introdução ao app
- Termos de uso e privacidade
- Permissões necessárias
- Tutorial de funcionalidades

### AuthScreen
- Login/Registro
- Autenticação social
- Recuperação de senha
- Validação de email

### MapScreen (Principal)
- Mapa com polígonos
- Controles de zoom
- Botão "Minha localização"
- Legenda de facções
- Menu de filtros

### SearchScreen
- Campo de busca
- Resultados em lista
- Mapa de resultados
- Histórico de buscas

### SubscriptionScreen
- Planos disponíveis
- Comparativo de recursos
- Processo de pagamento
- Status da assinatura

### ProfileScreen
- Dados do usuário
- Configurações
- Histórico de pagamentos
- Suporte e contato

### SettingsScreen
- Preferências do app
- Cache offline
- Notificações
- Sobre o app

## 📡 Integração com API

### HTTP Client (Dio)
```dart
class ApiService {
  final Dio _dio = Dio(BaseOptions(
    baseUrl: Config.apiUrl,
    connectTimeout: Duration(seconds: 5),
    receiveTimeout: Duration(seconds: 10),
  ));

  Future<List<Faction>> getFactions() async {
    final response = await _dio.get('/factions');
    return (response.data['data'] as List)
        .map((json) => Faction.fromJson(json))
        .toList();
  }
}
```

### WebSocket (Socket.io)
```dart
class RealtimeService {
  late IO.Socket socket;

  void connect() {
    socket = IO.io(Config.wsUrl, <String, dynamic>{
      'transports': ['websocket'],
      'autoConnect': false,
    });

    socket.on('geometry.updated', (data) {
      // Atualizar mapa em tempo real
    });

    socket.connect();
  }
}
```

## 🗺️ Integração com Mapas

### Mapbox Maps Flutter
```dart
class MapWidget extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return MapboxMap(
      styleString: MapboxStyles.MAPBOX_STREETS,
      initialCameraPosition: CameraPosition(
        target: LatLng(-3.7172, -38.5267), // Fortaleza
        zoom: 10,
      ),
      onMapCreated: (MapboxMapController controller) {
        _loadGeometries(controller);
      },
    );
  }
}
```

### Camadas Dinâmicas
- Source GeoJSON para polígonos
- Estilização por facção
- Clustering de incidentes
- Rótulos responsivos

## 💾 Cache e Offline

### Hive Database
```dart
@HiveType(typeId: 0)
class CachedGeometry extends HiveObject {
  @HiveField(0)
  final String id;
  
  @HiveField(1)
  final Map<String, dynamic> geoJson;
  
  @HiveField(2)
  final DateTime cachedAt;
}
```

### Estratégia de Cache
- Cache por cidade (limite configurável)
- TTL de 24 horas para dados
- Sync em background
- Indicadores de dados offline

## 🔔 Notificações Push

### Firebase Cloud Messaging
```dart
class NotificationService {
  static Future<void> initialize() async {
    await Firebase.initializeApp();
    
    FirebaseMessaging.onMessage.listen((RemoteMessage message) {
      // Mostrar notificação local
    });
    
    FirebaseMessaging.onBackgroundMessage(_firebaseMessagingBackgroundHandler);
  }
}
```

### Tipos de Notificação
- Mudança de área próxima
- Novos incidentes
- Atualizações de assinatura
- Alertas de segurança

## 🎨 Design System

### Material 3
- Componentes padronizados
- Color schemes dinâmicos
- Typography scales
- Adaptive layouts

### Temas
```dart
ThemeData lightTheme = ThemeData(
  useMaterial3: true,
  colorScheme: ColorScheme.fromSeed(
    seedColor: Colors.red,
    brightness: Brightness.light,
  ),
);

ThemeData darkTheme = ThemeData(
  useMaterial3: true,
  colorScheme: ColorScheme.fromSeed(
    seedColor: Colors.red,
    brightness: Brightness.dark,
  ),
);
```

## 🧪 Testes

### Estrutura de Testes
```bash
test/
├── unit/                 # Testes unitários
├── widget/               # Testes de widgets
├── integration/          # Testes de integração
└── golden/               # Golden tests
```

### Comandos
```bash
# Testes unitários
flutter test

# Testes com cobertura
flutter test --coverage

# Testes de integração
flutter test integration_test/

# Golden tests
flutter test --update-goldens
```

## 📦 Build e Deploy

### Android
```bash
# Debug APK
flutter build apk --debug

# Release APK
flutter build apk --release

# App Bundle (Play Store)
flutter build appbundle --release
```

### iOS
```bash
# iOS build
flutter build ios --release

# Archive for App Store
flutter build ipa --release
```

### Web
```bash
# Web build
flutter build web --release

# PWA optimized
flutter build web --pwa-strategy offline-first
```

### Desktop
```bash
# Windows
flutter build windows --release

# macOS
flutter build macos --release

# Linux
flutter build linux --release
```

## 🔧 Configuração de Ambiente

### pubspec.yaml (principais dependencies)
```yaml
dependencies:
  flutter:
    sdk: flutter
  
  # Estado
  flutter_riverpod: ^2.4.9
  riverpod_annotation: ^2.3.3
  
  # HTTP & WebSocket
  dio: ^5.4.0
  socket_io_client: ^2.0.3+1
  
  # Mapas
  mapbox_maps_flutter: ^1.0.0
  geolocator: ^10.1.0
  geocoding: ^2.1.1
  
  # Cache & Storage
  hive: ^2.2.3
  hive_flutter: ^1.1.0
  shared_preferences: ^2.2.2
  
  # UI & Navigation
  go_router: ^12.1.3
  flutter_animate: ^4.5.0
  
  # Utils
  intl: ^0.19.0
  url_launcher: ^6.2.2
  package_info_plus: ^4.2.0
```

## 📱 Platform-Specific Features

### Android
- Adaptive icons
- Splash screen API
- Background services
- Deep links

### iOS
- Widget extensions
- Shortcuts
- Background refresh
- Universal links

### Web
- PWA capabilities
- Service workers
- Web share API
- Responsive design

### Desktop
- Native menu bars
- System tray
- File system access
- Window management

## 🔐 Segurança

### Proteção de APIs
- Certificate pinning
- Token encryption
- Biometric authentication
- Secure storage

### Privacidade
- Location permissions
- Data minimization
- Analytics opt-out
- LGPD compliance