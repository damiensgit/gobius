// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "forge-std/Test.sol";
import {console} from "forge-std/console.sol";
import "../src/BulkTasks.sol"; // Assuming BulkTasks.sol is in contracts/src/

contract BulkTasksForkTest is Test {

    address constant LIVE_ENGINE_ADDRESS = address(0x9b51Ef044d3486A1fB0A2D55A6e0CeeAdd323E66); 
    address constant LIVE_BASE_TOKEN_ADDRESS = address(0x4a24B101728e07A52053c13FB4dB2BcF490CAbc3);

    bytes32 constant problematicTaskId = 0x9d49af40436ea3a29a12b63064d2ec5a0f3706eb4635af6d453d04cfed474121;

    BulkTasks public bulkTasksContract;

    function setUp() public {
        bulkTasksContract = new BulkTasks(LIVE_BASE_TOKEN_ADDRESS, LIVE_ENGINE_ADDRESS);
    }

    function test_GetBulkCombinedInfo_SingleProblematicTask_Forked() public {
        // require(address(bulkTasksContract) != address(0), "New BulkTasks contract not deployed"); // Or use the live one if comparing
        require(problematicTaskId != bytes32(0), "Problematic Task ID not set");

        bytes32[] memory taskIdsToTest = new bytes32[](1);
        taskIdsToTest[0] = problematicTaskId;
        taskIdsToTest[0] = 0x4fc1d4456f8d48bb6d16094b1097e18df89ac8d24e8c987de02fd46d6281a5d2;
        

        console.log("Testing getBulkCombinedTaskInfo with Task ID:");
        console.logBytes32(problematicTaskId);
        console.log("BulkTasks Contract Address:", address(bulkTasksContract));
                
        BulkTasks.TaskSolutionWithContestationInfo[] memory infos = bulkTasksContract.getBulkCombinedTaskInfo(taskIdsToTest);

        require(infos.length == 2, "Should return info for 2 tasks");
        console.log("Successfully called getBulkCombinedTaskInfo.");

        //console.log("---- Task Info for ID:", infos[0].taskId, "----");
        console.log("  Solution Exists:", infos[0].solutionExists);
        console.log("  Solution Validator:", infos[0].solutionValidator);
        console.log("  Solution Blocktime:", infos[0].solutionBlocktime);
        console.log("  Solution Claimed:", infos[0].solutionClaimed);
       // console.logBytes(infos[0].solutionCid); // Use if CID is small, otherwise can be too verbose

        console.log("  Contestation Exists:", infos[0].contestationExists);
        console.log("  Contestation Validator:", infos[0].contestationValidator);
        console.log("  Contestation Blocktime:", infos[0].contestationBlocktime);
        console.log("  Contestation Slash Amount:", infos[0].contestationSlashAmount);
    }

    function test_GetSolutions_SingleProblematicTask_Forked() public {
        bytes32 singleProblematicTaskId = 0x4fc1d4456f8d48bb6d16094b1097e18df89ac8d24e8c987de02fd46d6281a5d2;
        bytes32[] memory taskIdsToTest = new bytes32[](1);
        taskIdsToTest[0] = singleProblematicTaskId;

        console.log("Testing getSolutions with Task ID:");
        console.logBytes32(singleProblematicTaskId);

        console.log("BulkTasks Contract Address:", address(bulkTasksContract));

        IArbiusEngine.EngineSolution[] memory solutions = bulkTasksContract.getSolutions(taskIdsToTest);

        require(solutions.length == 1, "Should return info for 1 task");
        console.log("Successfully called getSolutions.");
        console.log("  Solution Validator:", solutions[0].validator);
        console.log("  Solution Blocktime:", solutions[0].blocktime);
        console.log("  Solution Claimed:", solutions[0].claimed);
    }
} 