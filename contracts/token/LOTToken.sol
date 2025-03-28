// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract LOTToken is ERC20, Ownable {
    constructor(uint256 initialSupply) ERC20("Lottery Token", "LOT") Ownable(msg.sender) {
        _mint(msg.sender, initialSupply);
    }
} 