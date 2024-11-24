//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IArbius {
    struct Solution {
        address validator;
        uint64 blocktime;
        bool claimed;
        bytes cid; // ipfs cid
    }

    function solutions(bytes32 taskid) external view returns (Solution memory);
}

/**
 @dev This contract is not meant to be deployed. Instead, use a static call with the
      deployment bytecode as payload.
 */
contract GetSolutionBlocktimeBatchRequest {
    constructor(
        bytes32[] calldata _taskids,
        address _arbius
    ) {
        uint256 taskLen = _taskids.length;

        // There is a max number of pool as a too big returned data times out the rpc
        uint64[] memory blocktimes = new uint64[](taskLen);

        // Query every pool balance
        for (uint256 i = 0; i < taskLen; i++) {
            result = IArbius(_arbius).solutions(_taskids[i]);
            blocktimes[i] = result.blocktime;
        }

        // ensure abi encoding, not needed here but increase reusability for different return types
        // note: abi.encode add a first 32 bytes word with the address of the original data
        bytes memory _abiEncodedData = abi.encode(blocktimes);

        assembly {
            // Return from the start of the data (discarding the original data address)
            // up to the end of the memory used
            let dataStart := add(_abiEncodedData, 0x20)
            return(dataStart, sub(msize(), dataStart))
        }
    }
}