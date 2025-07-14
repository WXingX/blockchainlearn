package main

import (
	"fmt"
	"sort"
)

func main() {
	// arr := [9]int{1, 2, 3, 5, 5, 6, 3, 6, 2}
	// res := singleNumber(arr[:])
	// fmt.Println(res)

	// res := isPalindrome(-121)
	// fmt.Println(res)

	// s := "{}()[{}]"
	// res := isValid(s)
	// fmt.Println(res)

	// arr := [3]int{1, 1, 2}
	// arr := [10]int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	// res := removeDuplicates(arr[:])
	// fmt.Println(res)
}

// 给定一个非空整数数组，除了某个元素只出现一次以外，其余每个元素均出现两次。找出那个只出现了一次的元素。
// 可以使用 for 循环遍历数组，结合 if 条件判断和 map 数据结构来解决，例如通过 map 记录每个元素出现的次数，
// 然后再遍历 map 找到出现次数为1的元素。
func singleNumber(nums []int) int {
	var numsMap map[int]int = make(map[int]int)
	for _, v := range nums {
		_, ok := numsMap[v]
		if ok {
			delete(numsMap, v)
		} else {
			numsMap[v] = v
		}
	}
	var firstKey int
	for k, _ := range numsMap {
		firstKey = k
		break
	}
	return numsMap[firstKey]
}

// 回文数
func isPalindrome(x int) bool {
	if x < 0 {
		return false
	}
	temp, sum := x, 0
	for {
		var q = temp % 10
		temp = temp / 10
		if temp == 0 && q == 0 {
			break
		}
		sum = sum*10 + q
	}

	fmt.Println(x)
	fmt.Println(sum)
	if sum == x {
		return true
	} else {
		return false
	}
}

// 有效的括号  给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串，判断字符串是否有效
func isValid(s string) bool {
	if s == "" || s[0] == ')' || s[0] == '}' || s[0] == ']' {
		return false
	}

	var stack []byte
	bytes := []byte(s)
	for _, v := range bytes {
		if len(stack) == 0 {
			stack = append(stack, v)
		} else if (stack[len(stack)-1] == '(' && v == ')') || (stack[len(stack)-1] == '{' && v == '}') || (stack[len(stack)-1] == '[' && v == ']') {
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, v)
		}
	}

	if len(stack) == 0 {
		return true
	} else {
		return false
	}
}

// 最长公共前缀  查找字符串数组中的最长公共前缀
func longestCommonPrefix(strs []string) string {
	rsLen := len(strs)
	if rsLen == 0 {
		return ""
	} else if rsLen == 1 {
		return strs[0]
	}

	var commonPrefix string = ""
	for i := 0; i < len(strs[0]) && i < len(strs[1]); i++ {
		if strs[0][i] != strs[1][i] {
			break
		}
		commonPrefix += string(strs[0][i])
	}
	if len(commonPrefix) == 0 {
		return commonPrefix
	}

	for i := 2; i < rsLen; i++ {
		var tempPrefix string = ""
		for j := 0; j < len(strs[i]) && j < len(commonPrefix); j++ {
			if strs[i][j] != commonPrefix[j] {
				break
			}
			tempPrefix += string(strs[i][j])
		}
		if len(tempPrefix) == 0 {
			return ""
		}
		commonPrefix = tempPrefix
	}

	return commonPrefix
}

func longestCommonPrefix2(strs []string) string {
	for i := 0; i < len(strs[0]); i++ {
		ch := strs[0][i]
		for j := 1; j < len(strs); j++ {

			if i >= len(strs[j]) || ch != strs[j][i] {
				return strs[0][:i]
			}
		}
	}
	return strs[0]
}

// 加一  给定一个由整数组成的非空数组所表示的非负整数，在该数的基础上加一
func plusOne(digits []int) []int {
	var res []int
	var carry int = 0
	for i := len(digits) - 1; i >= 0; i-- {
		var sum int = 0
		if i == len(digits)-1 {
			sum = digits[i] + 1
		} else if carry == 0 {
			res = append(digits[0:i+1], res...)
			break
		} else {
			sum = digits[i] + carry
		}

		if sum >= 10 {
			sum -= 10
			carry = 1
		} else {
			carry = 0
		}
		res = append([]int{sum}, res...)

	}
	if carry == 1 {
		res = append([]int{1}, res...)
	}

	return res
}

//	删除有序数组中的重复项
//
// 给你一个有序数组 nums ，请你原地删除重复出现的元素，使每个元素只出现一次，返回删除后数组的新长度。
// 不要使用额外的数组空间，你必须在原地修改输入数组并在使用 O(1) 额外空间的条件下完成。
// 可以使用双指针法，一个慢指针 i 用于记录不重复元素的位置，一个快指针 j 用于遍历数组，
// 当 nums[i] 与 nums[j] 不相等时，将 nums[j] 赋值给 nums[i + 1]，并将 i 后移一位。
func removeDuplicates(nums []int) int {
	for i, j := 0, 1; i < len(nums) && j < len(nums); {
		if nums[i] == nums[j] {
			j++
		} else if nums[i] != nums[j] {
			if i == 0 {
				nums = nums[j-1:]
			} else {
				nums = append(nums[0:i], nums[j-1:]...)
			}
			i++
			j = i + 1
		}
	}
	fmt.Println(nums)
	return len(nums)
}

func removeDuplicates2(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	count := 1
	preNum := nums[0]
	for _, v := range nums {
		if v != preNum {
			preNum = v
			nums[count] = v
			count += 1
		}
	}
	return count
}

// 合并区间 以数组 intervals 表示若干个区间的集合，其中单个区间为 intervals[i] = [starti, endi] 。
// 请你合并所有重叠的区间，并返回 一个不重叠的区间数组，该数组需恰好覆盖输入中的所有区间
// 可以先对区间数组按照区间的起始位置进行排序，然后使用一个切片来存储合并后的区间，遍历排序后的区间数组，将当前区间与切片中最后一个区间进行比较，如果有重叠，则合并区间；如果没有重叠，则将当前区间添加到切片中。
func merge(intervals [][]int) [][]int {
	if len(intervals) == 0 {
		return [][]int{}
	}

	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	result := make([][]int, 0)
	result = append(result, intervals[0])

	for i := 1; i < len(intervals); i++ {
		last := result[len(result)-1]
		curr := intervals[i]

		if curr[0] <= last[1] {
			// 有重叠，合并区间
			if last[1] < curr[1] {
				last[1] = curr[1]
			}
		} else {
			// 无重叠，直接添加到结果中
			result = append(result, curr)
		}
	}

	return result
}

// 给定一个整数数组 nums 和一个目标值 target，请你在该数组中找出和为目标值的那两个整数
func twoSum(nums []int, target int) []int {
	for i := 0; i < len(nums)-1; i++ {
		var r = target - nums[i]
		for j := i + 1; j < len(nums); j++ {
			if nums[j] == r {
				return []int{i, j}
			}
		}
	}
	return []int{}
}

func twoSum2(nums []int, target int) []int {
	m := map[int]int{}
	for k, v := range nums {
		m[v] = k
	}

	for k, v := range nums {
		if k2, ok := m[target-v]; ok && k != k2 {
			return []int{k, k2}
		}
	}

	return []int{}
}
