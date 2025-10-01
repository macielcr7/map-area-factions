-- Seed data for Map Area Factions
-- Phase 1 development data

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enable PostGIS extension
CREATE EXTENSION IF NOT EXISTS postgis;

-- Insert default admin user (password: admin123)
INSERT INTO users (id, name, email, password_hash, role, active, created_at, updated_at) 
VALUES (
    uuid_generate_v4(),
    'Administrator',
    'admin@mapfactions.com',
    '$2a$10$LLwd44QbRJejYZj/lKMXruwJLFlx3Y3k7QLEyRVGw3b8QJMCsYNNG',
    'admin',
    true,
    NOW(),
    NOW()
) ON CONFLICT (email) DO NOTHING;

-- Insert moderator user (password: moderator123)
INSERT INTO users (id, name, email, password_hash, role, active, created_at, updated_at) 
VALUES (
    uuid_generate_v4(),
    'Moderador',
    'moderator@mapfactions.com',
    '$2a$10$bX6MuXGGiunC11NvA1LRGeFn5KjR1ve5/jtHMX5Unz5G5JivA7mwS',
    'moderator',
    true,
    NOW(),
    NOW()
) ON CONFLICT (email) DO NOTHING;

-- Insert collaborator user (password: collaborator123) 
INSERT INTO users (id, name, email, password_hash, role, active, created_at, updated_at) 
VALUES (
    uuid_generate_v4(),
    'Colaborador',
    'collaborator@mapfactions.com',
    '$2a$10$Cp4jPQ8thaSeB0p54blnvefUgbJWYdD2cvGr1LJnAK5a0te0XHnSy',
    'collaborator',
    true,
    NOW(),
    NOW()
) ON CONFLICT (email) DO NOTHING;

-- Insert main Brazilian factions
INSERT INTO factions (id, name, acronym, color_hex, display_priority, active, created_at, updated_at)
VALUES 
    (uuid_generate_v4(), 'Terceiro Comando do Nordeste', 'TDN', '#ff0000', 10, true, NOW(), NOW()),
    (uuid_generate_v4(), 'Comando Vermelho', 'CV', '#cc0000', 9, true, NOW(), NOW()),
    (uuid_generate_v4(), 'Guardiões do Estado', 'GDE', '#0000ff', 8, true, NOW(), NOW()),
    (uuid_generate_v4(), 'Primeiro Comando da Capital', 'PCC', '#ffff00', 7, true, NOW(), NOW()),
    (uuid_generate_v4(), 'Família do Norte', 'FDN', '#008000', 6, true, NOW(), NOW()),
    (uuid_generate_v4(), 'Comando Classe A', 'CCA', '#ff8000', 5, true, NOW(), NOW()),
    (uuid_generate_v4(), 'Milícia', 'MIL', '#800080', 4, true, NOW(), NOW()),
    (uuid_generate_v4(), 'Neutro/Disputa', 'NEU', '#808080', 1, true, NOW(), NOW())
ON CONFLICT (acronym) DO NOTHING;

-- Insert some example regions (Fortaleza neighborhoods)
INSERT INTO regions (id, name, type, municipality, state, status, created_at, updated_at)
VALUES 
    (uuid_generate_v4(), 'Messejana', 'neighborhood', 'Fortaleza', 'CE', 'published', NOW(), NOW()),
    (uuid_generate_v4(), 'Centro', 'neighborhood', 'Fortaleza', 'CE', 'published', NOW(), NOW()),
    (uuid_generate_v4(), 'Aldeota', 'neighborhood', 'Fortaleza', 'CE', 'published', NOW(), NOW()),
    (uuid_generate_v4(), 'Meireles', 'neighborhood', 'Fortaleza', 'CE', 'published', NOW(), NOW()),
    (uuid_generate_v4(), 'Barra do Ceará', 'neighborhood', 'Fortaleza', 'CE', 'published', NOW(), NOW()),
    (uuid_generate_v4(), 'Montese', 'neighborhood', 'Fortaleza', 'CE', 'published', NOW(), NOW()),
    (uuid_generate_v4(), 'Parangaba', 'neighborhood', 'Fortaleza', 'CE', 'published', NOW(), NOW()),
    (uuid_generate_v4(), 'Antônio Bezerra', 'neighborhood', 'Fortaleza', 'CE', 'published', NOW(), NOW()),
    (uuid_generate_v4(), 'Passaré', 'neighborhood', 'Fortaleza', 'CE', 'published', NOW(), NOW()),
    (uuid_generate_v4(), 'Jangurussu', 'neighborhood', 'Fortaleza', 'CE', 'published', NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Insert system settings
INSERT INTO system_settings (id, data, created_at, updated_at)
VALUES (
    uuid_generate_v4(),
    '{"app_name": "Map Area Factions", "version": "1.0.0", "maintenance_mode": false, "max_geometries_per_user": 100, "max_file_size_mb": 10}',
    NOW(),
    NOW()
) ON CONFLICT DO NOTHING;

-- Insert example incidents
INSERT INTO incidents (id, type, description, location, latitude, longitude, geom, severity, status, author_id, created_at, updated_at)
VALUES 
    (uuid_generate_v4(), 'Violence', 'Relato de violência na região', 'Messejana, Fortaleza', -3.7974, -38.5432, ST_SetSRID(ST_MakePoint(-38.5432, -3.7974), 4326), 3, 'active', (SELECT id FROM users WHERE email = 'admin@mapfactions.com'), NOW(), NOW()),
    (uuid_generate_v4(), 'Drug Trafficking', 'Tráfico de drogas identificado', 'Centro, Fortaleza', -3.7172, -38.5432, ST_SetSRID(ST_MakePoint(-38.5432, -3.7172), 4326), 4, 'active', (SELECT id FROM users WHERE email = 'moderator@mapfactions.com'), NOW(), NOW()),
    (uuid_generate_v4(), 'Theft', 'Roubo reportado', 'Aldeota, Fortaleza', -3.7172, -38.5432, ST_SetSRID(ST_MakePoint(-38.5432, -3.7172), 4326), 2, 'resolved', (SELECT id FROM users WHERE email = 'collaborator@mapfactions.com'), NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Insert example subscriptions
INSERT INTO subscriptions (id, user_id, plan, status, start_date, end_date, price, currency, payment_method, auto_renew, created_at, updated_at)
VALUES 
    (uuid_generate_v4(), (SELECT id FROM users WHERE email = 'admin@mapfactions.com'), 'lifetime', 'active', NOW(), NOW() + INTERVAL '1 year', 299.90, 'BRL', 'credit_card', false, NOW(), NOW()),
    (uuid_generate_v4(), (SELECT id FROM users WHERE email = 'moderator@mapfactions.com'), 'monthly', 'active', NOW(), NOW() + INTERVAL '1 month', 29.90, 'BRL', 'pix', true, NOW(), NOW()),
    (uuid_generate_v4(), (SELECT id FROM users WHERE email = 'collaborator@mapfactions.com'), 'quarterly', 'pending', NOW(), NOW() + INTERVAL '3 months', 79.90, 'BRL', 'boleto', false, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Insert example reports
INSERT INTO reports (id, user_id, geometry_id, text, status, created_at, updated_at)
VALUES 
    (uuid_generate_v4(), (SELECT id FROM users WHERE email = 'collaborator@mapfactions.com'), NULL, 'Relato de segurança na região', 'open', NOW(), NOW()),
    (uuid_generate_v4(), (SELECT id FROM users WHERE email = 'moderator@mapfactions.com'), NULL, 'Verificação de precisão dos dados', 'in_review', NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Insert example audit logs
INSERT INTO audit_logs (id, user_id, entity, entity_id, action, diff, ip_address, user_agent, created_at)
VALUES 
    (uuid_generate_v4(), (SELECT id FROM users WHERE email = 'admin@mapfactions.com'), 'user', (SELECT id FROM users WHERE email = 'moderator@mapfactions.com'), 'create', '{"name": "Moderador", "email": "moderator@mapfactions.com"}', '192.168.1.1', 'Mozilla/5.0', NOW()),
    (uuid_generate_v4(), (SELECT id FROM users WHERE email = 'admin@mapfactions.com'), 'faction', (SELECT id FROM factions WHERE acronym = 'TDN'), 'create', '{"name": "Terceiro Comando do Nordeste", "acronym": "TDN"}', '192.168.1.1', 'Mozilla/5.0', NOW())
ON CONFLICT DO NOTHING;

-- Note: Geometry data will be inserted through the API since we need proper GeoJSON handling
-- Example geometries can be created via POST /api/v1/geometries with proper authentication