// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {ERC20Permit} from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Permit.sol";
import {IERC20Permit} from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Permit.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "../lottery.sol";

contract FakeU is ERC20, ERC20Permit, Ownable {
    constructor(uint256 initialSupply) ERC20("Fake USD", "FKU")  ERC20Permit("Fake USD") Ownable(msg.sender) {
        _mint(msg.sender, initialSupply);
    }
}

contract LOTToken is ERC20, ERC20Permit, Ownable {
    // 稳定币信息结构体
    struct StablecoinInfo {
        string name;        // 稳定币名称
        uint256 rate;      // 兑换比例 (1个稳定币可兑换的LOT数量)
        address receiver;  // 稳定币收款地址
        bool isSupported;  // 是否支持该稳定币
    }

    bool private released;

    // 稳定币信息映射 (稳定币地址 => 稳定币信息)
    mapping(address => StablecoinInfo) public stablecoins;

    // 支持的稳定币地址列表
    address[] public supportedStablecoinsList;

    // 事件
    event StablecoinUpdated(address indexed stablecoin, string name, uint256 rate, address receiver);
    event StablecoinRemoved(address indexed stablecoin);
    event TokensExchanged(address indexed user, address indexed stablecoin, uint256 stablecoinAmount, uint256 lotAmount);
    event LOTExchanged(address indexed user, address indexed stablecoin, uint256 lotAmount, uint256 stablecoinAmount);

    modifier onlyReleased() {
        require(released, "The function is not released yet.");
        _;
    }
    constructor(uint256 initialSupply) ERC20("Lottery Token", "LOT") ERC20Permit("Lottery Token") Ownable(msg.sender) {
        released = false;
        _mint(msg.sender, initialSupply);
    }

    /**
     * @dev 设置发布标记
     */
    function setReleased() external onlyOwner {
        released = true;
    }

    function getReleased() external view returns (bool) {
        return released;
    }

    /**
     * @dev 添加或更新稳定币信息
     * @param stablecoinAddress 稳定币地址
     * @param name 稳定币名称
     * @param rate 兑换比例 (1个稳定币可兑换的LOT数量)
     * @param receiver 稳定币收款地址
     */
    function setStablecoin(address stablecoinAddress, string calldata name, uint256 rate, address receiver) external onlyOwner {
        require(stablecoinAddress != address(0), "Invalid stablecoin address");
        require(bytes(name).length > 0, "Name cannot be empty");
        require(rate > 0, "Rate must be greater than 0");
        require(receiver != address(0), "Invalid receiver address");

        // 如果是新的稳定币，添加到支持列表
        if (!stablecoins[stablecoinAddress].isSupported) {
            supportedStablecoinsList.push(stablecoinAddress);
            stablecoins[stablecoinAddress].isSupported = true;
        }

        // 更新稳定币信息
        stablecoins[stablecoinAddress].name = name;
        stablecoins[stablecoinAddress].rate = rate;
        stablecoins[stablecoinAddress].receiver = receiver;

        emit StablecoinUpdated(stablecoinAddress, name, rate, receiver);
    }

    /**
     * @dev 移除稳定币支持
     * @param stablecoinAddress 要移除的稳定币地址
     */
    function removeStablecoin(address stablecoinAddress) external onlyOwner {
        require(stablecoins[stablecoinAddress].isSupported, "Stablecoin not supported");

        // 从支持列表中移除
        for (uint256 i = 0; i < supportedStablecoinsList.length; i++) {
            if (supportedStablecoinsList[i] == stablecoinAddress) {
                // 将最后一个元素移到当前位置，然后删除最后一个元素
                supportedStablecoinsList[i] = supportedStablecoinsList[supportedStablecoinsList.length - 1];
                supportedStablecoinsList.pop();
                break;
            }
        }

        // 清除映射中的支持标志
        stablecoins[stablecoinAddress].isSupported = false;

        emit StablecoinRemoved(stablecoinAddress);
    }

    /**
     * @dev 使用稳定币兑换LOT代币，支持permit授权
     * @param stablecoinAddress 稳定币地址
     * @param amount 稳定币数量
     * @param deadline 授权的截止时间（时间戳）
     * @param v 签名的 v 值
     * @param r 签名的 r 值
     * @param s 签名的 s 值
     * @return 兑换的LOT数量
     */
    function exchangeForLOT(
        address stablecoinAddress,
        uint256 amount,
        uint256 deadline,
        uint8 v,
        bytes32 r,
        bytes32 s
    ) external returns (uint256) {
        StablecoinInfo memory info = stablecoins[stablecoinAddress];
        require(info.isSupported, "Stablecoin not supported");
        require(amount > 0, "Amount must be greater than 0");

        // 使用 permit 授权
        IERC20Permit(stablecoinAddress).permit(msg.sender, address(this), amount, deadline, v, r, s);
        // 计算可兑换的LOT数量
        uint256 lotAmount = amount * info.rate;
    
        // 转移稳定币到指定的收款地址
        IERC20 stablecoinToken = IERC20(stablecoinAddress);
        require(stablecoinToken.transferFrom(msg.sender, info.receiver, amount), "Stablecoin transfer failed");
    
        // 铸造LOT代币给用户
        _mint(msg.sender, lotAmount);
    
        emit TokensExchanged(msg.sender, stablecoinAddress, amount, lotAmount);
    
        return lotAmount;
    }

    /**
     * @dev 获取支持的稳定币数量
     * @return 支持的稳定币数量
     */
    function getSupportedStablecoinsCount() external view returns (uint256) {
        return supportedStablecoinsList.length;
    }

    /**
     * @dev 使用LOT代币兑换稳定币
     * @param stablecoinAddress 稳定币地址
     * @param lotAmount LOT代币数量
     * @return 兑换的稳定币数量
     * 注意事项：
     *      1、 使用前需要 StablecoinInfo.receiver 地址对当前合约进行approve授权；
     *      2、 本产品正式发布前，此接口无法调用；
     *      3、 本产品正式发布后，需要调用Released接口设置发布标记；
     */
    function exchangeForStablecoin(address stablecoinAddress, uint256 lotAmount) external onlyReleased returns (uint256) {
        StablecoinInfo memory info = stablecoins[stablecoinAddress];
        require(info.isSupported, "Stablecoin not supported");
        require(lotAmount > 0, "Amount must be greater than 0");

        // 计算可兑换的稳定币数量
        // 例如：如果比例是 1:10 (1个稳定币 = 10 LOT)，那么 100 LOT = 10 稳定币
        uint256 stablecoinAmount = lotAmount / info.rate;
        require(stablecoinAmount > 0, "Exchange amount too small");

        // 检查用户是否有足够的LOT代币
        require(balanceOf(msg.sender) >= lotAmount, "Insufficient LOT balance");

        // 检查合约是否有足够的稳定币
        IERC20 stablecoinToken = IERC20(stablecoinAddress);
        address receiver = info.receiver;
        require(stablecoinToken.balanceOf(receiver) >= stablecoinAmount, "Insufficient stablecoin balance in contract");

        // 销毁用户的LOT代币
        _burn(msg.sender, lotAmount);

        // // 使用 permit 授权
        // IERC20Permit(stablecoinAddress).permit(receiver, address(this), stablecoinAmount, deadline, v, r, s);
        // 从收款地址转移稳定币到用户
        // 注意：收款地址必须先授权该合约调用transferFrom
        require(stablecoinToken.transferFrom(receiver, msg.sender, stablecoinAmount), "Stablecoin transfer failed");

        emit LOTExchanged(msg.sender, stablecoinAddress, lotAmount, stablecoinAmount);

        return stablecoinAmount;
    }

    /**
     * @dev 获取稳定币信息
     * @param stablecoinAddress 稳定币地址
     * @return 稳定币名称、兑换比例、收款地址和是否支持
     */
    function getStablecoinInfo(address stablecoinAddress) external view returns (string memory, uint256, address, bool) {
        StablecoinInfo memory info = stablecoins[stablecoinAddress];
        return (info.name, info.rate, info.receiver, info.isSupported);
    }

    function buy(address placeAddr, uint256 _amount, uint256[] calldata target) external returns (bool) {
        address buyer = msg.sender;

        if (buyer == placeAddr) {
            revert("Not allowed.");
        }
        uint256 total_price = LotteryManager(placeAddr).recordPlaceBet(buyer, _amount, target);
        if (balanceOf(buyer) < total_price) {
            revert("There is no enough token.");
        }
        if (!transfer(placeAddr, total_price)) {
            revert("Payment failed.");
        }
        return true;
    }
}