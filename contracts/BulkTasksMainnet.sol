// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;
import "./interfaces/IArbius.sol";

contract BulkTasks {

    constructor() {
        IBaseToken(0x4a24B101728e07A52053c13FB4dB2BcF490CAbc3).approve(0x9b51Ef044d3486A1fB0A2D55A6e0CeeAdd323E66, type(uint256).max);
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
                    0x9b51Ef044d3486A1fB0A2D55A6e0CeeAdd323E66,
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

            let input_data := mload(0x40) 
                  
            mstore(input_data, 0x506ea7de00000000000000000000000000000000000000000000000000000000) 

            let guard := add(1, calldatasize())
            for {let offset := commitments_.offset} lt(offset, guard) {offset := add(offset, 32)} {
                calldatacopy(add(input_data, 4), offset, 32)
                pop(call(
                    gas(),
                    0x9b51Ef044d3486A1fB0A2D55A6e0CeeAdd323E66,
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

