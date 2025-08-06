// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Counter {
    int256 public cnt;
    function add() public {
        cnt++;
    }

    function getCount() public view returns (int256){
        return cnt;
    }
}