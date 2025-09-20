-- Seed data for Map Area Factions
-- Phase 1 development data

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enable PostGIS extension
CREATE EXTENSION IF NOT EXISTS postgis;

-- Insert default admin user (password: admin123)
INSERT INTO users (id, name, email, password_hash, role, status, created_at, updated_at) 
VALUES (
    uuid_generate_v4(),
    'Administrator',
    'admin@mapfactions.com',
    '$2a$10$8/8M8Yn6fJQK1NQ9YZ2.x.zJrFV2Q1ZvCBGzxFCXr1zxpKMzW3iX2',
    'admin',
    'active',
    NOW(),
    NOW()
) ON CONFLICT (email) DO NOTHING;

-- Insert moderator user (password: moderator123)
INSERT INTO users (id, name, email, password_hash, role, status, created_at, updated_at) 
VALUES (
    uuid_generate_v4(),
    'Moderador',
    'moderator@mapfactions.com',
    '$2a$10$9/9M9Yn7fJQK2NQ9YZ3.x.zJrFV3Q2ZvCBGzxFCXr2zxpKMzW4iX3',
    'moderator',
    'active',
    NOW(),
    NOW()
) ON CONFLICT (email) DO NOTHING;

-- Insert collaborator user (password: collaborator123) 
INSERT INTO users (id, name, email, password_hash, role, status, created_at, updated_at) 
VALUES (
    uuid_generate_v4(),
    'Colaborador',
    'collaborator@mapfactions.com',
    '$2a$10$a/aMapaYn8fJQK3NQ9YZ4.x.zJrFV4Q3ZvCBGzxFCXr3zxpKMzW5iX4',
    'collaborator',
    'active',
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

-- Note: Geometry data will be inserted through the API since we need proper GeoJSON handling
-- Example geometries can be created via POST /api/v1/geometries with proper authentication