// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

// 1. 允许用户向合约地址发送以太币。
// 2. 记录每个捐赠者的地址和捐赠金额。
// 3. 允许合约所有者提取所有捐赠的资金。

contract BeggingContract{
    address public owner;

    // 记录每个捐赠者的捐赠金额
    mapping(address => uint256) public donations;
    // 所有捐赠者列表
    address[] public donors;

    uint256 public startTimeStamp;
    uint256 public duration;
   
    event Donation(address donor, uint256 amount);

    constructor(uint256 _duration) {
        owner = msg.sender; // 合约创建者为所有者
        startTimeStamp = block.timestamp; // 记录合约创建时间 单位为 秒
        duration = _duration; // 合约持续时间为365天
    }

    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can perform this action");
        _;
    }

    modifier onlyActive() {
        require(block.timestamp <= startTimeStamp + duration, "Contract is no longer active");
        _;
    }

    function containers(address addr) internal view returns (bool) {
        for (uint i = 0; i < donors.length; i++) {
            if (donors[i] == addr) {
                return true;
            }
        }
        return false;
    }

    function donate() external payable onlyActive {
        require(msg.value > 0, "Donation must be greater than 0");
        donations[msg.sender] += msg.value; // 记录捐赠金额
        if (!containers(msg.sender)) {
            donors.push(msg.sender); // 如果捐赠者不在列表中，添加到列表
        }
        emit Donation(msg.sender, msg.value); // 触发捐赠事件
    }

    // 允许用户向合约发送以太币
    receive() external payable {
        require(msg.value > 0, "Donation must be greater than 0");
        require(block.timestamp <= startTimeStamp + duration, "Contract is no longer active");
        donations[msg.sender] += msg.value; // 记录捐赠金额
        if (!containers(msg.sender)) {
            donors.push(msg.sender); // 如果捐赠者不在列表中，添加到列表
        }
        emit Donation(msg.sender, msg.value); // 触发捐赠事件
    }

    // 提取所有捐赠的资金到 msg.sender
    function withdraw() external onlyOwner payable {
        require(address(this).balance > 0, "No funds to withdraw");
        // 使用 transfer 方法将合约余额转给所有者
        // transfer ETH and revert on failure
        // payable(msg.sender).transfer(address(this).balance); // 将合约余额转给所有者

        // 使用send方法将合约余额转给所有者
        // send方法不会抛出异常，如果失败会返回false
        // bool success = payable(msg.sender).send(address(this).balance);
        // require(success, "Transfer failed"); // 确保转账成功

        // 使用call方法将合约余额转给所有者
        //  call : transfer ETH with data , return result of function and bool
        (bool success, ) = payable(msg.sender).call{value: address(this).balance}("");
        require(success, "Transfer failed"); // 确保转账成功
    }

    // 获取捐赠者的捐赠金额
    function getDonation(address donor) external view returns (uint256) {
        return donations[donor]; // 返回指定捐赠者的捐赠金额
    }

    // 获取前n名的捐赠者地址
    function getTopDonors(uint256 topn) external view returns (address[] memory) {
        address[] memory topDonors = new address[](topn);
        for (uint i = 0; i < topn && i < donors.length; i++) {
            //  获取 topDonors 中的最小捐赠金额及索引
            (uint256 minDonation, uint256 minIndex) = getMinDonation(topDonors);
            if (minDonation < donations[donors[i]]) {
                topDonors[minIndex] = donors[i];
            }
        }
        return topDonors;
    }

    function getMinDonation(address[] memory topDonors) internal view returns (uint256, uint256) {
        uint256 minDonation = type(uint256).max;
        uint256 minIndex = 0;
        for (uint256 i = 0; i < topDonors.length; i++) {
            if (donations[topDonors[i]] < minDonation) {
                minDonation = donations[topDonors[i]];
                minIndex = i;
            }
        }
        return (minDonation, minIndex);
    }
}