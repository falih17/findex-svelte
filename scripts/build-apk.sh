#!/usr/bin/env bash

set -e

echo "=========================================================="
echo "🚀 FINDEX ATTENDANCE - ANDROID APK BUILD SCRIPT"
echo "=========================================================="

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FRONTEND_DIR="$ROOT_DIR/services/frontend"
OUTPUT_DIR="$ROOT_DIR/dist"

MODE="${1:-debug}" # default to debug build (no signing keystore required)

echo "1. Checking directory..."
cd "$FRONTEND_DIR"

echo "2. Building SvelteKit Static SPA..."
npm run build

echo "3. Syncing assets with Capacitor Android..."
npx cap sync android

echo "4. Compiling Android APK ($MODE)..."
cd android

if [ "$MODE" = "release" ]; then
    ./gradlew assembleRelease
    APK_SOURCE="app/build/outputs/apk/release/app-release-unsigned.apk"
else
    ./gradlew assembleDebug
    APK_SOURCE="app/build/outputs/apk/debug/app-debug.apk"
fi

mkdir -p "$OUTPUT_DIR"
cp "$APK_SOURCE" "$OUTPUT_DIR/findex-attendance-$MODE.apk"

echo "=========================================================="
echo "✅ BUILD BERHASIL!"
echo "APK Siap diunduh di: $OUTPUT_DIR/findex-attendance-$MODE.apk"
echo "=========================================================="
