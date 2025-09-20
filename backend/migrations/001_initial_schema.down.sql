-- Migration: 001_initial_schema.down.sql
-- Rollback do schema inicial

-- Drop views
DROP VIEW IF EXISTS v_municipality_stats;
DROP VIEW IF EXISTS v_active_geometries;

-- Drop triggers
DROP TRIGGER IF EXISTS update_geometries_geom ON geometries;
DROP TRIGGER IF EXISTS update_reports_updated_at ON reports;
DROP TRIGGER IF EXISTS update_subscriptions_updated_at ON subscriptions;
DROP TRIGGER IF EXISTS update_incidents_updated_at ON incidents;
DROP TRIGGER IF EXISTS update_geometries_updated_at ON geometries;
DROP TRIGGER IF EXISTS update_regions_updated_at ON regions;
DROP TRIGGER IF EXISTS update_factions_updated_at ON factions;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop functions
DROP FUNCTION IF EXISTS update_geometry_from_geojson();
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables (in reverse order due to foreign keys)
DROP TABLE IF EXISTS reports;
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS incidents;
DROP TABLE IF EXISTS geometry_versions;
DROP TABLE IF EXISTS geometries;
DROP TABLE IF EXISTS regions;
DROP TABLE IF EXISTS factions;
DROP TABLE IF EXISTS users;

-- Drop enum types
DROP TYPE IF EXISTS report_status;
DROP TYPE IF EXISTS incident_status;
DROP TYPE IF EXISTS risk_level;
DROP TYPE IF EXISTS geometry_status;
DROP TYPE IF EXISTS subscription_status;
DROP TYPE IF EXISTS user_role;