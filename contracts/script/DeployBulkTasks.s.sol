// contracts/script/DeployBulkTasks.s.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {Script, console} from "forge-std/Script.sol";
import {BulkTasks} from "../src/BulkTasks.sol"; 

contract DeployBulkTasks is Script {
    function run() external returns (BulkTasks) {
        string memory deployEnv = vm.envString("DEPLOY_ENV");
        if (bytes(deployEnv).length == 0) {
            deployEnv = "local";
        }
        console.log("Using environment:", deployEnv);

        string memory configFilePath = string.concat("./config/config.", deployEnv, ".json");
        console.log("Reading config from:", configFilePath);

        string memory json = vm.readFile(configFilePath);

        address baseTokenAddr = vm.parseJsonAddress(json, ".baseTokenAddress");
        address engineAddr = vm.parseJsonAddress(json, ".engineAddress");

        require(baseTokenAddr != address(0), "BaseTokenAddr not found/invalid in JSON");
        require(engineAddr != address(0), "EngineAddr not found/invalid in JSON");

        console.log("Deploying BulkTasks with:");
        console.log("  Base Token:", baseTokenAddr);
        console.log("  Engine:", engineAddr);

        vm.startBroadcast();
        BulkTasks bulkTasks = new BulkTasks(baseTokenAddr, engineAddr);
        vm.stopBroadcast();

        console.log("BulkTasks deployed to:", address(bulkTasks));
        return bulkTasks;
    }
}