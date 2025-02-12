#!/bin/bash
set -e

echo "Installing development tools..."

# Install solc
echo "Installing solc..."
go run tools/solc/install.go

# Install abigen
echo "Installing abigen..."
go run tools/abigen/install.go

echo "Setup complete! Tools installed in ./bin directory"