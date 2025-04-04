// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;
import "./interfaces/IArbius.sol";

contract BulkTasks {
    address public basetoken;
    address public engine;

    constructor(address basetoken_, address engine_) {
        basetoken = basetoken_;
        engine = engine_;
        IBaseToken(basetoken).approve(engine, type(uint256).max);
    }

    function claimSolutions(bytes32[] calldata _taskids) public {
        address engine_ = engine;
        assembly {
            let input_data := mload(0x40)

            mstore(
                input_data,
                0x77286d1700000000000000000000000000000000000000000000000000000000
            )

            let guard := add(1, calldatasize())
            for {
                let offset := _taskids.offset
            } lt(offset, guard) {
                offset := add(offset, 32)
            } {
                calldatacopy(add(input_data, 4), offset, 32)
                pop(call(gas(), engine_, 0, input_data, 36, 0, 0))
            }
        }
    }

    function bulkSignalCommitment(bytes32[] calldata commitments_) public {
        address engine_ = engine;

        assembly {
            let input_data := mload(0x40)

            mstore(
                input_data,
                0x506ea7de00000000000000000000000000000000000000000000000000000000
            )

            let guard := add(1, calldatasize())
            for {
                let offset := commitments_.offset
            } lt(offset, guard) {
                offset := add(offset, 32)
            } {
                calldatacopy(add(input_data, 4), offset, 32)
                pop(
                    call(
                        gas(),
                        engine_,
                        0,
                        input_data,
                        36,
                        0,
                        0
                    )
                )
            }
        }
    }
}
