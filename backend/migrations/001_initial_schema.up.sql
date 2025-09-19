-- Migration: 001_initial_schema.up.sql
-- Criação do schema inicial com extensões PostGIS

-- Habilitar extensões PostGIS
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS postgis_topology;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enum types
CREATE TYPE user_role AS ENUM ('admin', 'moderator', 'collaborator', 'citizen');
CREATE TYPE subscription_status AS ENUM ('active', 'expired', 'cancelled', 'pending');
CREATE TYPE geometry_status AS ENUM ('draft', 'review', 'published');
CREATE TYPE risk_level AS ENUM ('low', 'medium', 'high');
CREATE TYPE incident_status AS ENUM ('active', 'resolved', 'expired');
CREATE TYPE report_status AS ENUM ('open', 'in_review', 'resolved', 'rejected');

-- Tabela de usuários
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role user_role DEFAULT 'citizen',
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabela de facções
CREATE TABLE factions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    acronym VARCHAR(10) NOT NULL UNIQUE,
    color_hex VARCHAR(7) NOT NULL, -- #RRGGBB
    display_priority INTEGER DEFAULT 0,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabela de regiões (municípios, bairros, etc.)
CREATE TABLE regions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'neighborhood', 'zone', 'custom'
    municipality VARCHAR(255) NOT NULL,
    state VARCHAR(2) NOT NULL,
    centroid GEOMETRY(POINT, 4326),
    bounds GEOMETRY(POLYGON, 4326),
    status geometry_status DEFAULT 'draft',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabela de geometrias (polígonos e polilinhas)
CREATE TABLE geometries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    region_id UUID REFERENCES regions(id) ON DELETE CASCADE,
    geojson JSONB NOT NULL,
    geometry_type VARCHAR(20) NOT NULL, -- 'polygon', 'polyline'
    risk_level risk_level DEFAULT 'low',
    faction_id UUID REFERENCES factions(id),
    validity_start TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    validity_end TIMESTAMP WITH TIME ZONE,
    source VARCHAR(255),
    notes TEXT,
    status geometry_status DEFAULT 'draft',
    geom GEOMETRY NOT NULL, -- Generated from geojson
    author_id UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabela de versões de geometrias (histórico)
CREATE TABLE geometry_versions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    geometry_id UUID REFERENCES geometries(id) ON DELETE CASCADE,
    geojson JSONB NOT NULL,
    diff JSONB,
    author_id UUID REFERENCES users(id),
    action VARCHAR(50) NOT NULL, -- 'create', 'update', 'delete'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabela de incidentes
CREATE TABLE incidents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type VARCHAR(100) NOT NULL,
    description TEXT,
    geom GEOMETRY NOT NULL, -- Point or Polygon
    start_time TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    end_time TIMESTAMP WITH TIME ZONE,
    severity INTEGER DEFAULT 1, -- 1-5
    status incident_status DEFAULT 'active',
    author_id UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabela de assinaturas
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    plan VARCHAR(50) NOT NULL, -- 'monthly', 'quarterly', 'lifetime'
    status subscription_status DEFAULT 'pending',
    provider VARCHAR(50) DEFAULT 'mercado_pago',
    external_id VARCHAR(255),
    start_date TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    end_date TIMESTAMP WITH TIME ZONE,
    amount DECIMAL(10,2),
    invoice_data JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabela de logs de auditoria
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    entity VARCHAR(100) NOT NULL,
    entity_id UUID NOT NULL,
    action VARCHAR(50) NOT NULL,
    diff JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabela de relatórios/denúncias
CREATE TABLE reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    geometry_id UUID REFERENCES geometries(id),
    incident_id UUID REFERENCES incidents(id),
    text TEXT NOT NULL,
    attachments JSONB, -- URLs dos anexos no S3
    status report_status DEFAULT 'open',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Índices espaciais (GIST)
CREATE INDEX idx_regions_centroid ON regions USING GIST (centroid);
CREATE INDEX idx_regions_bounds ON regions USING GIST (bounds);
CREATE INDEX idx_geometries_geom ON geometries USING GIST (geom);
CREATE INDEX idx_incidents_geom ON incidents USING GIST (geom);

-- Índices regulares
CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_role ON users (role);
CREATE INDEX idx_factions_active ON factions (active);
CREATE INDEX idx_geometries_status ON geometries (status);
CREATE INDEX idx_geometries_faction_id ON geometries (faction_id);
CREATE INDEX idx_geometries_region_id ON geometries (region_id);
CREATE INDEX idx_incidents_status ON incidents (status);
CREATE INDEX idx_incidents_type ON incidents (type);
CREATE INDEX idx_subscriptions_user_id ON subscriptions (user_id);
CREATE INDEX idx_subscriptions_status ON subscriptions (status);
CREATE INDEX idx_audit_logs_entity ON audit_logs (entity, entity_id);
CREATE INDEX idx_reports_status ON reports (status);

-- Trigger para atualizar updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Aplicar triggers
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_factions_updated_at BEFORE UPDATE ON factions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_regions_updated_at BEFORE UPDATE ON regions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_geometries_updated_at BEFORE UPDATE ON geometries FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_incidents_updated_at BEFORE UPDATE ON incidents FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_subscriptions_updated_at BEFORE UPDATE ON subscriptions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_reports_updated_at BEFORE UPDATE ON reports FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function para converter GeoJSON em Geometry
CREATE OR REPLACE FUNCTION update_geometry_from_geojson()
RETURNS TRIGGER AS $$
BEGIN
    NEW.geom = ST_SetSRID(ST_GeomFromGeoJSON(NEW.geojson->>'geometry'), 4326);
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger para atualizar geometry a partir do geojson
CREATE TRIGGER update_geometries_geom 
    BEFORE INSERT OR UPDATE ON geometries 
    FOR EACH ROW EXECUTE FUNCTION update_geometry_from_geojson();

-- Views para consultas otimizadas
CREATE VIEW v_active_geometries AS
SELECT 
    g.*,
    f.name as faction_name,
    f.acronym as faction_acronym,
    f.color_hex,
    r.name as region_name,
    r.municipality,
    r.state
FROM geometries g
LEFT JOIN factions f ON g.faction_id = f.id
LEFT JOIN regions r ON g.region_id = r.id
WHERE g.status = 'published' 
    AND (g.validity_end IS NULL OR g.validity_end > NOW())
    AND (f.active IS NULL OR f.active = true);

-- View para estatísticas por município
CREATE VIEW v_municipality_stats AS
SELECT 
    r.municipality,
    r.state,
    COUNT(g.id) as total_geometries,
    COUNT(CASE WHEN g.risk_level = 'high' THEN 1 END) as high_risk_areas,
    COUNT(CASE WHEN g.risk_level = 'medium' THEN 1 END) as medium_risk_areas,
    COUNT(CASE WHEN g.risk_level = 'low' THEN 1 END) as low_risk_areas,
    COUNT(DISTINCT g.faction_id) as total_factions
FROM regions r
LEFT JOIN geometries g ON r.id = g.region_id 
    AND g.status = 'published'
    AND (g.validity_end IS NULL OR g.validity_end > NOW())
GROUP BY r.municipality, r.state;