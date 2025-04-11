// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {VRFV2PlusClient} from "./VRFV2PlusClient.sol";
import {IVRFCoordinatorV2Plus} from "./VRFConsumerBaseV2Plus.sol";

contract VRFCoordinatorV2 is IVRFCoordinatorV2Plus{
    uint256 private s_requestId;
    mapping(uint256 => address) public s_requesters;

    event RandomWordsRequested(uint256 requestId, address requester, uint256 subId, bytes32 keyHash, uint32 callbackGasLimit, uint16 requestConfirmations, uint32 numWords, bytes extraArgs);
    //event RandomWordsRequested(uint256 indexed requestId);
   
    function requestRandomWords(VRFV2PlusClient.RandomWordsRequest calldata req) public returns (uint256) {
        s_requestId++;
        s_requesters[s_requestId] = msg.sender;
        emit RandomWordsRequested(s_requestId, msg.sender, req.subId, req.keyHash, req.callbackGasLimit, req.requestConfirmations, req.numWords, req.extraArgs);
        //emit RandomWordsRequested(s_requestId);
        return s_requestId;
    }

    function CallFullfillRandomWords(uint256 requestId, uint256[] calldata randomWords) public {
        address requester = s_requesters[requestId];
        if (requester == address(0)) {
            return;
        }
        (bool success, ) = requester.call(abi.encodeWithSignature("rawFulfillRandomWords(uint256,uint256[])", requestId, randomWords));
        require(success, "FulfillRandomWords failed.");
    }
}
