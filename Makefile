# Lumo Makefile

# Variables
BINARY_NAME=lumo
BUILD_DIR=build
VERSION=$(shell grep -oP 'Version = "\K[^"]+' pkg/version/version.go)
BUILD_DATE=$(shell date +%Y-%m-%d)
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GO_VERSION=$(shell go version | awk '{print $$3}')
LDFLAGS=-ldflags "-X github.com/agnath18/lumo/pkg/version.Version=${VERSION} -X github.com/agnath18/lumo/pkg/version.BuildDate=${BUILD_DATE} -X github.com/agnath18/lumo/pkg/version.GitCommit=${GIT_COMMIT} -X github.com/agnath18/lumo/pkg/version.GoVersion=${GO_VERSION}"
DEB_PACKAGE=${BUILD_DIR}/lumo_${VERSION}_amd64.deb

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	@echo "Building Lumo v${VERSION}..."
	@mkdir -p ${BUILD_DIR}
	@go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} cmd/lumo/main.go
	@echo "Build complete: ${BUILD_DIR}/${BINARY_NAME}"

# Install the binary
.PHONY: install
install: build
	@echo "Installing Lumo to /usr/bin/${BINARY_NAME}..."
	@sudo cp ${BUILD_DIR}/${BINARY_NAME} /usr/bin/${BINARY_NAME}
	@echo "Installation complete"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test ./...

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf ${BUILD_DIR}
	@echo "Clean complete"

# Clean Debian package files
.PHONY: clean-deb
clean-deb:
	@echo "Cleaning Debian package files..."
	@rm -rf debian
	@rm -f *.deb
	@rm -f ${BUILD_DIR}/*.deb
	@echo "Debian package files cleaned"

# Show version
.PHONY: version
version:
	@echo "Lumo v${VERSION}"
	@echo "Build date: ${BUILD_DATE}"
	@echo "Git commit: ${GIT_COMMIT}"
	@echo "Go version: ${GO_VERSION}"

# Build Debian package
.PHONY: deb
deb: build
	@echo "Building Debian package..."
	@./scripts/generate_control.sh
	@mkdir -p debian/usr/bin debian/usr/share/doc/lumo
	@cp ${BUILD_DIR}/${BINARY_NAME} debian/usr/bin/
	@cp README.md debian/usr/share/doc/lumo/
	@chmod 755 debian/usr/bin/${BINARY_NAME}
	@find debian -type d -exec chmod 755 {} \;
	@mkdir -p ${BUILD_DIR}
	@dpkg-deb --build debian ${DEB_PACKAGE}
	@echo "Debian package built: ${DEB_PACKAGE}"

# Help
.PHONY: help
help:
	@echo "Lumo Makefile"
	@echo ""
	@echo "Targets:"
	@echo "  all      - Build the binary (default)"
	@echo "  build    - Build the binary in the ${BUILD_DIR} directory"
	@echo "  install  - Install the binary to /usr/bin"
	@echo "  test     - Run tests"
	@echo "  clean    - Clean build artifacts from ${BUILD_DIR}"
	@echo "  deb      - Build a Debian package (.deb)"
	@echo "  clean-deb - Clean Debian package files"
	@echo "  version  - Show version information"
	@echo "  help     - Show this help"
