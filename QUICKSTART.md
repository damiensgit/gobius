# Quickstart Guide for Gobius Mining

This guide will walk you through setting up Gobius for mining on the Arbius network. Follow these steps to get started quickly.

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
   - Install IPFS by following the official guide at [https://docs.ipfs.tech/install/command-line/#install-official-binary-distributions](https://docs.ipfs.tech/install/command-line/#install-official-binary-distributions).
   - Initialize your IPFS node:
     ```bash
     ipfs init
     ```
   - Run the IPFS daemon:
     ```bash
     ipfs daemon
     ```
   - **Firewall**: Ensure ports 4001 (TCP/UDP), 5001 (TCP), and 8080 (TCP) are open in your firewall, especially if using a cloud provider.

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
         "0xdafc13d7e9529e278bf40038524f20841ba8879c0b5d7667d438cd60c64a118d": {
           "url": [
             "<cog-url>"  // Replace with your Cog model URL
           ]
         }
       }
     }
     // ... other optional settings ...
   }
   ```

   > ⚠️ **Important**: Never share or commit your private key. Keep your `config.json` file secure.

   > **Note**: Ensure your IPFS daemon is running and accessible at the configured URL. If you're running IPFS on a different host or port, adjust the URL accordingly.

   > **Cog Model URLs**: Replace `<cog-url>` with the actual URL of the Cog model you are using. You can add multiple URLs if needed.

### Deployed Contracts

## Running the Miner

1. Start the miner with your configuration:
   ```bash
   ./gobius --config config.json
   ```

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

2. Replace `<model-id>` with the ID of your model and `<cog-url>` with the actual URL. You can add multiple URLs if needed.

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

Once your GPU instance is running on RunPod or Vast.ai, you'll need to extract the correct URL to use in your configuration. Remember that the URL should end with `/predict`.

### RunPod

1. **Access the Pod**: Navigate to "My Pods" in the RunPod dashboard and select your running pod.
2. **Find the URL**: Go to the "Connect" tab. Look for the HTTP/HTTPS connection URL provided for port 5000. It will typically look like `https://<pod-id>-5000.proxy.runpod.net`.
3. **Update Config**: Place this full URL, **appending `/predict`**, in your `config.json` under the Cog model section:
    ```json
    "ml": {
      "strategy": "cog",
      "cog": {
        "<model-id>": { // Replace <model-id> with the actual model ID
          "url": [
            "https://<pod-id>-5000.proxy.runpod.net/predict" // Replace with your actual RunPod URL, ending in /predict
          ]
        }
      }
    }
    ```

### Vast.ai

1. **Access the Instance**: Go to your "Instances" page in the Vast.ai dashboard.
2. **Find the URL**: Click the "Connect" button for your instance. Look for the port forwarding URL mapped to port 5000. It often looks like `http://<numeric-address>:<port>` or a direct domain name if configured.
3. **Update Config**: Insert the full URL, **appending `/predict`**, into your `config.json` in the Cog model section:
    ```json
    "ml": {
      "strategy": "cog",
      "cog": {
        "<model-id>": { // Replace <model-id> with the actual model ID
          "url": [
            "http://<vast-instance-url>:5000/predict" // Replace with your actual Vast.ai URL and port, ending in /predict
          ]
        }
      }
    }
    ```

> **Note**: Ensure the URLs are accessible, include `http://` or `https://` as appropriate, end with `/predict`, and are correctly formatted to avoid connectivity issues. 