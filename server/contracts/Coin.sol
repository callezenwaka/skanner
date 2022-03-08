//SPDX-License-Identifier: Unlicense
pragma solidity ^0.8.3;

// import "hardhat/console.sol";

contract Coin {
  string public name = "PeaceCoin";
  string public symbol = "PC";
  uint public totalSupply = 1000000;
  address public owner;
  mapping(address => uint) balances;

  constructor() {
    balances[msg.sender] = totalSupply;
    owner = msg.sender;
  }

  function transfer(address receiver, uint amount) external {
    require(balances[msg.sender] >= amount, "Not enough coins");
    balances[msg.sender] -= amount;
    balances[receiver] += amount;
  }

  function balanceOf(address _owner) external view returns (uint) {
    return balances[_owner];
  }
}