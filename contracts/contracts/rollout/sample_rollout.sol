// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {VRFConsumerBaseV2Plus} from "../randGen/VRFConsumerBaseV2Plus.sol";
import {VRFV2PlusClient} from "../randGen/VRFV2PlusClient.sol";
import {IRolloutCall} from "../interface/rollout_if.sol";
import {IRolloutCallback} from "../interface/rollout_if.sol";

contract SimpleRollout is VRFConsumerBaseV2Plus, IRolloutCall {
    uint256 s_subscriptionId;
    IRolloutCallback s_callbacks;
    uint256 public rollout_epoch;
    uint256[] public rollout_results;
    address trigger;
    uint256 public requestID;

 
    // events
    event DiceRolled(uint256 indexed requestId, uint256 indexed epoch);
    event DiceLanded(uint256 indexed requestId, uint256[] indexed results);


    modifier onlyTrigger() {
        require(msg.sender == trigger, "Only trigger can call this function");
        _;
    }

    // constructor
    constructor(uint256 subscriptionId, address _vrfCoordinator, address _trigger) VRFConsumerBaseV2Plus(_vrfCoordinator) {
        s_subscriptionId = subscriptionId;
        rollout_epoch = 1;
        trigger = _trigger;
    }

    // rollDice function
    function rollDice() private returns (uint256 requestId) {
       requestId = s_vrfCoordinator.requestRandomWords(
            VRFV2PlusClient.RandomWordsRequest({
                keyHash: 0x787d74caea10b2b357790d5b5247c2f63d1d91572a9846f780606e4d953677ae,
                subId: s_subscriptionId,
                requestConfirmations: 3,
                callbackGasLimit: 40000,
                numWords: 3,
                // Set nativePayment to true to pay for VRF requests with Sepolia ETH instead of LINK
                extraArgs: VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: false}))
            })
        );
    }

    function rolloutCall(IRolloutCallback rolloutcb) external onlyTrigger override {
        delete rollout_results;
        requestID = rollDice();
        // store the callback
        s_callbacks = rolloutcb;
        emit DiceRolled(requestID, rollout_epoch);
        rollout_epoch++;
    }

    // fulfillRandomWords function
    function fulfillRandomWords(uint256 requestId, uint256[] calldata randomWords) internal override {
        for (uint i = 0; i < randomWords.length; i++) {
            rollout_results.push(randomWords[i] % 36 + 1);
        }
        s_callbacks.rolloutCallback(rollout_results);
        emit DiceLanded(requestId, rollout_results);
    }
}