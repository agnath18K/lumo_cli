#!/bin/bash
# Script to generate Debian control file

# Ensure the DEBIAN directory exists
mkdir -p debian/DEBIAN

# Get version from version.go
VERSION=$(grep -oP 'Version = "\K[^"]+' pkg/version/version.go)

# Create control file
cat > debian/DEBIAN/control << EOF
Package: lumo
Version: ${VERSION}
Section: utils
Priority: optional
Architecture: amd64
Depends: libc6 (>= 2.17)
Maintainer: agnath18K <agnath18@gmail.com>
Description: AI-powered CLI assistant
 Lumo is an AI-powered CLI assistant in Go that interprets
 natural language to execute commands. It helps users find
 relevant terminal commands and provides guidance for
 terminal tasks. Lumo integrates with Gemini, OpenAI, and
 Ollama APIs.
EOF

echo "Generated control file with version ${VERSION}"

# Create postinst script
cat > debian/DEBIAN/postinst << EOF
#!/bin/sh
# postinst script for lumo

set -e

# Make sure the binary is executable
chmod 755 /usr/bin/lumo

# Create log directory if it doesn't exist
mkdir -p /var/log/lumo
chmod 755 /var/log/lumo

exit 0
EOF

# Make postinst executable
chmod +x debian/DEBIAN/postinst

# Create prerm script
cat > debian/DEBIAN/prerm << EOF
#!/bin/sh
# prerm script for lumo

set -e

# Clean up any temporary files if needed
# (none for now)

exit 0
EOF

# Make prerm executable
chmod +x debian/DEBIAN/prerm

echo "Generated installation scripts"
