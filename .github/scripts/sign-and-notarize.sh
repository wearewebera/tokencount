#!/bin/bash
set -e

# This script is used in GitHub Actions to sign and notarize macOS binaries
# It expects the certificate to be in base64 format in the environment

BINARY_PATH="$1"
ENTITLEMENTS="$2"

echo "🔐 Setting up keychain..."

# Create a temporary keychain
KEYCHAIN_PATH="$RUNNER_TEMP/app-signing.keychain-db"
KEYCHAIN_PASSWORD="$(openssl rand -base64 32)"

# Create keychain
security create-keychain -p "$KEYCHAIN_PASSWORD" "$KEYCHAIN_PATH"
security set-keychain-settings -lut 21600 "$KEYCHAIN_PATH"
security unlock-keychain -p "$KEYCHAIN_PASSWORD" "$KEYCHAIN_PATH"

# Import certificate
echo "$MACOS_CERTIFICATE" | base64 --decode > certificate.p12
security import certificate.p12 -P "$MACOS_CERTIFICATE_PWD" -A -t cert -f pkcs12 -k "$KEYCHAIN_PATH"
security list-keychain -d user -s "$KEYCHAIN_PATH"
rm certificate.p12

# Find the certificate identity
echo "🔍 Finding certificate identity..."
IDENTITY=$(security find-identity -v -p codesigning "$KEYCHAIN_PATH" | grep "Developer ID Application" | head -1 | awk -F'"' '{print $2}')

if [ -z "$IDENTITY" ]; then
  echo "⚠️  No 'Developer ID Application' certificate found"
  echo "📋 Available certificates:"
  security find-identity -v -p codesigning "$KEYCHAIN_PATH"
  
  # Try to find any Apple Developer certificate as fallback
  IDENTITY=$(security find-identity -v -p codesigning "$KEYCHAIN_PATH" | grep -E "Apple Development|Developer ID" | head -1 | awk -F'"' '{print $2}')
  
  if [ -z "$IDENTITY" ]; then
    echo "❌ No suitable certificate found for code signing"
    exit 1
  fi
  
  echo "⚠️  Using certificate: $IDENTITY"
  echo "⚠️  Note: For distribution outside the App Store, you need a 'Developer ID Application' certificate"
fi

echo "✅ Found identity: $IDENTITY"

# Sign the binary
echo "🔏 Signing $BINARY_PATH..."
codesign --force \
  --options runtime \
  --entitlements "$ENTITLEMENTS" \
  --sign "$IDENTITY" \
  --timestamp \
  --keychain "$KEYCHAIN_PATH" \
  "$BINARY_PATH"

echo "✅ Signed successfully"

# Verify
codesign --verify --deep --strict --verbose=2 "$BINARY_PATH"

# Notarize
echo "📤 Notarizing..."
ditto -c -k --keepParent "$BINARY_PATH" "$BINARY_PATH.zip"

xcrun notarytool submit "$BINARY_PATH.zip" \
  --apple-id "$APPLE_ID" \
  --password "$NOTARIZATION_PASSWORD" \
  --team-id "$TEAM_ID" \
  --wait

# Staple
echo "📌 Stapling..."
xcrun stapler staple "$BINARY_PATH"

# Cleanup
rm -f "$BINARY_PATH.zip"
security delete-keychain "$KEYCHAIN_PATH"

echo "✨ Done!"