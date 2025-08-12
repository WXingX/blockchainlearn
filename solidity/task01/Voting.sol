// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract Voting {
    mapping(address => uint256) public votesMap;
    address[] public voteAddr;
    string public str = "abcde";
    mapping(bytes1 => int256) public RomanNum;
    string[] roman = ["M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"];
    int256[] values = [int256(1000), 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1];


    constructor() {
        RomanNum['M'] = 1000;
        RomanNum['D'] = 500;
        RomanNum['C'] = 100;
        RomanNum['L'] = 50;
        RomanNum['X'] = 10;
        RomanNum['V'] = 5;
        RomanNum['I'] = 1;
    }

    function vote(address candidateAddr,bool isYes) public {
        if (isYes) {
            votesMap[candidateAddr]++;
            if (!contains(candidateAddr)) {
                voteAddr.push(candidateAddr);
            }
        }
    }

    function contains(address candidateAddr) public view returns (bool) {
        for (uint256 i = 0; i < voteAddr.length; i++) {
            if (voteAddr[i] == candidateAddr) {
                return true;
            }
        }
        return false;
    }

    function getVotes(address candidateAddr) public view returns (uint256) {
        return votesMap[candidateAddr];
    }

    function resetVotes() public {
        for (uint256 i = 0; i < voteAddr.length; i++) {
           votesMap[voteAddr[i]] = 0;
        }
        
    }

    // 反转字符串
    function convert() public view returns(string memory) {
        bytes memory b = bytes(str);
        for (uint256 i = 0; i < b.length / 2; i++) {
            (b[i], b[b.length - 1 - i]) = (b[b.length - 1 - i], b[i]);
        }
        return string(b);
    }

    // 整数转罗马数字
    function intToRoman(int256 num) public view returns(string memory s) {
        for (uint256 i = 0; i < values.length; i++) {
            while (num >= values[i]) {
                s = string(abi.encodePacked(s, roman[i]));
                num -= values[i];
            }
        }
        return s;
    }

    // 罗马数字转整数
    function romanToInt(string calldata romans) public view returns(int256 sum) {
        bytes memory b = bytes(romans);
        uint256 len = b.length;
        for (uint256 i = 0; i < len - 1; i++) {
            if (RomanNum[b[i]] < RomanNum[b[i + 1]]) {
                sum -= RomanNum[b[i]];
            } else {
                sum += RomanNum[b[i]];
            }
        }
        return sum + RomanNum[b[len - 1]];
    }

    // 合并有序数组
    function meergeSortedArray(int256[] calldata nums1,int256[] calldata nums2) public pure returns(int256[] memory) {
        int256[] memory res = new int256[](uint256(nums1.length + nums2.length));
        uint256 i = 0;
        uint256 j = 0;
        uint256 k = 0;
        while (i < nums1.length && j < nums2.length) {
            if (nums1[i] < nums2[j]) {
                res[k++] = nums1[i++];
            } else {
                res[k++] = nums2[j++];
            }
        }
        while (i < nums1.length) {
            res[k++] = nums1[i++];
        }
        while (j < nums2.length) {
            res[k++] = nums2[j++];
        }
        return res;
    }

    //  二分查找
    function searchNum(int256[] calldata nums, int256 target) public pure returns(int256) {
        int256 index = -1;
        uint256 left = 0;
        uint256 right = nums.length - 1;
        for (; left <= right; ) {
            uint256 mid = (left + right) / 2;
            if (nums[mid] == target) {
                index = int256(mid);
                break;
            } else if (nums[mid] < target) {
                left = mid + 1;
            } else if (nums[mid] > target) {
                if (mid == 0) {
                    // 这里 如果是0 再减一 会溢出,会导致一直循环消耗gas
                    break;
                }
                right = mid - 1;
            }
        }

        return index;
    }
}