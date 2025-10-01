import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

class SearchPage extends ConsumerStatefulWidget {
  const SearchPage({super.key});

  @override
  ConsumerState<SearchPage> createState() => _SearchPageState();
}

class _SearchPageState extends ConsumerState<SearchPage> {
  final _searchController = TextEditingController();
  final List<String> _recentSearches = [
    'Messejana, Fortaleza-CE',
    'Centro, São Paulo-SP',
    'Copacabana, Rio de Janeiro-RJ',
    'Boa Viagem, Recife-PE',
  ];

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Buscar Localização'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
      ),
      body: Column(
        children: [
          // Search bar
          Padding(
            padding: const EdgeInsets.all(16),
            child: TextField(
              controller: _searchController,
              decoration: InputDecoration(
                hintText: 'Digite um bairro, cidade ou endereço...',
                prefixIcon: const Icon(Icons.search),
                suffixIcon: _searchController.text.isNotEmpty
                    ? IconButton(
                        icon: const Icon(Icons.clear),
                        onPressed: () {
                          _searchController.clear();
                          setState(() {});
                        },
                      )
                    : IconButton(
                        icon: const Icon(Icons.my_location),
                        onPressed: () {
                          _getCurrentLocation();
                        },
                      ),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
              ),
              onChanged: (value) {
                setState(() {});
              },
              onSubmitted: (value) {
                _performSearch(value);
              },
            ),
          ),

          // Quick location button
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: SizedBox(
              width: double.infinity,
              child: ElevatedButton.icon(
                onPressed: _getCurrentLocation,
                icon: const Icon(Icons.my_location),
                label: const Text('Usar minha localização'),
                style: ElevatedButton.styleFrom(
                  padding: const EdgeInsets.symmetric(vertical: 12),
                ),
              ),
            ),
          ),

          const SizedBox(height: 24),

          // Recent searches
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Padding(
                  padding: EdgeInsets.symmetric(horizontal: 16),
                  child: Text(
                    'Buscas recentes',
                    style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
                const SizedBox(height: 12),
                Expanded(
                  child: ListView.builder(
                    itemCount: _recentSearches.length,
                    itemBuilder: (context, index) {
                      final search = _recentSearches[index];
                      return ListTile(
                        leading: const Icon(Icons.history),
                        title: Text(search),
                        trailing: IconButton(
                          icon: const Icon(Icons.close),
                          onPressed: () {
                            setState(() {
                              _recentSearches.removeAt(index);
                            });
                          },
                        ),
                        onTap: () {
                          _performSearch(search);
                        },
                      );
                    },
                  ),
                ),
                
                // Popular locations
                const Padding(
                  padding: EdgeInsets.all(16),
                  child: Text(
                    'Locais populares',
                    style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
                
                _buildPopularLocation(
                  'Fortaleza, CE',
                  'Capital do Ceará',
                  Icons.location_city,
                ),
                _buildPopularLocation(
                  'São Paulo, SP',
                  'Capital de São Paulo',
                  Icons.location_city,
                ),
                _buildPopularLocation(
                  'Rio de Janeiro, RJ',
                  'Cidade Maravilhosa',
                  Icons.location_city,
                ),
                _buildPopularLocation(
                  'Recife, PE',
                  'Capital de Pernambuco',
                  Icons.location_city,
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildPopularLocation(String title, String subtitle, IconData icon) {
    return ListTile(
      leading: Icon(icon),
      title: Text(title),
      subtitle: Text(subtitle),
      onTap: () {
        _performSearch(title);
      },
    );
  }

  void _getCurrentLocation() {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('Obtendo sua localização...'),
        duration: Duration(seconds: 2),
      ),
    );

    // Simulate location detection
    Future.delayed(const Duration(seconds: 2), () {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Localização encontrada: Messejana, Fortaleza-CE'),
            backgroundColor: Colors.green,
          ),
        );
        context.pop();
      }
    });
  }

  void _performSearch(String query) {
    if (query.trim().isEmpty) return;

    // Add to recent searches if not already there
    if (!_recentSearches.contains(query)) {
      setState(() {
        _recentSearches.insert(0, query);
        if (_recentSearches.length > 10) {
          _recentSearches.removeLast();
        }
      });
    }

    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text('Buscando: $query'),
        action: SnackBarAction(
          label: 'Ver no mapa',
          onPressed: () {
            context.pop();
          },
        ),
      ),
    );

    // Simulate search and return to map
    Future.delayed(const Duration(seconds: 1), () {
      if (mounted) {
        context.pop();
      }
    });
  }
}