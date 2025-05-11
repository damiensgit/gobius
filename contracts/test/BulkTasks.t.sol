// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "forge-std/Test.sol";
import "../src/BulkTasks.sol"; // Adjust path if necessary

// Minimal interface for what BulkTasks calls on the Engine
interface IMockEngine {
    struct EngineContestation {
        address validator;
        uint64 blocktime;
        uint32 finish_start_index;
        uint256 slashAmount;
    }

    struct EngineSolution {
        address validator;
        uint64 blocktime;
        bool claimed;
    }

    function contestations(bytes32 taskId) external view returns (EngineContestation memory);
    function commitments(bytes32 commitment) external view returns (uint256);
    function solutions(bytes32 taskId) external view returns (EngineSolution memory);

    // Functions for test setup
    function mock_setCommitment(bytes32 commitment, uint256 blockNumber) external;
    function mock_setSolution(bytes32 taskId, address validator, uint64 blocktime, bool claimed) external;
    function mock_setContestation(bytes32 taskId, address validator, uint64 blocktime, uint32 finish_start_index, uint256 slashAmount) external;
    function mock_clearContestation(bytes32 taskId) external;
}

// Simple Mock ERC20 for testing purposes
contract MockERC20 {
    string public name = "Mock Token";
    string public symbol = "MTK";
    uint8 public decimals = 18;
    mapping(address => mapping(address => uint256)) public allowances;

    // Add other ERC20 functions if BulkTasks constructor needs them
    // For now, just being a contract might be enough.
    function approve(address spender, uint256 amount) public returns (bool) {
        allowances[msg.sender][spender] = amount;
        // Emit an Approval event if your tests or contract logic expect it
        // event Approval(address indexed owner, address indexed spender, uint256 value);
        // emit Approval(msg.sender, spender, amount);
        return true;
    }
}

contract MockEngine is IMockEngine {
    mapping(bytes32 => uint256) public _commitments;
    mapping(bytes32 => EngineSolution) public _solutions;
    mapping(bytes32 => EngineContestation) public _contestations;

    // --- Mock Setup Functions ---
    function mock_setCommitment(bytes32 commitment, uint256 blockNumber) external {
        _commitments[commitment] = blockNumber;
    }

    function mock_setSolution(bytes32 taskId, address validator, uint64 blocktime, bool claimed) external {
        _solutions[taskId] = EngineSolution(validator, blocktime, claimed);
    }

    function mock_setContestation(bytes32 taskId, address validator, uint64 blocktime, uint32 finish_start_index, uint256 slashAmount) external {
        _contestations[taskId] = EngineContestation(validator, blocktime, finish_start_index, slashAmount);
    }
    
    function mock_clearContestation(bytes32 taskId) external {
        delete _contestations[taskId];
    }

    // --- Engine View Functions (called by BulkTasks) ---
    function commitments(bytes32 commitment) external view override returns (uint256) {
        return _commitments[commitment];
    }

    function solutions(bytes32 taskId) external view override returns (EngineSolution memory) {
        return _solutions[taskId];
    }

    function contestations(bytes32 taskId) external view override returns (EngineContestation memory) {
        return _contestations[taskId];
    }
}

contract BulkTasksTest is Test {
    BulkTasks public bulkTasksContract;
    MockEngine public mockEngine;
    MockERC20 public mockBaseToken; // New

    function setUp() public {
        mockBaseToken = new MockERC20(); // Deploy Mock ERC20
        mockEngine = new MockEngine();
        bulkTasksContract = new BulkTasks(address(mockBaseToken), address(mockEngine)); // Use address of deployed mock token
    }

    // --- Test for getCommitments ---
    function test_GetCommitments_MultipleExisting() public {
        bytes32 c1 = keccak256("commit1");
        bytes32 c2 = keccak256("commit2");
        mockEngine.mock_setCommitment(c1, 100);
        mockEngine.mock_setCommitment(c2, 200);

        bytes32[] memory commitments = new bytes32[](2);
        commitments[0] = c1;
        commitments[1] = c2;

        uint256[] memory blockNumbers = bulkTasksContract.getCommitments(commitments);
        assertEq(blockNumbers.length, 2);
        assertEq(blockNumbers[0], 100);
        assertEq(blockNumbers[1], 200);
    }

    function test_GetCommitments_OneNonExisting() public {
        bytes32 c1 = keccak256("commit_exist");
        bytes32 c2 = keccak256("commit_not_exist");
        mockEngine.mock_setCommitment(c1, 100);

        bytes32[] memory commitments = new bytes32[](2);
        commitments[0] = c1;
        commitments[1] = c2;

        uint256[] memory blockNumbers = bulkTasksContract.getCommitments(commitments);
        assertEq(blockNumbers.length, 2);
        assertEq(blockNumbers[0], 100);
        assertEq(blockNumbers[1], 0); // Non-existing should return 0
    }

    function test_GetCommitments_EmptyInput() public {
        bytes32[] memory commitments = new bytes32[](0);
        uint256[] memory blockNumbers = bulkTasksContract.getCommitments(commitments);
        assertEq(blockNumbers.length, 0);
    }

    // --- Test for getSolutions ---
    function test_GetSolutions_MultipleExisting() public {
        bytes32 t1 = keccak256("task1_sol");
        bytes32 t2 = keccak256("task2_sol");
        address val1 = address(0x1);
        address val2 = address(0x2);

        mockEngine.mock_setSolution(t1, val1, 1000, false);
        mockEngine.mock_setSolution(t2, val2, 2000, true);

        bytes32[] memory taskIds = new bytes32[](2);
        taskIds[0] = t1;
        taskIds[1] = t2;

        IArbiusEngine.EngineSolution[] memory solutions = bulkTasksContract.getSolutions(taskIds);
        assertEq(solutions.length, 2);
        assertEq(solutions[0].validator, val1);
        assertEq(solutions[0].blocktime, 1000);
        assertEq(solutions[0].claimed, false);

        assertEq(solutions[1].validator, val2);
        assertEq(solutions[1].blocktime, 2000);
        assertEq(solutions[1].claimed, true);
    }

    function test_GetSolutions_OneNonExisting() public {
        bytes32 t1 = keccak256("task_exist_sol");
        bytes32 t2 = keccak256("task_not_exist_sol");
        address val1 = address(0x1);
        mockEngine.mock_setSolution(t1, val1, 1000, false);
        
        bytes32[] memory taskIds = new bytes32[](2);
        taskIds[0] = t1;
        taskIds[1] = t2;

        IArbiusEngine.EngineSolution[] memory solutions = bulkTasksContract.getSolutions(taskIds);
        assertEq(solutions.length, 2);
        assertEq(solutions[0].validator, val1); // Existing one
        assertEq(solutions[0].blocktime, 1000);
        
        assertEq(solutions[1].validator, address(0)); // Non-existing should be zeroed
        assertEq(solutions[1].blocktime, 0);
        assertEq(solutions[1].claimed, false);
    }

    function test_GetSolutions_EmptyInput() public {
        bytes32[] memory taskIds = new bytes32[](0);
        IArbiusEngine.EngineSolution[] memory solutions = bulkTasksContract.getSolutions(taskIds);
        assertEq(solutions.length, 0);
    }

    // --- Test for getBulkCombinedTaskInfo ---
    function test_GetCombined_SolutionExists_NoContestation() public {
        bytes32 t1 = keccak256("task_sol_only");
        address val = address(0xA1);
        mockEngine.mock_setSolution(t1, val, 3000, false);
        mockEngine.mock_clearContestation(t1); // Ensure no contestation

        bytes32[] memory taskIds = new bytes32[](1);
        taskIds[0] = t1;

        BulkTasks.TaskSolutionWithContestationInfo[] memory infos = bulkTasksContract.getBulkCombinedTaskInfo(taskIds);
        assertEq(infos.length, 1);
        assertEq(infos[0].taskId, t1);
        assertTrue(infos[0].solutionExists);
        assertEq(infos[0].solutionValidator, val);
        assertEq(infos[0].solutionBlocktime, 3000);
        assertFalse(infos[0].solutionClaimed);
        assertFalse(infos[0].contestationExists);
        assertEq(infos[0].contestationValidator, address(0));
    }

    function test_GetCombined_ContestationExists_NoSolution() public {
        bytes32 t1 = keccak256("task_cont_only");
        address valCont = address(0xB1);
        mockEngine.mock_setContestation(t1, valCont, 4000, 1, 1 ether);
        // No solution set for t1

        bytes32[] memory taskIds = new bytes32[](1);
        taskIds[0] = t1;

        BulkTasks.TaskSolutionWithContestationInfo[] memory infos = bulkTasksContract.getBulkCombinedTaskInfo(taskIds);
        assertEq(infos.length, 1);
        assertEq(infos[0].taskId, t1);
        assertFalse(infos[0].solutionExists);
        assertEq(infos[0].solutionValidator, address(0));
        assertTrue(infos[0].contestationExists);
        assertEq(infos[0].contestationValidator, valCont);
        assertEq(infos[0].contestationBlocktime, 4000);
    }

     function test_GetCombined_BothExist_SolutionClaimed() public {
        bytes32 t1 = keccak256("task_both_claimed");
        address valSol = address(0xC1);
        address valCont = address(0xD1);
        bytes memory c = abi.encodePacked("cid_both");

        mockEngine.mock_setSolution(t1, valSol, 5000, true); // Solution is claimed
        mockEngine.mock_setContestation(t1, valCont, 6000, 2, 2 ether);

        bytes32[] memory taskIds = new bytes32[](1);
        taskIds[0] = t1;
        
        BulkTasks.TaskSolutionWithContestationInfo[] memory infos = bulkTasksContract.getBulkCombinedTaskInfo(taskIds);
        assertEq(infos.length, 1);
        assertEq(infos[0].taskId, t1);
        assertTrue(infos[0].solutionExists);
        assertEq(infos[0].solutionValidator, valSol);
        assertTrue(infos[0].solutionClaimed); // Verify claimed status
        assertTrue(infos[0].contestationExists);
        assertEq(infos[0].contestationValidator, valCont);
    }

    function test_GetCombined_NeitherExists() public {
        bytes32 t1 = keccak256("task_neither");
        // No solution or contestation set for t1
        mockEngine.mock_clearContestation(t1);

        bytes32[] memory taskIds = new bytes32[](1);
        taskIds[0] = t1;

        BulkTasks.TaskSolutionWithContestationInfo[] memory infos = bulkTasksContract.getBulkCombinedTaskInfo(taskIds);
        assertEq(infos.length, 1);
        assertEq(infos[0].taskId, t1);
        assertFalse(infos[0].solutionExists);
        assertFalse(infos[0].contestationExists);
    }

    function test_GetCombined_EmptyInput() public {
        bytes32[] memory taskIds = new bytes32[](0);
        BulkTasks.TaskSolutionWithContestationInfo[] memory infos = bulkTasksContract.getBulkCombinedTaskInfo(taskIds);
        assertEq(infos.length, 0);
    }
} 