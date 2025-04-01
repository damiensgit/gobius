# Gobius

![gobius logo](/logo.jpg)

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

## Setup

Windows, Linux and macOS are supported

### Prerequisites

- [Go 1.22 or later](https://go.dev/dl/)
- Git

- **IPFS**: Gobius requires IPFS (InterPlanetary File System) for model and data storage.
  - Install IPFS: Follow the official installation guide at [https://docs.ipfs.tech/install/command-line/#install-official-binary-distributions](https://docs.ipfs.tech/install/command-line/#install-official-binary-distributions).
  - **Initialization**: After installing, initialize your IPFS node:
    ```bash
    ipfs init
    ```
  - **Daemon**: Run the IPFS daemon:
    ```bash
    ipfs daemon
    ```
  - **Firewall**: IPFS uses specific ports (typically 4001 TCP/UDP, 5001 TCP, 8080 TCP). Ensure these are open in your firewall, especially if running on a cloud provider or restricted network. Some hosting providers may block peer-to-peer traffic, which could affect IPFS operation. Consult your provider's documentation and configure firewall rules accordingly.

### System Setup

1. Install Go:
   ```bash
   # Ubuntu/Debian
   wget https://go.dev/dl/go1.22.9.linux-amd64.tar.gz
   sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.20.14.linux-amd64.tar.gz
   echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
   source ~/.bashrc

   # macOS (using Homebrew)
   brew install go

   # windows
   Download and install the latest windows installer from https://go.dev/dl/

   # Verify installation
   go version
   ```

### Building the Miner

1. Clone and build:
   ```bash
   # Clone the repository
   git clone --recursive https://github.com/damiensgit/gobius.git
   cd gobius

   # Build the miner (this will automatically download required packages)
   go build
   ```

### Building Contract Bindings

The project uses [Task](https://taskfile.dev) to manage contract binding generation. First, install Task:

```bash
go install github.com/go-task/task/v3/cmd/task@latest
```

1. Generate all contract bindings:
   ```bash
   task build:all
   ```

2. Or generate specific bindings:
   ```bash
   # Build just the voter interface
   task build:voter

   # Build just the engine contract
   task build:enginev5

   # Build bulk tasks contract
   task build:bulktasks
   
   # Build SQL schema bindings
   task build:sqlc
   ```

The bindings will be generated in the `./bindings` directory with appropriate subdirectories for each contract.

### SQL Schema Requirements

To build SQL bindings, you'll need [sqlc](https://sqlc.dev/) installed:

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

The SQL schema and configuration are located in the `./sql` directory. Running `task build:sqlc` will generate Go code from the SQL schema.

### Configuration Setup

1. Copy the example configuration:
   ```bash
   cp config.example.json config.json
   ```

2. Edit config.json with your settings:
   ```json
   {
     "rpc": "YOUR_RPC_ENDPOINT",  // e.g., "https://arb1.arbitrum.io/rpc"
     "privatekey": "YOUR_PRIVATE_KEY",  // Your wallet private key (without 0x prefix)
     "ipfs": {
       "strategy": "http_client",
       "http_client": {
         "url": "/ip4/127.0.0.1/tcp/5001"  // IPFS daemon address
       }
     }
     // ... other optional settings ...
   }
   ```

   > ⚠️ **Important**: Never share or commit your private key. Keep your config.json file secure.

   > **Note**: Ensure your IPFS daemon is running and accessible at the configured URL. If you're running IPFS on a different host or port, adjust the URL accordingly.

### Deployed Contracts

The miner makes use of a simple solidity contract that enables the submission of bulk claims and bulk commitments.

See ```/contracts``` folder for details.

There is a deployed contract on Arbitrum one available at ``0x75879250b1d43F8860Bd30C628E8606782a02a87``, but you may deploy your own e.g. for event tracking.

Use remix to deploy [this contract](contracts/BulkTasksMainnet.sol) and update the ./config/config.json field ``bulkTasksAddress`` with the new contract address.

### Running the Miner

Basic usage:
```bash
./gobius --config config.json
```

#### Command Line Flags

| Flag | Description | Default | Required |
|------|-------------|---------|----------|
| `--config` | Path to the configuration file | "config.json" | No |
| `--skipvalidation` | Skip safety checks and validation of the model and miner version | false | No |
| `--loglevel` | Set logging verbosity (1 = default) | 1 | No |
| `--testnet` | Run using testnet (1 = local, 2 = nova testnet) | 0 | No |
| `--taskscanner` | Scan blocks for unsolved tasks in pst 12 hours | 0 | No |

#### Example Commands

```bash
# Run with custom config file
./miner --config myconfig.json

# Run on testnet with increased logging
./miner --testnet 2 --loglevel 2

# Run task scanner
./miner --taskscanner 1
```

## For Developers

1. Clone the repository with submodules:
   ```bash
   # Clone with submodules
   git clone --recursive https://github.com/damiensgit/gobius.git
   cd gobius

   # Or if you've already cloned the repository:
   git submodule update --init --recursive
   ```

2. Run the setup script:
   ```bash
   # Linux/MacOS
   ./setup.sh

   # Windows
   setup.bat
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

3. **Go version issues**: If you encounter Go-related errors, ensure you're using Go 1.22 or later:
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
