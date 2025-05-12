// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IRolloutCallback} from "./interface/rollout_if.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract LotteryManager is IRolloutCallback {
    address public admin;
    address public owner;
    address public rolloutContract;
    uint256 public prizeRate;
    IERC20 public tokenContract;
    uint256 private validCount;
    string public name;
    uint256 public totalSupply;
    uint256 private g_TotalSupply;
    uint256 public price;   // Token price of each bet
    uint256 public epoch;

    
    struct Bet {
        address buyer;
        uint256 amount;
    }
    mapping(bytes => Bet[]) public betAmounts;
    bytes[] public allBets;

    enum ContractState { Ready, Distribute, Rollout, Terminal }
    ContractState private state;


    event LotteryResults(uint256[] results, uint256 epoch, uint256 timestamp);
    event RolloutCallbakTXFailed(address, address, uint256);
    event TransState(ContractState);

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

    modifier validPlaceBet(uint256[] calldata target) {
        require(target.length == validCount, "Invalid target length");
        _;
    }

    modifier onlyToken() {
        require(msg.sender == address(tokenContract), "Only token can call this function");
        _;
    }

    /**
     * @notice 构造函数，初始化合约的基本信息
     * @param _admin 管理员地址
     * @param _owner 拥有者地址
     * @param _rolloutContract 开奖合约地址
     * @param _name 彩票名称
     * @param supply 彩票总供应量
     * @param _price 每次投注的价格
     * @param _tokenContract 代币合约地址
     */
    constructor(address _admin, address _owner, address _rolloutContract, 
                string memory _name, uint256 supply, uint256 _price, address _tokenContract
            ) {
        admin = _admin;
        owner = _owner;
        rolloutContract = _rolloutContract;
        state = ContractState.Ready;
        name = _name;
        totalSupply = supply;
        g_TotalSupply = supply;
        price = _price;
        tokenContract = IERC20(_tokenContract);
        validCount = 3;
        epoch = 1;
        prizeRate = 50;
    }

    /**
     * @notice 处理开奖回调函数
     * @dev 该函数由开奖合约调用，用于处理开奖结果并分配奖金
     * @param _results 开奖结果数组，长度必须与 validCount 一致
     */
    function rolloutCallback(uint256[] calldata _results) external onlyRolloutContract override {
        require(state == ContractState.Rollout, "Contract is not in Rollout state");
        require(_results.length == validCount, "Invalid results length");

        bytes memory resultsAsBytes = abi.encode(_results);
        Bet[] memory bets = betAmounts[resultsAsBytes];
        //require(bets.length > 0, "No bets found for the given results");
        if (bets.length > 0) {
            uint256 sum = 0;
            for (uint i = 0; i < bets.length; i++) {
                sum += bets[i].amount;
            }
            uint256 prizePool = tokenContract.balanceOf(address(this));
            uint256 perPrize = prizePool * prizeRate / 100 / sum;
            require((perPrize * sum) < prizePool, "There is no enough token in Prize pool.");

            for (uint i = 0; i < bets.length; i++) {
                uint256 amount = perPrize * bets[i].amount;
                if (!tokenContract.transfer(bets[i].buyer, amount)){
                    emit RolloutCallbakTXFailed(address(this), bets[i].buyer, amount);
                }
            }
        }
        emit LotteryResults(_results, epoch, block.timestamp); // 触发事件
        setState(ContractState.Ready);
    }

    /**
     * @notice 用户投注记录函数
     * @dev 用户可以通过此函数进行投注，投注信息会存储在 betAmounts 中
     * @param _amount 投注数量
     * @param _target 投注目标，字节数组
     * @return 返回投注目标
     */
    function recordPlaceBet(address buyer, uint256 _amount, uint256[] calldata _target) external validPlaceBet(_target) onlyToken returns (uint256) {
        require(state == ContractState.Distribute, "Contract is not in Distribute state");
        require(totalSupply > 0, "All tokens have been sold out.");
 
        uint256 total_price = _amount * price;
        bool Notfound = true;
        bytes memory target = abi.encode(_target);
        Bet[] memory bets = betAmounts[target];
        for(uint i =0; i < bets.length; i++) {
            if (bets[i].buyer == buyer) {
                betAmounts[target][i].amount += _amount;
                Notfound = false;
                break;
            }
        }
        if (Notfound) {
            betAmounts[target].push(Bet({
                    buyer: buyer,
                    amount: _amount
                }));
            allBets.push(target);
        }
        // 因为 Solidity 有溢出检查，所以totalSupply只能用减法，不能用加法
        totalSupply -= _amount;

        return total_price;
    }

    /**
     * @notice 清理合约数据
     * @dev 清理奖池中的代币并删除所有投注记录
     */
    function clear() private {
        if (tokenContract.balanceOf(address(this)) > 0) {
            tokenContract.transfer(owner, tokenContract.balanceOf(address(this)));
        }
        
        if (allBets.length > 0) {
            for (uint i = 0; i < allBets.length; i++) {
                delete betAmounts[allBets[i]];
            }
            delete allBets;
        }

        totalSupply = g_TotalSupply;
    }

    /**
     * @notice 销毁合约
     * @dev 只有管理员可以调用，且合约状态必须为 Terminal 或 Ready
     */
    function destroy() external onlyAdmin {
        require(state == ContractState.Terminal || state == ContractState.Ready, "Contract is not in Terminal or Ready state");
        doDestroy();
    }

    function doDestroy() private {
        clear();
        selfdestruct(payable(admin));
    }

    /**
     * @notice 重置合约状态为 Ready
     * @dev 只有管理员可以调用，且合约状态必须为 Terminal
     */
    function doReady() private {
        require(state == ContractState.Rollout, "Contract is not in Rollout state");
        clear();
        epoch++;
    }

    /**
     * @notice 转换合约状态
     * @dev 只有管理员可以调用
     * @param _state 要设置的目标状态
     */
    function transState(ContractState _state) external onlyAdmin {
        setState(_state);
    }

    /**
     * @notice 设置合约状态
     * @dev 内部函数，用于限制状态转换规则
     * @param _state 要设置的目标状态
     */
    function setState(ContractState _state) private {
        ContractState m_state = state;
        if (m_state == ContractState.Ready && _state == ContractState.Distribute) {
            // ready -> distribute    
        } else if (m_state == ContractState.Distribute && _state == ContractState.Rollout) {
            // distribute -> rollout
        } else if (m_state == ContractState.Rollout && _state == ContractState.Ready) {
            // rollout -> ready
            doReady();
        } else if (m_state == ContractState.Ready && _state == ContractState.Terminal) {
            // ready -> terminal
            doDestroy();
        } else {
            revert("Invalid state transition");
        }
        state = _state;
        emit TransState(_state);
    }

    /**
     * @notice 获取当前合约状态
     * @return 返回当前合约状态
     */
    function getState() external view returns (ContractState) {
        return state;
    }
}