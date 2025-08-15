// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract MyERC20 {
    string public name;
    string public symbol;
    address public owner;

    uint256 private _totalSupply;

    mapping (address account => uint256 balance) private balances;

    mapping (address account => mapping(address spender => uint256)) public allowances;

    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);

    constructor(string memory _name, string memory _symbol) {
        name = _name;
        symbol = _symbol;
        owner = msg.sender;
    }

    function totalSupply() public view returns (uint256) {
        return _totalSupply;
    }

    function balanceOf(address account) public view returns (uint256) {
        return balances[account];
    }

    function transfer(address to, uint256 value) public returns (bool) {
        address from = msg.sender;
        require(balances[from] >= value, "Not enough balance");
        //  这里没检查溢出
        balances[from] -= value;
        balances[to] += value;
        emit Transfer(from, to, value);
        return true;
    }

    // approve
    function approve(address spender, uint256 value) public returns (bool) {
        require(spender != address(0), "Invalid spender address");
        require(balances[msg.sender] >= value, "Not enough balance");
        address _owner = msg.sender;
        allowances[_owner][spender] = value;
        emit Approval(_owner, spender, value);
        return true;
    }

    function allowance(address _owner, address spender) public view returns (uint256) {
        require(_owner != address(0), "Invalid owner address");
        require(spender != address(0), "Invalid spender address");
        return allowances[_owner][spender];
    }

    function transferFrom(address from, address to, uint256 value) public returns (bool) {
        require(from != address(0), "Invalid from address");
        require(to != address(0), "Invalid to address");
        require(balances[from] >= value, "Not enough balance");
        require(allowances[from][msg.sender] >= value, "Allowance exceeded");

        allowances[from][msg.sender] -= value;
        balances[from] -= value;
        balances[to] += value;

        emit Transfer(from, to, value);
        return true;
    }

    function mint(address account, uint256 value) public returns (bool) {
        require(msg.sender == owner, "Only owner can mint");
        require(account != address(0), "Invalid account address");

        _totalSupply += value;
        balances[account] += value;

        emit Transfer(address(0), account, value);
        return true;
    }
}