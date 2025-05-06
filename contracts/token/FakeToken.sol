// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {ERC20Permit} from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Permit.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract FakeU is ERC20, ERC20Permit, Ownable {
    constructor(uint256 initialSupply) ERC20("Fake USD", "FKU")  ERC20Permit("Fake USD") Ownable(msg.sender) {
        _mint(msg.sender, initialSupply);
    }
}