# macOS Code Signing Setup

This guide explains how to set up code signing and notarization for the tokencount binaries.

## Prerequisites

1. Apple Developer Account ($99/year)
2. Developer ID Application certificate
3. App-specific password for notarization

## Step 1: Create Developer ID Certificate

**IMPORTANT**: You need a "Developer ID Application" certificate for distributing apps outside the Mac App Store. Regular "Apple Development" certificates won't work for notarization.

1. Sign in to [Apple Developer](https://developer.apple.com)
2. Go to Certificates, Identifiers & Profiles
3. Click the + button to create a new certificate
4. **Select "Developer ID Application"** (NOT "Apple Development")
   - This is under the "Software" section
   - Required for apps distributed outside the Mac App Store
5. Follow the instructions to create a CSR using Keychain Access
6. Download and install the certificate

If you only have an "Apple Development" certificate:
- The signing will work but notarization will fail
- Users will still see Gatekeeper warnings
- You need to create a proper "Developer ID Application" certificate

## Step 2: Find Your Certificate Name

```bash
security find-identity -v -p codesigning
```

Look for something like:
```
"Developer ID Application: Your Name (TEAMID)"
```

## Step 3: Create App-Specific Password

1. Go to [appleid.apple.com](https://appleid.apple.com)
2. Sign in and go to "Sign-In and Security"
3. Select "App-Specific Passwords"
4. Generate a new password for "tokencount-notarization"
5. Save this password securely

## Step 4: Export Certificate for GitHub Actions

```bash
# Export your certificate to a .p12 file
security export -k ~/Library/Keychains/login.keychain-db \
  -t identities \
  -f pkcs12 \
  -o certificate.p12 \
  -P "your-export-password"

# Convert to base64 for GitHub secrets
base64 -i certificate.p12 -o certificate.base64

# Copy the base64 content
cat certificate.base64 | pbcopy

# Clean up
rm certificate.p12 certificate.base64
```

## Step 5: Configure GitHub Secrets

Go to your repository Settings > Secrets and variables > Actions, and add:

1. `MACOS_CERTIFICATE`: The base64 certificate content (from step 4)
2. `MACOS_CERTIFICATE_PWD`: The password you used when exporting (step 4)
3. `MACOS_CERTIFICATE_NAME`: The certificate name from step 2
4. `APPLE_ID`: Your Apple ID email
5. `NOTARIZATION_PASSWORD`: The app-specific password from step 3
6. `TEAM_ID`: Your Apple Developer Team ID (visible in developer portal)

## Local Signing

To sign locally before creating a release:

```bash
# Set environment variables
export DEVELOPER_ID="Developer ID Application: Your Name (TEAMID)"
export APPLE_ID="your-email@example.com"
export NOTARIZATION_PASSWORD="xxxx-xxxx-xxxx-xxxx"
export TEAM_ID="YOURTEAMID"

# Build the binary
go build -o tokencount

# Sign and notarize
./scripts/sign-macos.sh tokencount scripts/entitlements.plist
```

## Verification

After signing, verify the binary:

```bash
# Check signature
codesign -dv --verbose=4 tokencount

# Check notarization
spctl -a -vvv -t install tokencount
```

You should see:
```
tokencount: accepted
source=Notarized Developer ID
```

## Troubleshooting

### "Unable to build chain to self-signed root for signer"
- Make sure you have the Apple Developer certificates installed
- Download from: https://www.apple.com/certificateauthority/

### "The username or password was incorrect"
- Ensure you're using an app-specific password, not your Apple ID password
- Check that the Apple ID has accepted the latest developer agreements

### "Team ID not found"
- Find your Team ID in the Apple Developer portal under Membership