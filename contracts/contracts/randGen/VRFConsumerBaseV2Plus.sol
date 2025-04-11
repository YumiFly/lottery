// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {VRFV2PlusClient} from "./VRFV2PlusClient.sol";

interface IVRFCoordinatorV2Plus {
  function requestRandomWords(VRFV2PlusClient.RandomWordsRequest calldata req) external returns (uint256 requestId);
}

abstract contract VRFConsumerBaseV2Plus {
    error OnlyCoordinatorCanFulfill(address have, address want);
    address private owner;
    IVRFCoordinatorV2Plus public s_vrfCoordinator;

    constructor(address _vrfCoordinator) {
        owner = msg.sender;
        s_vrfCoordinator = IVRFCoordinatorV2Plus(_vrfCoordinator);
    }

    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this function");
        _;
    }

    function fulfillRandomWords(uint256 requestId, uint256[] calldata randomWords) internal virtual;

    function rawFulfillRandomWords(uint256 requestId, uint256[] calldata randomWords) external {
        if (msg.sender != address(s_vrfCoordinator)) {
            revert OnlyCoordinatorCanFulfill(msg.sender, address(s_vrfCoordinator));
        }
        fulfillRandomWords(requestId, randomWords);
    }
}