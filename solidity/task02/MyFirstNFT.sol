// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";


contract MyFirstNFT is ERC721, ERC721URIStorage, Ownable {
   
    uint256 private _nextTokenId; // 下一个 NFT 的 ID


    constructor() ERC721("MyFirstNFT", "MFN") Ownable(msg.sender) {
    }


    function mintNFT(address recipient, string memory tokenURI) public onlyOwner {
        require(recipient != address(0), "Cannot mint to the zero address");
        require(bytes(tokenURI).length != 0, "Token URI cannot be empty");
        uint256 tokenId = _nextTokenId++;
        _mint(recipient, tokenId);
        _setTokenURI(tokenId, tokenURI);
    }
}