// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;
import "./interfaces/IArbius.sol";

contract BulkTasks {

    constructor() {
        IBaseToken(0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0).approve(0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9, type(uint256).max);
    }

    function claimSolutions(bytes32[] calldata _taskids) public {
        assembly {
             let input_data := mload(0x40) 
                  
            mstore(input_data, 0x77286d1700000000000000000000000000000000000000000000000000000000) 

            let guard := add(1, calldatasize()) 
            for {let offset := _taskids.offset} lt(offset, guard) {offset := add(offset, 32)} {
                calldatacopy(add(input_data, 4), offset, 32)
                pop(call(
                    gas(),
                    0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9,
                    0,
                    input_data,
                    36,
                    0,
                    0
                ))
            }
        }
    }

    function bulkSignalCommitment(bytes32[] calldata commitments_) public {
         assembly {

            //0x506ea7decaf9a90acda804f400bf1238754038eaec53eebdcb1ea409029b0eff8afd00cf
            let input_data := mload(0x40) 
                  
            mstore(input_data, 0x506ea7de00000000000000000000000000000000000000000000000000000000) 

            let guard := add(1, calldatasize())
            for {let offset := commitments_.offset} lt(offset, guard) {offset := add(offset, 32)} {
                calldatacopy(add(input_data, 4), offset, 32)
                pop(call(
                    gas(),
                    0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9,
                    0,
                    input_data,
                    36,
                    0,
                    0
                ))
            }
        }
    }

} 

