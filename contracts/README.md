# BulkTasks Contract

This directory contains the Solidity contracts and deployment scripts for the BulkTasks project.

## Prerequisites

- [Foundry](https://getfoundry.sh/) installed.
- Environment variables set for deployment (see Deployment section).

## Building

To compile the contracts, navigate to this `contracts` directory in your terminal and run:

```bash
forge build
```

This will compile the contracts and place the artifacts in the `artifacts` directory (or as configured in `foundry.toml`).

## Deployment

The `BulkTasks` contract is deployed using a Forge script (`script/DeployBulkTasks.s.sol`). This script reads configuration (Base Token and Engine addresses) from JSON files located in the root `config` directory based on the target environment.

### Prerequisites for Deployment

Before deploying, you **must** set the following environment variables in your terminal:

1.  `DEPLOY_ENV`: Specifies the deployment environment and determines which config file to use (`config.<ENV>.json`). Valid values typically include `local`, `testnet`, `sepolia`, `mainnet`.
2.  `RPC_URL`: The RPC endpoint URL for the target network.
3.  `PRIVATE_KEY`: The private key of the account that will deploy the contract. **Never commit this key.**
4.  `ETHERSCAN_API_KEY` (Optional, but required for verification): Your Etherscan API key if you intend to verify the contract on Etherscan (usually needed for `mainnet`).

**Example Environment Variable Setup:**

*   **Linux/macOS (Bash/Zsh):**
    ```bash
    export DEPLOY_ENV="sepolia"
    export RPC_URL="YOUR_SEPOLIA_RPC_URL"
    export PRIVATE_KEY="YOUR_DEPLOYER_PRIVATE_KEY"
    export ETHERSCAN_API_KEY="YOUR_ETHERSCAN_KEY"
    ```
*   **Windows (PowerShell):**
    ```powershell
    $env:DEPLOY_ENV = "sepolia"
    $env:RPC_URL = "YOUR_SEPOLIA_RPC_URL"
    $env:PRIVATE_KEY = "YOUR_DEPLOYER_PRIVATE_KEY"
    $env:ETHERSCAN_API_KEY = "YOUR_ETHERSCAN_KEY"
    ```

### Running the Deployment Script

Navigate to this `contracts` directory in your terminal.

*   **For non-mainnet deployment (e.g., Sepolia):**
    Set `DEPLOY_ENV=sepolia` (or your target testnet) and run:
    ```bash
    # Linux/macOS/Git Bash:
    forge script script/DeployBulkTasks.s.sol:DeployBulkTasks --rpc-url $RPC_URL --private-key $PRIVATE_KEY --broadcast
    
    # Windows PowerShell:
    forge script script/DeployBulkTasks.s.sol:DeployBulkTasks --rpc-url $env:RPC_URL --private-key $env:PRIVATE_KEY --broadcast
    
    # Windows Command Prompt:
    forge script script/DeployBulkTasks.s.sol:DeployBulkTasks --rpc-url %RPC_URL% --private-key %PRIVATE_KEY% --broadcast
    ```

*   **For mainnet deployment (with verification):**
    Set `DEPLOY_ENV=mainnet` and run:
    ```bash
    # Linux/macOS/Git Bash:
    forge script script/DeployBulkTasks.s.sol:DeployBulkTasks --rpc-url $RPC_URL --private-key $PRIVATE_KEY --broadcast --verify
    
    # Windows PowerShell:
    forge script script/DeployBulkTasks.s.sol:DeployBulkTasks --rpc-url $env:RPC_URL --private-key $env:PRIVATE_KEY --broadcast --verify
    
    # Windows Command Prompt:
    forge script script/DeployBulkTasks.s.sol:DeployBulkTasks --rpc-url %RPC_URL% --private-key %PRIVATE_KEY% --broadcast --verify
    ```
    *(Ensure `ETHERSCAN_API_KEY` is set.)*

The script will read the appropriate `config.<DEPLOY_ENV>.json` file, extract the necessary addresses, deploy the `BulkTasks` contract, and log the deployment address.
