#!/bin/bash
# ProxyDoctor Setup Script
# Installs dependencies, runs tests, and builds binaries

set -e

echo "🔧 ProxyDoctor Setup"
echo "===================="
echo ""

# Check Go version
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo "✅ Found Go: $GO_VERSION"
echo ""

# Download and verify dependencies
echo "📦 Downloading dependencies..."
go mod download
go mod verify
echo "✅ Dependencies verified"
echo ""

# Run tests
echo "🧪 Running tests..."
if go test -v ./...; then
    echo "✅ All tests passed!"
else
    echo "❌ Tests failed. Please fix before proceeding."
    exit 1
fi
echo ""

# Build binaries
echo "🏗️  Building binaries..."
mkdir -p bin

if go build -o bin/proxydoctor ./cmd/cli; then
    echo "✅ CLI built: bin/proxydoctor"
else
    echo "❌ Failed to build CLI"
    exit 1
fi

if go build -o bin/proxydoctor-server ./cmd/server; then
    echo "✅ Server built: bin/proxydoctor-server"
else
    echo "❌ Failed to build server"
    exit 1
fi
echo ""

# Success
echo "🎉 Setup complete!"
echo ""
echo "Next steps:"
echo "  • Run CLI:     ./run.sh cli --help"
echo "  • Run server:  ./run.sh server"
echo "  • Run tests:   go test -v ./..."
echo ""
