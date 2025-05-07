// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;
import "arbius-contracts/interfaces/IArbius.sol"; 

contract BulkTasks {
    address public basetoken;
    IArbius public engine;

    constructor(address basetoken_, address engine_) {
        basetoken = basetoken_;
        engine = IArbius(engine_);
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
        returns (IArbius.Solution[] memory)
    {
        uint256 len = taskids_.length;
        IArbius.Solution[] memory solutionInfos = new IArbius.Solution[](
            len
        );

        for (uint256 i = 0; i < len; i++) {
            solutionInfos[i] = engine.solutions(taskids_[i]);
        }
        return solutionInfos;
    }
}
