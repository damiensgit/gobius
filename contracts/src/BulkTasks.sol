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
    }

    function contestations(bytes32 taskId) external view returns (EngineContestation memory);
    function claimSolution(bytes32 taskid) external;
    function signalCommitment(bytes32 commitment) external;
    function commitments(bytes32 commitment) external view returns (uint256);
    function solutions(bytes32 taskId) external view returns (address validator, uint64 blocktime, bool claimed);
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
            try engine.claimSolution(_taskids[i]) {
            } catch {
            }
        }
    }

    function bulkSignalCommitment(bytes32[] calldata commitments_) public {
        uint256 len = commitments_.length; // Cache array length
        for (uint256 i = 0; i < len; i++) {
            try engine.signalCommitment(commitments_[i]) {
            } catch {
            }
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
    /// @return An array of Solution structs containing details for each task ID (without CIDs).
    function getSolutions(bytes32[] calldata taskids_)
        external
        view
        returns (IArbiusEngine.EngineSolution[] memory) // This EngineSolution struct in IArbiusEngine does not have cid
    {
        uint256 len = taskids_.length;
        IArbiusEngine.EngineSolution[] memory solutionInfos = new IArbiusEngine.EngineSolution[](
            len
        );

        for (uint256 i = 0; i < len; i++) {
            // Call the modified engine.solutions() that doesn't return cid
            // and populate the EngineSolution struct which also doesn't have cid.
            (address sValidator, uint64 sBlocktime, bool sClaimed) = engine.solutions(taskids_[i]);
            solutionInfos[i].validator = sValidator;
            solutionInfos[i].blocktime = sBlocktime;
            solutionInfos[i].claimed = sClaimed;
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
        returns (IArbiusEngine.EngineContestation[] memory)
    {
        uint256 len = taskids_.length;
        IArbiusEngine.EngineContestation[] memory contestationInfos = new IArbiusEngine.EngineContestation[](len);

        for (uint256 i = 0; i < len; i++) {
            IArbiusEngine.EngineContestation memory engineContestation = engine.contestations(taskids_[i]);

            contestationInfos[i].validator = engineContestation.validator;
            contestationInfos[i].blocktime = engineContestation.blocktime;
            contestationInfos[i].finish_start_index = engineContestation.finish_start_index; 
            contestationInfos[i].slashAmount = engineContestation.slashAmount;
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

            (address sValidator, uint64 sBlocktime, bool sClaimed) = engine.solutions(currentTaskId);
            allInfos[i].solutionValidator = sValidator;
            allInfos[i].solutionBlocktime = sBlocktime;
            allInfos[i].solutionClaimed = sClaimed;
            allInfos[i].solutionExists = sBlocktime > 0; 
            
            IArbiusEngine.EngineContestation memory cont = engine.contestations(currentTaskId);
            allInfos[i].contestationValidator = cont.validator;
            allInfos[i].contestationBlocktime = cont.blocktime;
            allInfos[i].contestationFinishStartIndex = cont.finish_start_index;
            allInfos[i].contestationSlashAmount = cont.slashAmount;
            allInfos[i].contestationExists = cont.validator != address(0);
        }
        return allInfos;
    }
}
