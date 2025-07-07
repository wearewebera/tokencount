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

# Sign the binary
echo "🔏 Signing $BINARY_PATH..."
codesign --force \
  --options runtime \
  --entitlements "$ENTITLEMENTS" \
  --sign "$MACOS_CERTIFICATE_NAME" \
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