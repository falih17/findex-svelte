-- ==========================================================
-- 01: INITIALIZE POSTGRESQL EXTENSIONS
-- ==========================================================

-- UUID and Crypto for Secure Identifiers
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- PostGIS for Geolocation & Geofencing (Spatial queries)
CREATE EXTENSION IF NOT EXISTS "postgis";

-- pgvector for AI Face Biometrics Embeddings (512-D vectors)
CREATE EXTENSION IF NOT EXISTS "vector";

-- Verify extensions
DO $$
BEGIN
    RAISE NOTICE 'Extensions initialized successfully: uuid-ossp, pgcrypto, postgis, vector';
END $$;
