# Base image: PostgreSQL 16 with pgvector on Debian 12 Bookworm (Multi-arch: ARM64 & AMD64)
FROM pgvector/pgvector:pg16

LABEL maintainer="Findex Engineering <dev@findex.local>"
LABEL description="PostgreSQL 16 with pgvector and PostGIS for SaaS Presensi Findex"

# Install PostGIS 3 for PostgreSQL 16 directly from official PGDG repository
RUN apt-get update && apt-get install -y --no-install-recommends \
    postgresql-16-postgis-3 \
    postgresql-16-postgis-3-scripts \
    && rm -rf /var/lib/apt/lists/*
