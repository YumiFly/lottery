// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

interface IRolloutCall {
    function rolloutCall(IRolloutCallback rolloutcb) external;
}

interface IRolloutCallback {
    function rolloutCallback(uint256[] calldata results) external;
}