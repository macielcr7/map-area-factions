class Faction {
  final String id;
  final String name;
  final String description;
  final String color;
  final bool active;
  final DateTime createdAt;

  const Faction({
    required this.id,
    required this.name,
    required this.description,
    required this.color,
    required this.active,
    required this.createdAt,
  });

  factory Faction.fromJson(Map<String, dynamic> json) {
    return Faction(
      id: json['id'] as String,
      name: json['name'] as String,
      description: json['description'] as String,
      color: json['color'] as String,
      active: json['active'] as bool,
      createdAt: DateTime.parse(json['created_at'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'description': description,
      'color': color,
      'active': active,
      'created_at': createdAt.toIso8601String(),
    };
  }
}