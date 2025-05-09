// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;
import "arbius-contracts/interfaces/IBaseToken.sol";

interface IArbiusEngine {
    // Define the struct EXACTLY as it is in the Engine contract storage
    struct EngineContestation {
        address validator;
        uint64 blocktime;
        uint32 finish_start_index; // Match actual field name
        uint256 slashAmount;     // Match actual field name & type
    }
    
    struct EngineSolution { // Match actual engine solution struct
        address validator;
        uint64 blocktime;
        bool claimed;
        bytes cid;
    }

    function contestations(bytes32 taskId) external view returns (EngineContestation memory);
    function claimSolution(bytes32 taskid) external;
    function signalCommitment(bytes32 commitment) external;
    function commitments(bytes32 commitment) external view returns (uint256);
    function solutions(bytes32 taskId) external view returns (EngineSolution memory);
}


contract BulkTasks {
    address public basetoken;
    IArbiusEngine public engine;

    struct TaskSolutionWithContestationInfo {
        bytes32 taskId;
        // Solution part
        address solutionValidator;
        uint64 solutionBlocktime;
        bool solutionClaimed;
        bytes solutionCid;
        bool solutionExists;
        // Contestation part
        address contestationValidator;
        uint64 contestationBlocktime;
        uint32 contestationFinishStartIndex;
        uint256 contestationSlashAmount;
        bool contestationExists;
    }

    constructor(address basetoken_, address engine_) {
        basetoken = basetoken_;
        engine = IArbiusEngine(engine_);
        IBaseToken(basetoken).approve(engine_, type(uint256).max);
    }

    function claimSolutions(bytes32[] calldata _taskids) public {
        uint256 len = _taskids.length; // Cache array length
        for (uint256 i = 0; i < len; i++) {
            engine.claimSolution(_taskids[i]); // Call claimSolution directly
        }
    }

    function bulkSignalCommitment(bytes32[] calldata commitments_) public {
        uint256 len = commitments_.length; // Cache array length
        for (uint256 i = 0; i < len; i++) {
            engine.signalCommitment(commitments_[i]); // Call signalCommitment directly
        }
    }

    /// @notice Get the block number for multiple commitments.
    /// @param commitments_ Array of commitment hashes to query.
    /// @return An array of block numbers corresponding to each commitment hash (0 if not committed).
    function getCommitments(bytes32[] calldata commitments_)
        external
        view
        returns (uint256[] memory)
    {
        uint256 len = commitments_.length;
        uint256[] memory blockNumbers = new uint256[](len);

        for (uint256 i = 0; i < len; i++) {
            blockNumbers[i] = engine.commitments(commitments_[i]);
        }
        return blockNumbers;
    }

    /// @notice Get solution information for multiple task IDs.
    /// @param taskids_ Array of task IDs to query.
    /// @return An array of Solution structs containing details for each task ID.
    function getSolutions(bytes32[] calldata taskids_)
        external
        view
        returns (IArbiusEngine.EngineSolution[] memory)
    {
        uint256 len = taskids_.length;
        IArbiusEngine.EngineSolution[] memory solutionInfos = new IArbiusEngine.EngineSolution[](
            len
        );

        for (uint256 i = 0; i < len; i++) {
            solutionInfos[i] = engine.solutions(taskids_[i]);
        }
        return solutionInfos;
    }

    /// @notice Get contestation information for multiple task IDs.
    /// @param taskids_ Array of task IDs to query.
    /// @return An array of ContestationData structs containing details for each task ID.
    /// @dev Returns zero values if no contestation exists for a task ID.
    function getBulkContestations(bytes32[] calldata taskids_)
        external
        view
        // Return the locally defined struct array
        returns (IArbiusEngine.EngineContestation[] memory)
    {
        uint256 len = taskids_.length;
        IArbiusEngine.EngineContestation[] memory contestationInfos = new IArbiusEngine.EngineContestation[](len);

        for (uint256 i = 0; i < len; i++) {
            // Fetch the data from the engine's public mapping
            // Solidity allows assigning structs if they are compatible,
            // but to be explicit and safe, assign field by field if direct assignment fails.
            // Assuming engine.contestations returns a struct with the same fields:
            IArbiusEngine.EngineContestation memory engineContestation = engine.contestations(taskids_[i]);

            // Assign fields to our local struct
            contestationInfos[i].validator = engineContestation.validator;
            contestationInfos[i].blocktime = engineContestation.blocktime;
            contestationInfos[i].finish_start_index = engineContestation.finish_start_index; // Match field name from Engine
            contestationInfos[i].slashAmount = engineContestation.slashAmount; // Match field name from Engine
        }
        return contestationInfos;
    }

      function getBulkCombinedTaskInfo(bytes32[] calldata taskIds)
        external
        view
        returns (TaskSolutionWithContestationInfo[] memory)
    {
        uint256 len = taskIds.length;
        TaskSolutionWithContestationInfo[] memory allInfos = new TaskSolutionWithContestationInfo[](len);

        for (uint i = 0; i < len; i++) {
            bytes32 currentTaskId = taskIds[i];
            allInfos[i].taskId = currentTaskId;

            // Fetch Solution Data
            IArbiusEngine.EngineSolution memory sol = engine.solutions(currentTaskId);
            allInfos[i].solutionValidator = sol.validator;
            allInfos[i].solutionBlocktime = sol.blocktime;
            allInfos[i].solutionClaimed = sol.claimed;
            allInfos[i].solutionCid = sol.cid;
            allInfos[i].solutionExists = sol.blocktime > 0;

            // Fetch Contestation Data
            IArbiusEngine.EngineContestation memory cont = engine.contestations(currentTaskId);
            allInfos[i].contestationValidator = cont.validator;
            allInfos[i].contestationBlocktime = cont.blocktime;
            allInfos[i].contestationFinishStartIndex = cont.finish_start_index; 
            allInfos[i].contestationSlashAmount = cont.slashAmount; 
            // A contestation "exists" if its validator is not the zero address
            allInfos[i].contestationExists = cont.validator != address(0);
        }
        return allInfos;
    }
}
