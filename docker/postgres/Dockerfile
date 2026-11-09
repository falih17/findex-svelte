# Base image: PostgreSQL 16 with PostGIS 3.4
FROM postgis/postgis:16-3.4

# Metadata
LABEL maintainer="Findex Engineering <dev@findex.local>"
LABEL description="PostgreSQL 16 with PostGIS 3.4 and pgvector extension for SaaS Presensi"

# Install build dependencies for pgvector
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    git \
    postgresql-server-dev-16 \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Clone and compile pgvector from source (v0.8.0)
RUN git clone --branch v0.8.0 https://github.com/pgvector/pgvector.git /tmp/pgvector \
    && cd /tmp/pgvector \
    && make \
    && make install \
    && rm -rf /tmp/pgvector

# Clean build packages to reduce image size
RUN apt-get purge -y --auto-remove \
    build-essential \
    git \
    postgresql-server-dev-16 \
    && rm -rf /var/lib/apt/lists/*
