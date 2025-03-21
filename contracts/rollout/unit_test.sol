// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IRolloutCallback} from "https://raw.githubusercontent.com/YumiFly/lottery/refs/heads/main/contracts/interface/rollout_if.sol";

contract UnitTest is IRolloutCallback {
    uint256[] public results;
    uint256 public callbackCount;

    function rolloutCallback(uint256[] calldata _results) external override {
        results = _results;
        callbackCount++;
    }
}