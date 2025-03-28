// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IRolloutCallback} from "./interface/rollout_if.sol";

contract LotteryManager is IRolloutCallback {
    address public admin;
    address public owner;
    address public rolloutContract;
    address public tokenContract;
    uint256 private validCount;
    string public name;
    uint256 public totalSupply;
    uint256 public price;   // Token price of each bet
    
    struct Bet {
        address buyer;
        uint256 amount;
    }
    mapping(bytes => Bet[]) public betAmounts;

    enum ContractState { Ready, Distribute, Rollout, Terminal }
    ContractState private state;

    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this function");
        _;
    }

    modifier onlyAdmin() {
        require(msg.sender == admin, "Only admin can call this function");
        _;
    }
    
    modifier onlyRolloutContract() {
        require(msg.sender == rolloutContract, "Only rollout contract can call this function");
        _;
    }

    modifier validPlaceBet(bytes calldata target) {
        require(target.length == validCount, "Invalid target length");
        _;
    }
    
    constructor(address _admin, address _owner, address _rolloutContract, 
                string memory _name, uint256 supply, uint256 _price, address _tokenContract
            ) {
        admin = _admin;
        owner = _owner;
        rolloutContract = _rolloutContract;
        state = ContractState.Ready;
        name = _name;
        totalSupply = supply;
        price = _price;
        tokenContract = _tokenContract;
        validCount = 3;
    }

    function rolloutCallback(uint256[] calldata _results) external onlyRolloutContract override {
        require(state == ContractState.Rollout, "Contract is not in Rollout state");
        require(_results.length == validCount, "Invalid results length");

        uint8[3] memory results = [uint8(_results[0]), uint8(_results[1]), uint8(_results[2])];
        // 将 _results 转换为 bytes
        bytes memory resultsAsBytes = abi.encode(results);
    
        // 在 betAmounts 中找到对应的 buyer 和 amount
        Bet[] memory bets = betAmounts[resultsAsBytes];
        require(bets.length > 0, "No bets found for the given results");
    
        for (uint i = 0; i < bets.length; i++) {
            // 使用tokenContract合约，将赢家的Token转移到赢家的Token地址中
            
        }

        // 开奖结束，进入 Terminal 状态
        setState(ContractState.Terminal);
    }

    function placeBet(address buyer, uint256 _amount, bytes calldata target) external validPlaceBet(target) returns (bytes memory) {
        require(state == ContractState.Distribute, "Contract is not in Distribute state");
        
        uint256 total_price = _amount * price;
        // check balance of buyer >= total_price
        // Todo: implement the above logic

        // transfer total_price from buyer to this contract
        // Todo: implement the above logic
        
        betAmounts[target].push(Bet({
                buyer: buyer,
                amount: _amount
            }));

        return target;
    }

    function destroy() external onlyAdmin {
        require(state == ContractState.Terminal || state == ContractState.Ready, "Contract is not in Terminal or Ready state");
        selfdestruct(payable(admin));
    }

    function transState(ContractState _state) external onlyAdmin {
        setState(_state);
    }

    function setState(ContractState _state) private {
        ContractState m_state = state;
        if (m_state == ContractState.Ready && _state != ContractState.Distribute) {
            revert("when in Ready state, only Distribute state is allowed");
        } else if (m_state == ContractState.Distribute && _state != ContractState.Rollout) {
            revert("when in Distribute state, only Rollout state is allowed");
        } else if (m_state == ContractState.Rollout && _state != ContractState.Terminal) {
            revert("when in Rollout state, only Terminal state is allowed");
        } else if (m_state == ContractState.Terminal && _state != ContractState.Ready) {
            revert("when in Terminal state, only Ready state is allowed");
        }
        state = _state;
    }

    function getState() external view returns (ContractState) {
        return state;
    }
}