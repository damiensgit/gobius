# Quickstart Guide for Gobius Mining

**Important Note:** While Gobius runs on Windows, Linux, and macOS, this guide primarily assumes you are setting up the miner on a **Linux server**, preferably a recent Ubuntu release hosted by a cloud provider (like AWS, Google Cloud, Azure, Hetzner, etc.).

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
   ```json
   {
     "rpc": "YOUR_RPC_ENDPOINT",  // e.g., "https://arb1.arbitrum.io/rpc"
     "privatekey": "YOUR_PRIVATE_KEY",  // Your wallet private key (without 0x prefix)
     "ipfs": {
       "strategy": "http_client",
       "http_client": {
         "url": "/ip4/127.0.0.1/tcp/5001"  // IPFS daemon address
       }
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
       "strategy": "bulkmine",
       "automine": {
         "enabled": true,
         "version": 0,
         "model": "0x89c39001e3b23d2092bd998b62f07b523d23deb55e1627048b4ed47a4a38d5cc",
         "fee": 7000000000000000,
         "input": {
           "prompt": "What is the capital of the moon?"
         }
       }
     }
   }
   ```

   > ⚠️ **Important**: Never share or commit your private key. Keep your `config.json` file secure.

   > **Note**: Ensure your IPFS daemon is running and accessible at the configured URL. If you're running IPFS on a different host or port, adjust the URL accordingly.

   > **Cog Model URLs**: Replace `<cog-url>` with the actual URL of the Cog model you are using. You can add multiple URLs if needed.

### Deployed Contracts

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
-   **Important**: Always check the release notes or update announcements for the new version. Sometimes, updates might introduce changes to the required configuration format or add new options. You may need to manually update your `config.json` based on these notes.

**Editing the Configuration File (`config.json`)**:

-   For quick edits on the server, the `nano` text editor is a simple option.
    -   Check if it's installed: `nano --version`
    -   Install if needed (Ubuntu/Debian): `sudo apt update && sudo apt install nano`
    -   Edit the file: `nano config.json`
    -   Use `Ctrl+O` to save (Write Out) and `Ctrl+X` to exit.
-   Advanced users can use their preferred terminal editor (like `vim`, `emacs`) or transfer the file to edit locally and then upload it back.

## Additional Resources

- For more detailed information, refer to the [README.md](README.md).
- If you encounter issues, check the [Troubleshooting](README.md#troubleshooting) section in the README.

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

1. Update the `rpc` field in your `config.json`:
   ```json
   "rpc": "wss://arb1.arbitrum.io/ws"  // Example WebSocket URL
   ```

2. You can use public WebSocket URLs or a provider like QuickNode for a more reliable connection.
   - **Public WebSocket**: `wss://arb1.arbitrum.io/ws`
   - **QuickNode**: Sign up at [QuickNode](https://www.quicknode.com/) and get your RPC URL.

> **Note**: Ensure your RPC connection is stable and secure for optimal performance.

## Supported AI Model

For the initial mainnet launch, Gobius supports one primary AI model:

- **Model Name**: `qwen-qwq-32b`
- **Model ID**: `0x89c39001e3b23d2092bd998b62f07b523d23deb55e1627048b4ed47a4a38d5cc`

This Model ID is used in the `config.json` file when specifying the Cog model URL. Ensure you are configuring your GPU instance and miner settings for this specific model.

## Setting Up a GPU Instance on RunPod and Vast.ai

To run the Cog model `r8.im/kasumi-1/qwen-qwq-32b`, you need to set up a GPU instance with the following specifications:

**Important**: You **must** use an NVIDIA A100 GPU with 80GB of VRAM. No other GPU type or VRAM configuration is supported for this model.

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

> **Note**: Ensure your instance meets the model's resource requirements for optimal performance.

## Extracting the URL for Launched Instances

Once your GPU instance is running on RunPod or Vast.ai, you'll need to extract the correct URL to use in your configuration. Remember that the URL should end with `/predictions`.

### RunPod

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

### Vast.ai

1.  **Access the Instance**: Go to your "Instances" page in the Vast.ai dashboard.
2.  **Find the URL**: Click the "Connect" button for your instance. Look for the port forwarding URL mapped to port 5000. It often looks like `http://<numeric-address>:<port>` or a direct domain name if configured.
3.  **Update Config**: Insert the full URL, **appending `/predictions`**, into your `config.json` in the Cog model section:
    ```json
    "ml": {
      "strategy": "cog",
      "cog": {
        "0x89c39001e3b23d2092bd998b62f07b523d23deb55e1627048b4ed47a4a38d5cc": { // Ensure this is the correct model ID
          "url": [
            "http://<vast-instance-url>:5000/predictions" // Replace with your actual Vast.ai URL and port, ending in /predictions
          ]
        }
      }
    }
    ```

> **Note**: Ensure the URLs are accessible, include `http://` or `https://` as appropriate, end with the correct endpoint (`/predictions`), and are correctly formatted to avoid connectivity issues.

## Initial Mining Setup: Auto-Mining

The recommended initial setup for Gobius is **auto-mining**. This means your miner will automatically generate its own tasks and then solve them.

This is configured in the `strategies` section of your `config.json` file. The example configuration (`config.example.json`) is already set up for this:

```json
"strategies": {
  "model": "0x89c39001e3b23d2092bd998b62f07b523d23deb55e1627048b4ed47a4a38d5cc", // Ensure this matches the supported model ID
  "strategy": "bulkmine", // Use bulkmine for auto-mining
  "automine": {
    "enabled": true, // Enable auto-mining
    "version": 0,
    "model": "0x89c39001e3b23d2092bd998b62f07b523d23deb55e1627048b4ed47a4a38d5cc", // Ensure this matches the supported model ID
    "fee": 7000000000000000, // Default fee
    "input": {
      "prompt": "What is the capital of the moon?" // You can customize this prompt
    }
  }
},
```

**Key Settings for Auto-Mining:**

-   `strategy`: Must be set to `"bulkmine"`.
-   `automine.enabled`: Must be set to `true`.
-   `model` (in both `strategies` and `automine`): Ensure this matches the Model ID of the supported model (`0x89c39001e3b23d2092bd998b62f07b523d23deb55e1627048b4ed47a4a38d5cc`).
-   `automine.input.prompt`: This is the main setting you might want to change. It defines the input for the tasks your miner generates. The default prompt is just an example.

For starting, you typically only need to ensure the `model` IDs are correct and potentially customize the `prompt`.

## Running the Miner

Once all the prerequisites, building, and configuration steps are complete, you can start the miner:

1.  Start the miner with your configuration:
    ```bash
    ./gobius --config config.json
    ```

Congratulations! Your Gobius miner should now be running.



## Troubleshooting Common Issues

Here are some common problems users might encounter during setup:

**1. Connection Issues (SSH/Server)**

*   **Symptom**: Cannot connect to the server via SSH.
*   **Causes**: Wrong IP address or username, firewall blocking port 22, incorrect password or SSH key setup.
*   **Solution**: Double-check your connection details. Ensure your cloud provider's firewall allows SSH (port 22). Verify your SSH key setup if using key-based authentication.

*   **Symptom**: Gobius stops running after closing the terminal/SSH session.
*   **Cause**: Not running Gobius (or `ipfs daemon`) inside `screen` or `tmux`.
*   **Solution**: Always start Gobius and the IPFS daemon inside a `screen` or `tmux` session and detach properly before closing your SSH connection. Refer to the guide's intro section for basic commands.

**2. Prerequisite & Build Issues**

*   **Symptom**: `go` or `git` command not found.
*   **Cause**: Go or Git is not installed, or not added to the system's PATH.
*   **Solution**: Re-follow the installation steps for Go and Git in the Prerequisites section, ensuring the PATH is updated correctly (`source ~/.bashrc` or restart terminal).

*   **Symptom**: `go build` fails with errors.
*   **Causes**: Incorrect Go version (< 1.22), network issues preventing dependency downloads, corrupted download, missing Git submodules (if `git clone` wasn't run with `--recursive`).
*   **Solution**: Verify Go version (`go version`). Check internet connectivity. Try deleting the `gobius` directory and re-cloning with `git clone --recursive`. Run `go clean -modcache` and try `go build` again.

*   **Symptom**: `ipfs init` or `ipfs daemon` fails.
*   **Causes**: Permission issues (try with `sudo` if appropriate, though generally not recommended long-term), port conflicts (another service using 4001 or 5001), insufficient disk space, corrupted IPFS repository.
*   **Solution**: Check for running processes using the ports (`sudo lsof -i :4001`, `sudo lsof -i :5001`). Check disk space (`df -h`). Try removing the IPFS repo (`rm -rf ~/.ipfs`) and re-running `ipfs init --profile server`.

**3. Configuration (`config.json`) Errors**

*   **Symptom**: Gobius fails to start immediately, possibly mentioning JSON parsing errors or missing keys.
*   **Cause**: Syntax errors in `config.json` (e.g., missing comma, extra comma, incorrect bracket/brace, wrong quote type), incorrect file permissions.
*   **Solution**: Carefully review your `config.json` for syntax errors. Use a JSON validator (many online tools available) to check the structure. Ensure the file has read permissions for the user running Gobius.

*   **Symptom**: Gobius starts but fails with errors related to RPC connection.
*   **Cause**: Incorrect `rpc` URL in `config.json` (must be a WebSocket `wss://` or local IPC path, not `http://`), dead/invalid endpoint, firewall blocking outbound connection.
*   **Solution**: Verify the `rpc` URL format and ensure it's active. Use the public `wss://` endpoint or one from a provider like QuickNode as suggested.

*   **Symptom**: Gobius starts but fails with authentication or signature errors.
*   **Cause**: Incorrect `privatekey` in `config.json`. Key is missing, has `0x` prefix (it shouldn't), or is simply the wrong key for the intended wallet.
*   **Solution**: Double-check the `privatekey`. Ensure it's the raw private key **without** the `0x` prefix. **NEVER share this key.**

*   **Symptom**: Gobius starts but fails with IPFS connection errors.
*   **Cause**: IPFS daemon not running, incorrect `ipfs.http_client.url` in `config.json` (if IPFS is remote or on a non-default port), firewall blocking connection to port 5001.
*   **Solution**: Ensure the `ipfs daemon` is running (use `ps aux | grep ipfs`). Verify the URL in `config.json` matches where the daemon is listening (default is `/ip4/127.0.0.1/tcp/5001` for local). Check firewall rules if IPFS is remote.

*   **Symptom**: Gobius runs but cannot process tasks, possibly with Cog/ML errors.
*   **Cause**: Incorrect Cog model URL in `ml.cog.<model-id>.url` (typo, wrong IP/port, missing `/predictions`, wrong model ID as the key), GPU instance (RunPod/Vast.ai) not running or inaccessible, firewall blocking connection from Gobius server to GPU instance on port 5000.
*   **Solution**: Verify the Cog URL in `config.json`, ensuring it matches the running GPU instance's accessible URL and ends in `/predictions`. Ensure the GPU instance is running and accessible. Check firewalls on both the Gobius server (outbound) and GPU instance host (inbound on port 5000).

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
