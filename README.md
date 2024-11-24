# Gobius

A high-performance, multi-GPU mining client for the [Arbius protocol](https://arbius.ai/), written in Go.

Gobius enables decentralized machine learning by allowing miners to participate in the Arbius network - an on-chain, reproducible AI model execution platform similar to services like Midjourney or ChatGPT, but in a censorship-resistant way.

## Features

- **Multi-GPU Support**
  - Scales from single GPU to 100+ GPUs
  - Battle-tested with large GPU clusters
  - Efficient resource management

- **Cross-Platform**
  - Single binary deployment
  - Supports Windows, Linux and macOS
  - Built with Go for high performance

- **Smart Contract Integration**
  - Integrates with Arbitrum One network
  - Handles task queuing and solution submission
  - Supports bulk operations for gas optimization

- **Advanced Logging & Monitoring**
  - Terminal UI with real-time logging
  - Dynamic color support

- **Model Support**
  - Supports multiple AI models
  - Integrates with IPFS for model/data storage
  - Uses Cog for standardized model interfaces

- **Development Tools**
  - Built-in debugging support
  - Configuration via JSON
  - Testnet support for development

## Development Setup

Windows, Linux and macOS are supported

### Prerequisites

- Go 1.20 or later
- Git

### Installation

1. Clone the repository with submodules:
   ```bash
   # Clone with submodules
   git clone --recursive https://github.com/<tbc>/gobius
   cd gobius

   # Or if you've already cloned the repository:
   git submodule update --init --recursive
   ```

2. Run the setup script:
   ```bash
   # Linux/MacOS
   ./scripts/setup.sh

   # Windows
   scripts\setup.bat
   ```

   This will install required tools in the `./bin` directory:
   - solc (Solidity Compiler)
   - abigen (Ethereum bindings generator)

### Git Submodules

This project uses Git submodules for managing external dependencies. The following submodules are included:

```
external/
└── arbius/          # Arbius v4+ contracts
    └── contract/      # Smart contract source code, needed to create go bindings
```

#### Working with Submodules

1. **Initial clone with submodules**:
   ```bash
   git clone --recursive https://github.com/your/project
   ```

2. **Update submodules to their latest commits**:
   ```bash
   git submodule update --remote
   ```

3. **Initialize submodules (if cloned without --recursive)**:
   ```bash
   git submodule update --init --recursive
   ```

4. **Check submodule status**:
   ```bash
   git submodule status
   ```

### Configuration

The following environment variables can be used to customize the installation:

| Variable | Description | Default |
|----------|-------------|---------|
| `SOLC_VERSION` | Solidity compiler version | 0.8.19 |
| `ABIGEN_VERSION` | Abigen/go-ethereum version | 1.13.15 |

Example:
```bash
# Install specific versions
SOLC_VERSION=0.8.19 ABIGEN_VERSION=1.13.4 ./scripts/setup.sh
```

### Directory Structure

```
project/
├── bin/              # Compiled tools (solc, abigen)
├── contracts/        # Solidity smart contracts
├── external/         # Git submodules
│   └── arbius/       # External contract dependencies
├── scripts/         
│   ├── setup.sh     # Unix setup script
│   └── setup.bat    # Windows setup script
└── tools/
    ├── solc/        # Solc installer
    └── abigen/      # Abigen installer
```

### Troubleshooting

1. **Wrong directory**: Make sure to run the setup script from the project root directory.

2. **Permission denied**: On Unix systems, you might need to make the script executable:
   ```bash
   chmod +x scripts/setup.sh
   ```

3. **Go version issues**: If you encounter Go-related errors, ensure you're using Go 1.20 or later:
   ```bash
   go version
   ```

4. **Path issues**: The tools are installed in the `./bin` directory. Add this to your path or use the full path when executing:
   ```bash
   ./bin/solc --version
   ./bin/abigen --version
   ```

5. **Submodule issues**: If you see empty external directories or get submodule errors:
   ```bash
   # Reset and update submodules
   git submodule sync
   git submodule update --init --recursive --force
   ```

### Development Workflow

1. Write your Solidity contracts in the `contracts/` directory.

2. Generate Go bindings:
   ```bash
   go generate ./...
   ```

3. Build your project:
   ```bash
   go build
   ```

### Updating Dependencies

1. **Update Go dependencies**:
   ```bash
   go get -u ./...
   go mod tidy
   ```

2. **Update git submodules**:
   ```bash
   git submodule update --remote
   git add external/
   git commit -m "Update external dependencies"
   ```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. By contributing to this project, you agree to license your work under the terms of the MIT License.
