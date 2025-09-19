-- Initialization script for PostgreSQL with PostGIS
-- This script is run when the container starts for the first time

-- Create the main database if it doesn't exist
SELECT 'CREATE DATABASE map_factions'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'map_factions')\gexec

-- Connect to the database
\c map_factions;

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS postgis_topology;
CREATE EXTENSION IF NOT EXISTS postgis_raster;
CREATE EXTENSION IF NOT EXISTS fuzzystrmatch;
CREATE EXTENSION IF NOT EXISTS postgis_tiger_geocoder;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS btree_gist;

-- Set timezone
SET timezone = 'UTC';

-- Optimize PostgreSQL for geospatial workloads
ALTER SYSTEM SET shared_preload_libraries = 'pg_stat_statements';
ALTER SYSTEM SET max_connections = 200;
ALTER SYSTEM SET shared_buffers = '256MB';
ALTER SYSTEM SET effective_cache_size = '1GB';
ALTER SYSTEM SET maintenance_work_mem = '64MB';
ALTER SYSTEM SET checkpoint_completion_target = 0.9;
ALTER SYSTEM SET wal_buffers = '16MB';
ALTER SYSTEM SET default_statistics_target = 100;
ALTER SYSTEM SET random_page_cost = 1.1;
ALTER SYSTEM SET effective_io_concurrency = 200;

-- Enable query statistics
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

-- Create a read-only user for monitoring
CREATE USER monitoring_user WITH PASSWORD 'monitoring_password';
GRANT CONNECT ON DATABASE map_factions TO monitoring_user;
GRANT USAGE ON SCHEMA public TO monitoring_user;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO monitoring_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO monitoring_user;

COMMENT ON DATABASE map_factions IS 'Database for Map Area Factions application with PostGIS extensions enabled';