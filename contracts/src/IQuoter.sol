// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

interface IQuoter {
    function profitLevel() external view returns (uint);
    function GetAIUSPrice() external returns(uint, uint);

    event ProfitLevelChanged(uint _profitLevel);
}