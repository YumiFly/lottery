// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {VRFConsumerBaseV2Plus} from "./VRFConsumerBaseV2Plus.sol";
import {VRFV2PlusClient} from "./VRFV2PlusClient.sol";

contract Simple is VRFConsumerBaseV2Plus {
    uint256 s_subscriptionId;
    address public vrfCoordinator; // 使用 public 变量
    bytes32 s_keyHash = 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80;
    uint32 callbackGasLimit = 40000;
    uint16 requestConfirmations = 3;
    uint32 numWords =  1;

    mapping(uint256 => address) public s_rollers;
    mapping(address => uint256) public s_results;

    // variables
    uint256 private constant ROLL_IN_PROGRESS = 0;
    // events
    event DiceRolled(uint256 indexed requestId, address indexed roller);
    event DiceLanded(uint256 indexed requestId, uint256 indexed result);

    // constructor
    constructor(uint256 subscriptionId, address _vrfCoordinator) VRFConsumerBaseV2Plus(_vrfCoordinator) { // 接收 vrfCoordinator 地址
        s_subscriptionId = subscriptionId;
        vrfCoordinator = _vrfCoordinator; // 设置 vrfCoordinator 地址
    }

    // rollDice function
    function rollDice(address roller) public onlyOwner returns (uint256 requestId) {
        //require(s_results[roller] == 0, "Already rolled");
        // Will revert if subscription is not set and funded.
        
       requestId = s_vrfCoordinator.requestRandomWords(
            VRFV2PlusClient.RandomWordsRequest({
                keyHash: s_keyHash,
                subId: s_subscriptionId,
                requestConfirmations: requestConfirmations,
                callbackGasLimit: callbackGasLimit,
                numWords: numWords,
                // Set nativePayment to true to pay for VRF requests with Sepolia ETH instead of LINK
                extraArgs: VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: false}))
            })
        );

        s_rollers[requestId] = roller;
        s_results[roller] = ROLL_IN_PROGRESS;
        emit DiceRolled(requestId, roller);
    }

    // fulfillRandomWords function
    function fulfillRandomWords(uint256 requestId, uint256[] calldata randomWords) internal override {

        // transform the result to a number between 1 and 2000 inclusively
        uint256 d20Value = randomWords[0] % 2000 + 1;

        // assign the transformed value to the address in the s_results mapping variable
        s_results[s_rollers[requestId]] = d20Value;

        // emitting event to signal that dice landed
        emit DiceLanded(requestId, d20Value);
    }
}