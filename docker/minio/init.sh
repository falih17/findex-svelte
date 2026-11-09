#!/bin/sh
# Wait for MinIO to start
until /usr/bin/mc alias set myminio http://minio:9000 minioadmin minioadminsecret; do
    echo "...Waiting for MinIO to become ready..."
    sleep 2
done

echo "MinIO is ready. Creating buckets..."

# Create buckets if they don't exist
/usr/bin/mc mb --ignore-existing myminio/findex-biometrics
/usr/bin/mc mb --ignore-existing myminio/findex-attendances
/usr/bin/mc mb --ignore-existing myminio/findex-documents

# Set anonymous download policies for public image serving (if desired for attendance proof)
/usr/bin/mc anonymous set download myminio/findex-biometrics
/usr/bin/mc anonymous set download myminio/findex-attendances

echo "MinIO buckets created successfully!"
exit 0
