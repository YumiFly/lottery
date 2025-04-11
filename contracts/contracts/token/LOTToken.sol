// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "../lottery.sol";

contract LOTToken is ERC20, Ownable {
    constructor(uint256 initialSupply) ERC20("Lottery Token", "LOT") Ownable(msg.sender) {
        _mint(msg.sender, initialSupply);
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