# Quickstart Guide for Gobius Mining

**Important Note:** While Gobius runs on Windows, Linux, and macOS, this guide primarily assumes you are setting up the miner on a **Linux server**, preferably a recent Ubuntu release hosted by a cheap cloud provider (e.g., Vultr, Hetzner, etc.).

Running the miner on a server ensures it can operate 24/7. To keep Gobius running even after you disconnect your SSH session, you should run it inside a terminal multiplexer like `screen` or `tmux`. Here's a very basic guide:

1.  **Connect to Your Server via SSH**:
    Open a terminal on your local machine and use the `ssh` command:
    ```bash
    ssh your_username@your_server_ip_address
    ```
    Replace `your_username` and `your_server_ip_address` accordingly. You might need to provide a password or use an SSH key.

2.  **Install `screen` or `tmux` (if needed)**:
    Most Linux distributions come with `screen`. `tmux` might need installation:
    ```bash
    # For Ubuntu/Debian
    sudo apt update
    sudo apt install tmux
    ```

3.  **Start a Session**:
    -   **Using `screen`**:
        ```bash
        screen
        ```
        Press Enter or Space when prompted.
    -   **Using `tmux`**:
        ```bash
        tmux new-session -s gobius
        ```

4.  **Run Gobius Inside the Session**: Navigate to your Gobius directory (e.g., `cd gobius`) and start the miner as described later in this guide (`./gobius --config config.json`).

5.  **Detach from the Session**: This leaves Gobius running in the background.
    -   **`screen`**: Press `Ctrl+a` then `d`.
    -   **`tmux`**: Press `Ctrl+b` then `d`.

6.  **Reattach Later**: To check on your miner:
    -   **`screen`**: `screen -r`
    -   **`tmux`**: `tmux attach-session -t gobius`

This guide will now walk you through setting up Gobius for mining on the Arbius network. Follow these steps to get started quickly.


## Prerequisites

1. **Go**: Ensure you have Go 1.22 or later installed. Follow the installation instructions for your operating system:
   - **Ubuntu/Debian**:
     ```bash
     wget https://go.dev/dl/go1.22.9.linux-amd64.tar.gz
     sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.22.9.linux-amd64.tar.gz
     echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
     source ~/.bashrc
     ```
   - **macOS** (using Homebrew):
     ```bash
     brew install go
     ```
   - **Windows**: Download and install the latest Go installer from [https://go.dev/dl/](https://go.dev/dl/).

2. **Git**: Install Git if you haven't already. You can download it from [https://git-scm.com/](https://git-scm.com/).

3. **IPFS**: Gobius requires IPFS for model and data storage.

   **Recommendation**: It's highly recommended to run the IPFS daemon on the **same server** where you are running Gobius. This simplifies network configuration.

   - Install IPFS by following the official guide at [https://docs.ipfs.tech/install/command-line/#install-official-binary-distributions](https://docs.ipfs.tech/install/command-line/#install-official-binary-distributions).
   - Initialize your IPFS node using the `server` profile. This helps prevent potential network issues by disabling local discovery mechanisms that might be flagged by some hosting providers:
     ```bash
     ipfs init --profile server
     ```
   - Run the IPFS daemon (ideally within its own `screen` or `tmux` session):
     ```bash
     ipfs daemon
     ```
     **Identifying the API Address**: When the daemon starts, look for lines in the output similar to:
     `RPC API server listening on /ip4/127.0.0.1/tcp/5001`
     This tells you the address the daemon is using. The default `/ip4/127.0.0.1/tcp/5001` is what you'll typically use in `config.json` if running IPFS locally.
   - **Firewall**: If running IPFS on the same machine as Gobius, typically no extra firewall rules are needed for Gobius to connect to IPFS. IPFS itself requires port 4001 (TCP/UDP) open for external peering, but the local connection from Gobius uses port 5001.
   - **Advanced Setup**: If you choose to run IPFS on a separate server, you must configure your firewalls to allow Gobius to connect to the IPFS daemon's API port (usually 5001 TCP) on the remote server. You will also need to update the `ipfs.http_client.url` in your `config.json` accordingly (e.g., `/ip4/<remote-ipfs-ip>/tcp/5001`).

## Building Gobius

1. Clone the Gobius repository:
   ```bash
   git clone --recursive https://github.com/damiensgit/gobius.git
   cd gobius
   ```

2. Build the miner:
   ```bash
   go build
   ```
   **Note**: The first time you run `go build`, Go will automatically download all the necessary package dependencies. This might take a few minutes depending on your internet connection. Subsequent builds will be much faster.

## Configuration

1. Copy the example configuration:
   ```bash
   cp config.example.json config.json
   ```

2. Edit `config.json` with your settings:

   **Understanding the Required Accounts**

   Gobius uses a few different Ethereum accounts (EOAs) for distinct purposes. Using separate accounts helps with security, clarity, and potentially managing transaction nonces (transaction counters) during high-volume operations:

   1.  **Main Operational Account (`blockchain.private_key`)**: 
        -   **Location**: The field in the `blockchain` section: `"blockchain": { "private_key": "YOUR_KEY" }`.
        -   **Purpose**: This account is used to:
            -   Pay gas fees (in ETH) for all transactions.
            -   Pay for validator stake top-ups. When Gobius detects a validator needs more AIUS for stake, it **uses this account** to send the AIUS to the validator. **This means your main account must hold enough AIUS** to satisfy the minimum stake requirements for all your validators.
            -   If no specific bulk task accounts are specified (see below), this account will also pay for task submission fees.

   2.  **Validator Account(s) (`validator_config.private_keys`)**:
        -   **Location**: The `private_keys` array within the `validator_config` section.
        -   **Purpose**: These accounts:
            -   Hold the required **AIUS stake** for mining.
            -   Sign solutions/claims generated by the miner.
            -   Receive AIUS mining rewards.
            -   **Do NOT pay for their own stake top-ups** (that comes from the main account).
        -   You can configure **multiple validator accounts** if you are running a larger operation.

   3.  **Bulk Task Submission Account(s) (`batchtasks.private_keys`)**: 
        -   **Purpose**: Primarily for managing transaction nonces during high-volume operations.
        -   **Configuration**: The `config.json` includes a dedicated `batchtasks` section containing a `private_keys` array.
        -   **Behavior**: 
            -   **If Configured**: If you add one or more private keys to the `batchtasks.private_keys` array, Gobius will use **only these accounts** for:
                -   Sending bulk task/claim transactions. 
                -   Paying the task submission fees (the AIUS fee specified in `automine.fee`).
                -   Claiming completed solutions.
            -   **If Empty (Default)**: If the `batchtasks.private_keys` array is left **empty** (`[]`), Gobius will **default** to using the **Main Operational Account** for all of the above operations.
        -   **Recommendation**: For most users, the default behavior is recommended.

   **Recommendation**: For security and clarity, it's strongly advised to use **separate, dedicated wallets** for the Main Operational role and each Validator role, rather than reusing a single personal wallet for everything.

   **Exporting Your Private Key (Security Warning!)**

   To configure Gobius, you need the private key for both your main operational wallet (the `privatekey` field) and your validator wallet(s) (the `validator_config.private_keys` field). Exporting private keys carries **significant security risks**. Anyone who obtains your private key has **full control** over your wallet and its assets. **Never share your private key or store it insecurely.** Only proceed if you understand the risks.

   **How to Export from MetaMask (Example):**

   1. Open MetaMask and select the account you want to use (either the main operational account or a validator account).
   2. Click the three vertical dots (`⋮`) next to the account name.
   3. Select "Account details".
   4. Click the "Show private key" button (or similar wording, like "Export Private Key").
   5. Enter your MetaMask password to confirm.
   6. Your private key will be displayed. **Carefully copy this key.**
   7. **Crucially**: Paste the key into your `config.json` file **without the leading `0x`**. For example, if MetaMask shows `0xabc123...`, you put only `abc123...` in the configuration file.

   **Other Wallets**: The process for other EVM wallets (like Trust Wallet, Rabby Wallet, etc.) is generally similar, usually found within account settings or security options. Consult your specific wallet's official documentation for precise steps.

   **Recommendation**: For enhanced security, consider creating **new, separate wallets** specifically for use with Gobius (one for operations/fees, and one or more for validators) rather than using your primary wallet containing significant funds.

   **Editing `config.json` - Important Note:**

   The JSON snippet below shows the **most common fields** you will need to edit in your `config.json` file. Your actual `config.json` (copied from `config.example.json`) contains **many more settings**. You should **edit the fields shown below within your full `config.json` file**, leaving the other fields at their default values unless you understand their purpose. **Do NOT simply copy and paste this entire snippet** - it is only a partial example highlighting key configuration points.

   ```json
   {
     "db_path": "storage.db",
     "privatekey": "YOUR_MAIN_OPERATIONAL_PRIVATE_KEY_NO_0X",  // Key for sending tx, paying gas (ETH)
     "ipfs": {
       "strategy": "http_client",
       "http_client": {
         "url": "/ip4/127.0.0.1/tcp/5001"  // Default for local IPFS. Replace 127.0.0.1 if remote.
       }
     },
     "blockchain": {
       "private_key": "YOUR_MAIN_OPERATIONAL_PRIVATE_KEY_NO_0X",  // Same key as top-level privatekey
       "rpc_url": "wss://arbitrum-one-rpc.publicnode.com",  // WebSocket Arbitrum RPC URL
       "sender_rpc_url": "",  // Optional separate RPC for transaction sending
       "client_rpc_urls": [],  // Optional additional RPCs for redundancy
       "cache_nonce": false,
       "basefee_x": 2  // Multiplier for base fee estimation
     },
     "ml": {
       "strategy": "cog",
       "cog": {
         "0x89c39001e3b23d2092bd998b62f07b523d23deb55e1627048b4ed47a4a38d5cc": {
           "url": [
             "<cog-url>"  // Replace with your Cog model URL
           ]
         }
       }
     },
     "strategies": {
       "model": "0x89c39001e3b23d2092bd998b62f07b523d23deb55e1627048b4ed47a4a38d5cc",
       "strategy": "automine",
       "automine": {
         "owner": "0x1234567890123456789012345678901234567890",
         "version": 0,
         "model": "0x89c39001e3b23d2092bd998b62f07b523d23deb55e1627048b4ed47a4a38d5cc",
         "fee": 7000000000000000,
         "input": {
           "prompt": "What is the capital of the moon?"
         }
       }
     },
     "validator_config": {
       "private_keys": ["YOUR_VALIDATOR_PRIVATE_KEY_NO_0X"], // Key(s) for validator(s)
       "stake_check": true,
       "stake_check_interval": "120s", // How often to check stake/health
       "min_basetoken_threshold": 10, // Min AIUS to keep on validator
       "stake_buffer_percent": 2, // Target stake buffer % above min stake
       "stake_buffer_topup_percent": 1, // Top-up trigger % above min stake
       "initial_stake": 0 // Set > 0 to manually manage stake amount (disables auto top-up)
       // ... other advanced settings ...
     },
     "batchtasks": {
       "private_keys": []
     }
   }
   ```

   > **Note on Settings**: While `config.json` contains many settings, most can be left at their default values (copied from `config.example.json`) for initial setup. The defaults are designed to work well for basic auto-mining. Many options are for advanced users and fine-tuning specific behaviors. For details on all available options and their effects, advanced users can refer to the configuration loading logic in the source code (primarily within the `config/` directory).

   > **Database (`db_path`)**: Gobius uses a SQLite database (default: `storage.db`) to keep track of its state (pending commitments, solutions, claims, etc.). This allows Gobius to stop and restart safely, resuming where it left off. **Do not delete this file** unless you intend to reset the miner's state.

   > ⚠️ **Important**: Never share or commit your private key. Keep your `config.json` file secure.

   > **Note**: Ensure your IPFS daemon is running and accessible at the configured URL (`ipfs.http_client.url`). If you're running IPFS on a different host or port, adjust the URL accordingly. The default `/ip4/127.0.0.1/tcp/5001` assumes IPFS is running locally.

   > **Cog Model URLs**: Replace `<cog-url>` with the actual URL of the Cog model you are using. You can add multiple URLs if needed.

   **Understanding the `blockchain` Section**:

   The `blockchain` section is crucial for configuring how Gobius connects to the Arbitrum One network and manages transaction sending:

   -   **`private_key`**: This should match the top-level `privatekey` field and contains the private key for your main operational account (without the `0x` prefix). This account pays gas fees and sends transactions.
   -   **`rpc_url`**: The RPC endpoint Gobius uses to connect to the Arbitrum One network. This **must be a WebSocket URL** (starting with `wss://`), not an HTTP URL. Reliable options include:
        -   Public WebSocket endpoints: `wss://arbitrum-one-rpc.publicnode.com`
        -   Providers like QuickNode, Infura, Alchemy, etc. (requires account setup)
   -   **`sender_rpc_url`** (optional): A separate RPC endpoint specifically for sending transactions. Leave empty to use the main `rpc_url` for everything.
   -   **`client_rpc_urls`** (optional): An array of backup RPC endpoints for redundancy. Only necessary for high-availability setups.
   -   **`basefee_x`**: A multiplier applied to base fee estimates for gas pricing. The default of 2 is suitable for most users.

   > **Note**: Stable, reliable RPC connections are essential for Gobius to work properly. If using public endpoints, you might experience occasional connection issues. For serious mining operations, consider a paid RPC service like QuickNode for better reliability.

   **Recommended Validator AIUS Balance**:

   Your validator wallet (the address corresponding to the `private_keys` in `validator_config`) needs sufficient AIUS tokens to function correctly.

   -   **Minimum Stake**: The primary requirement is holding the network's minimum stake. As of writing, this is **420,926,548,100,086,163,465 wei** (approximately **421 AIUS**). You can always check the current minimum stake directly on the Arbius network or through community resources.
   -   **Buffer for Fees and Top-ups**: It is highly recommended to keep an **additional buffer** of AIUS in the validator wallet, beyond the minimum stake. A buffer of **at least 10%** (around 42 AIUS or more) is suggested.
   -   **Why the Buffer?**
        -   **Task Fees**: If using `automine`, each generated task requires a fee (`automine.fee`), which depletes the validator's balance.
        -   **Gas Costs**: All transactions (submitting solutions, claiming rewards, topping up stake) require gas fees, paid in **ETH (on Arbitrum)**.
        -   **Stake Top-ups**: Gobius might automatically top up the stake if it falls slightly (due to penalties or configuration). The buffer ensures funds are available for this.
   -   **`min_basetoken_threshold`**: This setting in `validator_config` helps enforce a minimum reserve, but you should still ensure the *total* balance is sufficient for both the stake and the buffer.

   **Understanding `validator_config`**:

   This section controls the behavior of your validator account(s). A validator is an Ethereum wallet (Externally Owned Account - EOA) that holds the required AIUS stake (minimum stake) and is used by Gobius to submit solutions and claim rewards.

   ```json
   "validator_config": {
     "private_keys": ["YOUR_VALIDATOR_PRIVATE_KEY_NO_0X"], // Key(s) for validator(s)
     "stake_check": true,
     "stake_check_interval": "120s", // How often to check stake/health
     "min_basetoken_threshold": 10, // Min AIUS to keep on validator
     "stake_buffer_percent": 2, // Target stake buffer % above min stake
     "stake_buffer_topup_percent": 1, // Top-up trigger % above min stake
     "initial_stake": 0 // Set > 0 to manually manage stake amount (disables auto top-up)
     // ... other advanced settings ...
   }
   ```

   -   **`private_keys`**: This is crucial. Add the private key(s) for your validator wallet(s) here as strings within the square brackets `[]`. You can add multiple keys if running multiple validators. **Remember to remove the leading `0x` from the private key.**
   -   **Stake Management (Defaults)**: By default (`initial_stake: 0`), Gobius automatically manages your stake based on percentages:
        -   It tries to maintain a stake level that is `stake_buffer_percent` above the network's minimum requirement.
        -   If the stake drops to only `stake_buffer_topup_percent` above the minimum, it will automatically top up the stake back to the `stake_buffer_percent` level.
   -   **Stake Management (Manual - `initial_stake`)**: If you set `initial_stake` to a value greater than 0 (representing an amount of AIUS), Gobius will attempt to top up the validator's stake to this specific amount. **Important:** This disables the automatic percentage-based top-ups, meaning you are responsible for ensuring the stake doesn't fall below the minimum required by the network.
   -   **`min_basetoken_threshold`**: Sets a minimum amount of AIUS tokens to always keep in the validator wallet. This reserve is used for future transaction fees (like stake top-ups or task submission fees) ensuring the validator doesn't run out of funds for essential operations.
   -   **`stake_check_interval`**: Determines how frequently Gobius checks the validator's stake level and performs health checks (default is every 120 seconds).

   **Important Note on Solution Rate Limit**: The Arbius network currently enforces a rate limit on solution submissions: **each validator account can only successfully submit 1 solution per second**. This means that to submit a batch containing *N* solutions, that specific validator must have been inactive (not submitted any solutions) for at least *N* seconds prior to the batch submission. If you attempt to submit a batch too quickly after a previous submission, the transaction may fail due to this rate limit. This is a key consideration when configuring batch sizes and observing miner behavior.

## Deployed Contract Addresses

Gobius interacts with several smart contracts deployed on the Arbitrum One network. While these are generally configured correctly in the base configuration, here are the key addresses for reference:

-   **Base Token (AIUS)**: `0x4a24B101728e07A52053c13FB4dB2BcF490CAbc3`
-   **Engine**: `0x9b51Ef044d3486A1fB0A2D55A6e0CeeAdd323E66`
-   **Bulk Tasks**: `0x75879250b1d43F8860Bd30C628E8606782a02a87`
-   **Voter**: `0x80E9B3dA81258705eC7C3DC89a799b78f2c68968`
-   **VeStaking**: `0x1c0a14cAC52ebDe9c724e5627162a90A26B85E15`
-   **Arbius Router**: `0x6caF23DC6dAe5EdCc53e4F330B61c420e6c71565`

These addresses are typically loaded automatically and do not require changes for basic operation. They are provided for reference only.

The `Bulk Tasks` contract is part of the Gobius project and has been deployed to mainnet for convenience. While using the provided address is recommended, advanced users can deploy their own version of the contract (source available in the `/contracts` folder). If you deploy your own, you **must** update the `bulkTasksAddress` field in your `config/config.json` file with the new contract address.

## Updating Gobius

When new versions of Gobius are released, you'll want to update your local repository and rebuild the miner.

**Before Updating**: It is **highly recommended** to back up your `config.json` file AND your database file (`storage.db` by default, check your `db_path` setting). Updates might include database schema changes (migrations).

1.  **Navigate to the Gobius Directory**:
    Make sure you are in the `gobius` directory on your server:
    ```bash
    cd path/to/gobius
    ```

2.  **Pull the Latest Changes**:
    Use `git pull` to download the latest code from the main repository:
    ```bash
    git pull origin main
    ```
    (Or replace `main` with the specific branch name if you are tracking a different one).

3.  **Rebuild the Miner**:
    Run the build command again:
    ```bash
    go build
    ```
    This will compile the updated code into a new `./gobius` executable.

**Configuration File Safety**:

-   Your `config.json` file will **not** be overwritten by the `git pull` command, as long as you created it by copying `config.example.json` and did not directly edit the example file itself.
-   Your database file (`storage.db` or custom `db_path`) will also **not** be overwritten by `git pull`.
-   **Important**: Always check the release notes or update announcements for the new version. Sometimes, updates might introduce changes to the required configuration format or add new options. Updates may also perform automatic database migrations upon first run - backing up the database is crucial in case of migration issues.

**Editing the Configuration File (`config.json`)**:

-   For quick edits on the server, the `nano` text editor is a simple option.
    -   Check if it's installed: `nano --version`
    -   Install if needed (Ubuntu/Debian): `sudo apt update && sudo apt install nano`
    -   Edit the file: `nano config.json`
    -   Use `Ctrl+O` to save (Write Out) and `Ctrl+X` to exit.
-   Advanced users can use their preferred terminal editor (like `vim`, `emacs`) or transfer the file to edit locally and then upload it back.

## Additional Resources

- For more detailed information, refer to the [README.md](README.md).
- If you encounter issues, check the [Troubleshooting](#troubleshooting-common-issues) section in this guide.

## Cog Setup

To use Cog models, you need to specify the URLs of the models in your configuration.

1. Locate the `ml` section in your `config.json`:
   ```json
   "ml": {
     "strategy": "cog",
     "cog": {
       "<model-id>": {
         "url": [
           "<cog-url>"  // Replace with your Cog model URL
         ]
       }
     }
   }
   ```

2. Replace `<model-id>` with the ID of your model (`0x89c39001e3b23d2092bd998b62f07b523d23deb55e1627048b4ed47a4a38d5cc`) and `<cog-url>` with the actual URL. You can add multiple URLs if needed.

## Arbitrum RPC Setup

Gobius requires a WebSocket or IPC connection to the Arbitrum network.

1. Update the `blockchain.rpc_url` field in your `config.json`:
   ```json
   "blockchain": {
     "rpc_url": "wss://arbitrum-one-rpc.publicnode.com"  // Example WebSocket URL
   }
   ```

2. You can use public WebSocket URLs or a provider like QuickNode for a more reliable connection.
   - **Public WebSocket**: `wss://arbitrum-one-rpc.publicnode.com`
   - **QuickNode**: Sign up at [QuickNode](https://www.quicknode.com/) and get your WebSocket RPC URL.

> **Note**: Ensure your RPC connection is stable and secure for optimal performance.

## Supported AI Model

For the initial mainnet launch, Gobius supports one primary AI model:

- **Model Name**: `qwen-qwq-32b`
- **Model ID**: `0x89c39001e3b23d2092bd998b62f07b523d23deb55e1627048b4ed47a4a38d5cc`

This Model ID is used in the `config.json` file when specifying the Cog model URL. Ensure you are configuring your GPU instance and miner settings for this specific model.

## Setting Up a GPU Instance on RunPod and Vast.ai

To run the Cog model `r8.im/kasumi-1/qwen-qwq-32b`, you need to set up a GPU instance with the following specifications:

**Important**: You **must** use an NVIDIA A100 GPU with 80GB of VRAM. No other GPU type or VRAM configuration is supported for this model.

**Note on Multiple GPUs**: While you might rent instances with multiple GPUs, the current Cog model (`r8.im/kasumi-1/qwen-qwq-32b`) will only utilize **one** GPU for processing tasks. Additional GPUs on the instance will not be used by the model.

### RunPod Setup

1. **Create an Account**: Sign up or log in to your RunPod account.

2. **Launch a Pod**:
   - Choose a GPU instance with an **NVIDIA A100 80GB** GPU.
   - Set the container disk size to **60GB**.
   - Set the volume disk size to **100GB**.
   - Ensure port **5000** is exposed for the application.

3. **Deploy the Model**:
   - Use the image `r8.im/kasumi-1/qwen-qwq-32b` to deploy your model.

### Vast.ai Setup

1. **Create an Account**: Sign up or log in to your Vast.ai account.

2. **Rent a GPU**: 
   - Select a machine with an **NVIDIA A100 80GB** GPU.
   - Configure the container disk to **60GB** and volume disk to **100GB**.
   - Expose port **5000** for external access.

3. **Deploy the Model**:
   - Use the image `r8.im/kasumi-1/qwen-qwq-32b` to deploy your model.

   > **Initialization Time**: Note that it can take several minutes for the Vast.ai instance to fully initialize, download the image, and start the Cog service after you click "Rent". Be patient and wait for the status to indicate it's running.

**Advanced (CLI)**:

If you prefer using the [Vast.ai CLI](https://vast.ai/docs/cli/getting-started), you can create a compatible instance using a specific offer ID. Find suitable offer IDs using the following command (filtering for NVIDIA A100s with 80GB VRAM and sufficient disk space):

```bash
vastai search offers 'gpu_name in ["A100_SXM4", "A100_PCI"] gpu_ram>=80 disk_space>100'
```

Replace `<OFFER_ID>` with your desired ID:

```bash
vastai create instance <OFFER_ID> --image r8.im/kasumi-1/qwen-qwq-32b:latest --env '-p 5000:5000 --gpus=all' --disk 100 --ssh --direct
```

This command sets up the instance with the required image, exposes port 5000, allocates disk space, and enables SSH access.

> **Note**: Ensure your instance meets the model's resource requirements for optimal performance.

### Extracting the URL for Launched Instances

Once your GPU instance is running (either on RunPod or Vast.ai), you need to find its public URL to configure Gobius.

#### RunPod

1.  **Access the Pod**: Navigate to "My Pods" in the RunPod dashboard and select your running pod.
2.  **Find the URL**: Go to the "Connect" tab. Look for the HTTP/HTTPS connection URL provided for port 5000. It will typically look like `https://<pod-id>-5000.proxy.runpod.net`.
3.  **Update Config**: Place this full URL, **appending `/predictions`**, in your `config.json` under the Cog model section:
    ```json
    "ml": {
      "strategy": "cog",
      "cog": {
        "0x89c39001e3b23d2092bd998b62f07b523d23deb55e1627048b4ed47a4a38d5cc": { // Ensure this is the correct model ID
          "url": [
            "https://ab122cde0ty69a-5000.proxy.runpod.net/predictions" // Example RunPod URL, ending in /predictions
          ]
        }
      }
    }
    ```

#### Vast.ai

1. Go to your Vast.ai Instances dashboard.
2. Find your running instance.
3. **Check Ports and IP**: Look for the section displaying connection information, specifically the IP address and port mappings. You need the **public IP address** and the **external port** that maps to the internal container port **5000/tcp**. It will often look something like `YOUR_PUBLIC_IP:EXTERNAL_PORT -> 5000/tcp`.
4. **Construct the URL**: Combine the public IP and the *external* port number to form the base URL (e.g., `http://YOUR_PUBLIC_IP:EXTERNAL_PORT`).
5. **Append Endpoint**: Add `/predictions` to the end.
   Example: `http://123.45.67.89:12345/predictions`
6. Update the `url` field in your `config.json` under `ml.cog.0x...` with this complete URL.

> **Note**: Ensure the URLs are accessible, include `http://` or `https://` as appropriate, end with the correct endpoint (`/predictions`), and are correctly formatted to avoid connectivity issues.

## Initial Mining Setup: Auto-Mining

The recommended initial setup and default strategy for Gobius is **auto-mining**. This means your miner will automatically generate its own tasks using the specified model and then solve them. This is the simplest way to get started.

This is configured in the `strategies` section of your `config.json` file. The example configuration (`config.example.json`) is already set up for this default behavior:

```json
"strategies": {
  "model": "0x89c39001e3b23d2092bd998b62f07b523d23deb55e1627048b4ed47a4a38d5cc",
  "strategy": "automine",
  "automine": {
    "owner": "0x1234567890123456789012345678901234567890",
    "version": 0,
    "model": "0x89c39001e3b23d2092bd998b62f07b523d23deb55e1627048b4ed47a4a38d5cc",
    "fee": 7000000000000000,
    "input": {
      "prompt": "What is the capital of the moon?"
    }
  }
},
"batchtasks": {
  "enabled": true,
  "private_keys": []
}
```

**Key Settings for Auto-Mining:**

-   `strategies.strategy`: Must be set to `"automine"`.
-   `strategies.automine.owner`: **Required**. Set this to your Ethereum wallet address (must include the `0x` prefix). You receive rewards for tasks you create.
-   `model` (in both `strategies` and `automine`): Ensure this matches the supported Model ID.
-   `automine.fee`: Fee paid per task submission (in AIUS wei). This depletes validator AIUS balance.
-   `automine.input.prompt`: Customizable prompt for generated tasks.

**Interaction with `batchtasks`**: 

-   For `automine` to successfully submit the tasks it generates, the `batchtasks` section must be configured and **enabled** (`"batchtasks": { "enabled": true }`). This allows Gobius to batch and submit the task commitments created by `automine`.
-   **Disabling Batch Tasks**: You might temporarily set `batchtasks.enabled` to `false` if you want to stop submitting *new* tasks (even if `automine` is configured) and focus on clearing the backlog of existing commitments/solutions and claims. This can be useful if you plan to stop mining and want to ensure all pending work is processed and claimed cleanly before shutting down.

## Running the Miner

Once all the prerequisites, building, and configuration steps are complete, you can start the miner:

1. Start the miner with your configuration:
   ```bash
   ./gobius --config config.json
   ```

Congratulations! Your Gobius miner should now be running.

## Troubleshooting Common Issues

Here are some common problems users might encounter during setup:

**1. Connection Issues (SSH/Server)**

*   **Symptom**: Cannot connect to the server via SSH.
*   **Cause**: Wrong IP address or username, firewall blocking port 22, incorrect password or SSH key setup.
*   **Solution**: Double-check your connection details. Ensure your cloud provider's firewall allows SSH (port 22). Verify your SSH key setup if using key-based authentication.

*   **Symptom**: Gobius stops running after closing the terminal/SSH session.
*   **Cause**: Not running Gobius (or `ipfs daemon`) inside `screen` or `tmux`.
*   **Solution**: Always start Gobius and the IPFS daemon inside a `screen` or `tmux` session and detach properly before closing your SSH connection. Refer to the guide's intro section for basic commands.

**2. Prerequisite & Build Issues**

*   **Symptom**: `go` or `git` command not found.
*   **Cause**: Go or Git is not installed, or not added to the system's PATH.
*   **Solution**: Re-follow the installation steps for Go and Git in the Prerequisites section, ensuring the PATH is updated correctly (`source ~/.bashrc` or restart terminal).

*   **Symptom**: `go build` fails with errors.
*   **Cause**: Incorrect Go version (< 1.22), network issues preventing dependency downloads, corrupted download, missing Git submodules (if `git clone` wasn't run with `--recursive`).
*   **Solution**: Verify Go version (`go version`). Check internet connectivity. Try deleting the `gobius` directory and re-cloning with `git clone --recursive`. Run `go clean -modcache` and try `go build` again.

*   **Symptom**: `ipfs init` or `ipfs daemon` fails.
*   **Cause**: Permission issues (try with `sudo` if appropriate, though generally not recommended long-term), port conflicts (another service using 4001 or 5001), insufficient disk space, corrupted IPFS repository.
*   **Solution**: Check for running processes using the ports (`sudo lsof -i :4001`, `sudo lsof -i :5001`). Check disk space (`df -h`). Try removing the IPFS repo (`rm -rf ~/.ipfs`) and re-running `ipfs init --profile server`.

*   **Symptom**: Gobius fails to start with a fatal error mentioning `CGO_ENABLED=0` and `go-sqlite3 requires cgo to work`.
    Example log: `{"level":"fatal","error":"Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work..."}`
*   **Cause**: Gobius was built on a system missing a C compiler (`gcc`), which is required by the `go-sqlite3` library. Go automatically disables CGO when a C compiler is not found during the build.
*   **Solution**: Install the necessary C compiler and build tools, then rebuild Gobius.
    ```bash
    # On Debian/Ubuntu
    sudo apt update && sudo apt install gcc build-essential
    
    # Then navigate back to the gobius directory and rebuild
    cd path/to/gobius
    go build
    ```

**3. Configuration (`config.json`) Errors**

*   **Symptom**: Gobius fails to start immediately, possibly mentioning JSON parsing errors or missing keys.
*   **Cause**: Syntax errors in `config.json` (e.g., missing comma, extra comma, incorrect bracket/brace, wrong quote type), incorrect file permissions.
*   **Solution**: Carefully review your `config.json` for syntax errors. Use a JSON validator (many online tools available) to check the structure. Ensure the file has read permissions for the user running Gobius.

*   **Symptom**: Gobius starts but fails with errors related to RPC connection.
*   **Cause**: Incorrect `blockchain.rpc_url` in `config.json` (must be a WebSocket `wss://` or local IPC path, not `http://`), dead/invalid endpoint, firewall blocking outbound connection.
*   **Solution**: Verify the `blockchain.rpc_url` format and ensure it's active. Use the public `wss://` endpoint or one from a provider like QuickNode as suggested.

*   **Symptom**: Gobius starts but fails with authentication or signature errors.
*   **Cause**: Incorrect private keys in `config.json`. Keys are missing, have `0x` prefix (they shouldn't), or are simply the wrong keys for the intended wallets.
*   **Solution**: Double-check all private keys. Ensure they're raw private keys **without** the `0x` prefix. **NEVER share these keys.**

*   **Symptom**: Gobius starts but fails with IPFS connection errors.
*   **Cause**: IPFS daemon not running, incorrect `ipfs.http_client.url` in `config.json` (if IPFS is remote or on a non-default port), firewall blocking connection to port 5001.
*   **Solution**: Ensure the `ipfs daemon` is running (use `ps aux | grep ipfs`). Verify the URL in `config.json` matches where the daemon is listening (default is `/ip4/127.0.0.1/tcp/5001` for local). Check firewall rules if IPFS is remote.

*   **Symptom**: Gobius runs but cannot process tasks, possibly with Cog/ML errors.
*   **Cause**: Incorrect Cog model URL in `ml.cog.<model-id>.url` (typo, wrong IP/port, missing `/predictions`, wrong model ID as the key), GPU instance (RunPod/Vast.ai) not running or inaccessible, firewall blocking connection from Gobius server to GPU instance.
*   **Solution**: Verify the Cog URL in `config.json`, ensuring it matches the running GPU instance's accessible URL and ends in `/predictions`. Ensure the GPU instance is running and accessible. Check firewalls on both the Gobius server (outbound) and GPU instance host (inbound port 5000).

**4. GPU Instance Issues (RunPod/Vast.ai)**

*   **Symptom**: Cannot connect to the Cog model URL / URL gives errors.
*   **Cause**: Instance not set up correctly (wrong GPU type - must be A100 80GB, port 5000 not exposed), Cog container failed to start, incorrect URL copied.
*   **Solution**: Double-check the instance setup against the guide's requirements (A100 80GB, port 5000 exposed, correct image `r8.im/kasumi-1/qwen-qwq-32b`). Check the instance logs on RunPod/Vast.ai for container errors. Re-verify the URL extraction steps, ensuring it includes `https://` or `http://` and ends in `/predictions`.

**5. Update Issues**

*   **Symptom**: `git pull` fails with conflicts.
*   **Cause**: You made local changes to tracked files that conflict with incoming updates.
*   **Solution**: Stash your changes (`git stash`), pull (`git pull origin main`), then re-apply your stash (`git stash pop`) and resolve conflicts manually. Alternatively, commit your changes to a separate branch before pulling.

*   **Symptom**: Gobius behaves unexpectedly after an update.
*   **Cause**: Forgetting to rebuild (`go build`) after pulling changes, or configuration requirements changed in the update.
*   **Solution**: Always run `go build` after `git pull`. Read the update's release notes carefully for any required changes to your `config.json`.

*   **Symptom**: Gobius fails to start after an update, potentially with database errors.
*   **Cause**: An automatic database migration required by the update failed.
*   **Solution**: Check the Gobius logs for specific database error messages. Report the issue with logs to the developers. If necessary, you may need to restore your database from the backup you made *before* updating and seek further assistance.
