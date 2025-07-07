#!/bin/bash
set -e

# This script signs and notarizes macOS binaries
# Prerequisites:
# 1. Valid Apple Developer ID Application certificate in Keychain
# 2. App-specific password for notarization
# 3. xcrun notarytool configured

if [ $# -lt 1 ]; then
    echo "Usage: $0 <binary-path> [entitlements-file]"
    exit 1
fi

BINARY_PATH="$1"
ENTITLEMENTS="${2:-}"

# Your Apple Developer ID - you'll need to update this
DEVELOPER_ID="Developer ID Application: Your Name (TEAMID)"
BUNDLE_ID="com.wearewebera.tokencount"

echo "🔏 Signing $BINARY_PATH..."

if [ -n "$ENTITLEMENTS" ]; then
    codesign --force --options runtime --entitlements "$ENTITLEMENTS" --sign "$DEVELOPER_ID" --timestamp "$BINARY_PATH"
else
    codesign --force --options runtime --sign "$DEVELOPER_ID" --timestamp "$BINARY_PATH"
fi

echo "✅ Signed successfully"

# Verify the signature
echo "🔍 Verifying signature..."
codesign --verify --deep --strict --verbose=2 "$BINARY_PATH"

# Create a zip for notarization
ZIP_PATH="${BINARY_PATH}.zip"
echo "📦 Creating zip for notarization..."
ditto -c -k --keepParent "$BINARY_PATH" "$ZIP_PATH"

# Notarize
echo "📤 Submitting for notarization..."
xcrun notarytool submit "$ZIP_PATH" \
    --apple-id "$APPLE_ID" \
    --password "$NOTARIZATION_PASSWORD" \
    --team-id "$TEAM_ID" \
    --wait

echo "🎉 Notarization complete!"

# Staple the notarization
echo "📌 Stapling notarization..."
xcrun stapler staple "$BINARY_PATH"

# Clean up
rm -f "$ZIP_PATH"

echo "✨ Done! Binary is signed and notarized."