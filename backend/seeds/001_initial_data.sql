-- Seed data para desenvolvimento
-- Dados de exemplo para facções, regiões e usuários

-- Inserir facções conhecidas
INSERT INTO factions (id, name, acronym, color_hex, display_priority, active) VALUES
-- Facções do Ceará
(uuid_generate_v4(), 'Terceiro Comando de Fortaleza', 'TDN', '#ff0000', 10, true),
(uuid_generate_v4(), 'Guardiões do Estado', 'GDE', '#0000ff', 9, true),
(uuid_generate_v4(), 'Família do Norte', 'FDN', '#ff8800', 8, true),

-- Facções nacionais
(uuid_generate_v4(), 'Comando Vermelho', 'CV', '#cc0000', 7, true),
(uuid_generate_v4(), 'Primeiro Comando da Capital', 'PCC', '#ffff00', 6, true),
(uuid_generate_v4(), 'Terceiro Comando Puro', 'TCP', '#00ff00', 5, true),

-- Outras organizações
(uuid_generate_v4(), 'Área Neutra', 'NEUTRA', '#808080', 1, true),
(uuid_generate_v4(), 'Área em Disputa', 'DISPUTA', '#ff6600', 2, true);

-- Inserir regiões (Fortaleza como exemplo)
INSERT INTO regions (id, name, type, municipality, state, centroid, bounds, status) VALUES
-- Bairros de Fortaleza
(uuid_generate_v4(), 'Centro', 'neighborhood', 'Fortaleza', 'CE', 
 ST_SetSRID(ST_MakePoint(-38.5267, -3.7172), 4326),
 ST_SetSRID(ST_MakePolygon(ST_GeomFromText('LINESTRING(-38.5367 -3.7072, -38.5167 -3.7072, -38.5167 -3.7272, -38.5367 -3.7272, -38.5367 -3.7072)')), 4326),
 'published'),

(uuid_generate_v4(), 'Messejana', 'neighborhood', 'Fortaleza', 'CE',
 ST_SetSRID(ST_MakePoint(-38.4851, -3.8099), 4326),
 ST_SetSRID(ST_MakePolygon(ST_GeomFromText('LINESTRING(-38.4951 -3.7999, -38.4751 -3.7999, -38.4751 -3.8199, -38.4951 -3.8199, -38.4951 -3.7999)')), 4326),
 'published'),

(uuid_generate_v4(), 'Barra do Ceará', 'neighborhood', 'Fortaleza', 'CE',
 ST_SetSRID(ST_MakePoint(-38.5855, -3.6983), 4326),
 ST_SetSRID(ST_MakePolygon(ST_GeomFromText('LINESTRING(-38.5955 -3.6883, -38.5755 -3.6883, -38.5755 -3.7083, -38.5955 -3.7083, -38.5955 -3.6883)')), 4326),
 'published'),

(uuid_generate_v4(), 'Aldeota', 'neighborhood', 'Fortaleza', 'CE',
 ST_SetSRID(ST_MakePoint(-38.4923, -3.7389), 4326),
 ST_SetSRID(ST_MakePolygon(ST_GeomFromText('LINESTRING(-38.5023 -3.7289, -38.4823 -3.7289, -38.4823 -3.7489, -38.5023 -3.7489, -38.5023 -3.7289)')), 4326),
 'published'),

(uuid_generate_v4(), 'Grande Bom Jardim', 'neighborhood', 'Fortaleza', 'CE',
 ST_SetSRID(ST_MakePoint(-38.5782, -3.7798), 4326),
 ST_SetSRID(ST_MakePolygon(ST_GeomFromText('LINESTRING(-38.5882 -3.7698, -38.5682 -3.7698, -38.5682 -3.7898, -38.5882 -3.7898, -38.5882 -3.7698)')), 4326),
 'published'),

-- Outras cidades
(uuid_generate_v4(), 'Copacabana', 'neighborhood', 'Rio de Janeiro', 'RJ',
 ST_SetSRID(ST_MakePoint(-43.2075, -22.9068), 4326),
 ST_SetSRID(ST_MakePolygon(ST_GeomFromText('LINESTRING(-43.2175 -22.8968, -43.1975 -22.8968, -43.1975 -22.9168, -43.2175 -22.9168, -43.2175 -22.8968)')), 4326),
 'published'),

(uuid_generate_v4(), 'Vila Madalena', 'neighborhood', 'São Paulo', 'SP',
 ST_SetSRID(ST_MakePoint(-46.6927, -23.5626), 4326),
 ST_SetSRID(ST_MakePolygon(ST_GeomFromText('LINESTRING(-46.7027 -23.5526, -46.6827 -23.5526, -46.6827 -23.5726, -46.7027 -23.5726, -46.7027 -23.5526)')), 4326),
 'published');

-- Inserir usuários de exemplo
INSERT INTO users (id, name, email, password_hash, role, status) VALUES
-- Admin
(uuid_generate_v4(), 'Administrator', 'admin@mapfactions.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin', 'active'),

-- Moderadores
(uuid_generate_v4(), 'Moderador Ceará', 'mod-ce@mapfactions.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'moderator', 'active'),
(uuid_generate_v4(), 'Moderador SP', 'mod-sp@mapfactions.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'moderator', 'active'),

-- Colaboradores
(uuid_generate_v4(), 'Colaborador Fortaleza', 'colab-for@mapfactions.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'collaborator', 'active'),

-- Cidadãos
(uuid_generate_v4(), 'Cidadão Teste', 'cidadao@exemplo.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'citizen', 'active');

-- Inserir geometrias de exemplo (polígonos de áreas)
-- Note: Estas são coordenadas de exemplo e devem ser ajustadas com dados reais

-- Polígono de exemplo para Centro de Fortaleza
INSERT INTO geometries (
    id, region_id, geojson, geometry_type, risk_level, faction_id, 
    validity_start, source, notes, status, author_id
) 
SELECT 
    uuid_generate_v4(),
    r.id,
    '{"type":"Feature","geometry":{"type":"Polygon","coordinates":[[[-38.5367,-3.7072],[-38.5167,-3.7072],[-38.5167,-3.7272],[-38.5367,-3.7272],[-38.5367,-3.7072]]]},"properties":{"name":"Centro - Área TDN"}}',
    'polygon',
    'medium',
    f.id,
    NOW(),
    'Seed Data',
    'Área histórica de atuação do TDN no centro de Fortaleza',
    'published',
    u.id
FROM regions r, factions f, users u
WHERE r.name = 'Centro' 
    AND f.acronym = 'TDN' 
    AND u.role = 'admin'
LIMIT 1;

-- Mais exemplos de geometrias
INSERT INTO geometries (
    id, region_id, geojson, geometry_type, risk_level, faction_id, 
    validity_start, source, notes, status, author_id
) 
SELECT 
    uuid_generate_v4(),
    r.id,
    '{"type":"Feature","geometry":{"type":"Polygon","coordinates":[[[-38.4951,-3.7999],[-38.4751,-3.7999],[-38.4751,-3.8199],[-38.4951,-3.8199],[-38.4951,-3.7999]]]},"properties":{"name":"Messejana - Área CV"}}',
    'polygon',
    'high',
    f.id,
    NOW(),
    'Seed Data',
    'Área de forte atuação do CV em Messejana',
    'published',
    u.id
FROM regions r, factions f, users u
WHERE r.name = 'Messejana' 
    AND f.acronym = 'CV' 
    AND u.role = 'admin'
LIMIT 1;

INSERT INTO geometries (
    id, region_id, geojson, geometry_type, risk_level, faction_id, 
    validity_start, source, notes, status, author_id
) 
SELECT 
    uuid_generate_v4(),
    r.id,
    '{"type":"Feature","geometry":{"type":"Polygon","coordinates":[[[-38.5882,-3.7698],[-38.5682,-3.7698],[-38.5682,-3.7898],[-38.5882,-3.7898],[-38.5882,-3.7698]]]},"properties":{"name":"Bom Jardim - Área GDE"}}',
    'polygon',
    'high',
    f.id,
    NOW(),
    'Seed Data',
    'Território tradicionalmente controlado pelo GDE',
    'published',
    u.id
FROM regions r, factions f, users u
WHERE r.name = 'Grande Bom Jardim' 
    AND f.acronym = 'GDE' 
    AND u.role = 'admin'
LIMIT 1;

-- Área neutra (Aldeota - bairro nobre)
INSERT INTO geometries (
    id, region_id, geojson, geometry_type, risk_level, faction_id, 
    validity_start, source, notes, status, author_id
) 
SELECT 
    uuid_generate_v4(),
    r.id,
    '{"type":"Feature","geometry":{"type":"Polygon","coordinates":[[[-38.5023,-3.7289],[-38.4823,-3.7289],[-38.4823,-3.7489],[-38.5023,-3.7489],[-38.5023,-3.7289]]]},"properties":{"name":"Aldeota - Área Neutra"}}',
    'polygon',
    'low',
    f.id,
    NOW(),
    'Seed Data',
    'Bairro nobre com baixa incidência de atividade de facções',
    'published',
    u.id
FROM regions r, factions f, users u
WHERE r.name = 'Aldeota' 
    AND f.acronym = 'NEUTRA' 
    AND u.role = 'admin'
LIMIT 1;

-- Inserir alguns incidentes de exemplo
INSERT INTO incidents (id, type, description, geom, start_time, severity, status, author_id) 
SELECT 
    uuid_generate_v4(),
    'Bloqueio de Via',
    'Bloqueio na Av. Bezerra de Menezes',
    ST_SetSRID(ST_MakePoint(-38.5267, -3.7172), 4326),
    NOW() - INTERVAL '2 hours',
    3,
    'active',
    u.id
FROM users u WHERE u.role = 'admin' LIMIT 1;

INSERT INTO incidents (id, type, description, geom, start_time, end_time, severity, status, author_id) 
SELECT 
    uuid_generate_v4(),
    'Confronto',
    'Confronto entre facções rivais resolvido',
    ST_SetSRID(ST_MakePoint(-38.4851, -3.8099), 4326),
    NOW() - INTERVAL '1 day',
    NOW() - INTERVAL '20 hours',
    4,
    'resolved',
    u.id
FROM users u WHERE u.role = 'moderator' LIMIT 1;

-- Inserir assinaturas de exemplo
INSERT INTO subscriptions (id, user_id, plan, status, amount, start_date, end_date)
SELECT 
    uuid_generate_v4(),
    u.id,
    'monthly',
    'active',
    14.90,
    NOW(),
    NOW() + INTERVAL '1 month'
FROM users u WHERE u.email = 'cidadao@exemplo.com';

-- Logs de auditoria de exemplo
INSERT INTO audit_logs (id, user_id, entity, entity_id, action, ip_address, user_agent)
SELECT 
    uuid_generate_v4(),
    u.id,
    'geometry',
    g.id,
    'create',
    '127.0.0.1'::inet,
    'Mozilla/5.0 (Test Browser)'
FROM users u, geometries g 
WHERE u.role = 'admin' 
LIMIT 3;