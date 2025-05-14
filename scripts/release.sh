#!/bin/bash
# Script to help with releasing new versions of Lumo

# Check if a version is provided
if [ $# -ne 1 ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 1.0.1"
    exit 1
fi


NEW_VERSION=$1

# Update version in version.go
sed -i "s/Version = \"[0-9.]*\"/Version = \"$NEW_VERSION\"/" pkg/version/version.go

echo "Updated version to $NEW_VERSION in pkg/version/version.go"

# Build the Debian package
make clean
make clean-deb
make deb

echo "Release process completed for version $NEW_VERSION"
echo "Debian package created: build/lumo_${NEW_VERSION}_amd64.deb"
echo ""
echo "Next steps:"
echo "1. Test the package: sudo dpkg -i build/lumo_${NEW_VERSION}_amd64.deb"
echo "2. Commit the changes: git commit -am \"Release version $NEW_VERSION\""
echo "3. Tag the release: git tag -a v$NEW_VERSION -m \"Version $NEW_VERSION\""
echo "4. Push the changes: git push && git push --tags"
