/*
@author: sk
@date: 2022/10/6
*/
package main

import (
	"sort"
)

/*
牌型：
单张，对子，三带一
顺子，双顺，3顺各带1    (都是至少5张)
炸弹   可以打上面的任意(上面的互不相关)
*/

// GetBestCards 从src 获取最小的大于tar的牌 放回nil代表没有结果
func GetBestCards(tar []int, src []int) []int {
	tar = SortCard(tar)
	src = SortCard(src)
	srcCounts := GetCount(src)
	res := DoGetBestCards(tar, srcCounts)
	// 如果最终结果为nil 且 对方不是炸弹  尝试一下炸弹
	if res == nil && ParseBomb(tar) == -1 {
		res = GetMinBomb(-1, srcCounts)
	}
	return res
}

func DoGetBestCards(tar []int, srcCounts map[int]int) []int {
	if len(tar) == 0 { //随便出
		// 出最长的
		res := make([]int, 0)
		temp := GetLongestStraight3(srcCounts)
		if len(temp) > len(res) {
			res = temp
		}
		temp = GetLongestStraight2(srcCounts)
		if len(temp) > len(res) {
			res = temp
		}
		temp = GetLongestStraight1(srcCounts)
		if len(temp) > len(res) {
			res = temp
		}
		if len(res) > 0 {
			return res
		}
		// 再尝试3带1 对子，单牌
		res = GetMinCard3(-1, srcCounts)
		if len(res) > 0 {
			return res
		}
		res = GetMinCard2(-1, srcCounts)
		if len(res) > 0 {
			return res
		}
		return GetMinCard1(-1, srcCounts)
	}
	value := ParseBomb(tar)
	if value != -1 { // 炸弹  只能对炸弹
		return GetMinBomb(value, srcCounts)
	}
	value = ParseCard1(tar)
	if value != -1 {
		return GetMinCard1(value, srcCounts)
	}
	value = ParseCard2(tar)
	if value != -1 {
		return GetMinCard2(value, srcCounts)
	}
	value = ParseCard3(tar)
	if value != -1 {
		return GetMinCard3(value, srcCounts)
	}
	value = ParseStraight1(tar)
	if value != -1 {
		return GetMinStraight1(value, len(tar), srcCounts)
	}
	value = ParseStraight2(tar)
	if value != -1 {
		return GetMinStraight2(value, len(tar)/2, srcCounts)
	}
	value = ParseStraight3(tar)
	if value != -1 {
		return GetMinStraight3(value, len(tar)/4, srcCounts)
	}
	return nil
}

func GetMinStraight3(value int, l int, counts map[int]int) []int {
	arr := GetMinStraight(value, l, counts, 3)
	if arr == nil {
		return nil
	}
	set := NewSet()
	res := make([]int, 0)
	for i := 0; i < len(arr); i++ {
		res = append(res, arr[i], arr[i], arr[i])
		set.Add(arr[i])
	}
	for num, count := range counts {
		if set.Has(num) {
			count -= 3
		}
		for i := 0; i < count; i++ {
			res = append(res, num)
			l--
			if l == 0 {
				return res
			}
		}
	}
	return nil
}

func GetMinStraight2(value int, l int, counts map[int]int) []int {
	res := GetMinStraight(value, l, counts, 2)
	if res == nil {
		return nil
	}
	return append(res, res...)
}

func GetMinStraight1(value int, l int, counts map[int]int) []int {
	return GetMinStraight(value, l, counts, 1)
}

func GetMinStraight(minValue int, l int, counts map[int]int, minCount int) []int {
	arr := make([]int, 0)
	for num, count := range counts {
		if num > minValue && count > minCount {
			arr = append(arr, num)
		}
	}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i] < arr[j]
	})
	res := make([]int, 0)
	last := arr[0] - 1
	for i := 0; i < len(arr); i++ {
		if arr[i] != last+1 {
			last = arr[i]
			res = make([]int, 0)
		}
		res = append(res, arr[i])
		if len(res) == l {
			return res
		}
	}
	return nil
}

func GetMinBomb(minValue int, counts map[int]int) []int {
	res := 2233
	for num, count := range counts {
		if count == 4 && num > minValue && num < res {
			res = num
		}
	}
	if res == 2233 {
		return nil
	}
	return []int{res, res, res, res}
}

func ParseBomb(nums []int) int {
	if len(nums) != 4 {
		return -1
	}
	value := nums[0]
	for i := 1; i < 4; i++ {
		if value != nums[i] {
			return -1
		}
	}
	return value
}

// 获取大于 min的最小单张
func GetMinCard1(min int, counts map[int]int) []int {
	res := GetMinCard(min, counts, 1)
	if res == -1 {
		return nil
	}
	return []int{res}
}

func GetMinCard(min int, counts map[int]int, minCount int) int {
	minValue := min
	res := 2233
	for num, count := range counts {
		if count >= minCount && num > minValue && num < res {
			res = num
		}
	}
	if res == 2233 {
		return -1
	}
	return res
}

// 获取大于 min的最小对子
func GetMinCard2(min int, counts map[int]int) []int {
	res := GetMinCard(min, counts, 2)
	if res == -1 {
		return nil
	}
	return []int{res, res}
}

// 获取大于 min的最小三带1
func GetMinCard3(min int, counts map[int]int) []int {
	res := GetMinCard(min, counts, 3)
	if res == -1 {
		return nil
	}
	arr := make([]int, 3)
	for i := 0; i < len(arr); i++ {
		arr[i] = res
	}
	for num, count := range counts {
		if num == res {
			count -= 3
		}
		if count > 0 {
			arr = append(arr, num)
			return arr
		}
	}
	return nil
}

// 最长3顺
func GetLongestStraight3(counts map[int]int) []int {
	temp := GetLongestStraight(counts, 3)
	if len(temp) < 2 {
		return nil
	}
	res := make([]int, 0)
	set := NewSet()
	for i := 0; i < len(temp); i++ {
		res = append(res, temp[i], temp[i], temp[i])
		set.Add(temp[i])
	}
	needNum := len(temp)
	for num, count := range counts {
		if set.Has(num) {
			count -= 3
		}
		for i := 0; i < count; i++ {
			needNum--
			res = append(res, num)
			if needNum == 0 {
				return res
			}
		}
	}
	return nil
}

// 最长双顺
func GetLongestStraight2(counts map[int]int) []int {
	res := GetLongestStraight(counts, 2)
	if len(res) < 3 {
		return nil
	}
	res = append(res, res...)
	return res
}

// 最长单顺
func GetLongestStraight1(counts map[int]int) []int {
	res := GetLongestStraight(counts, 1)
	if len(res) < 5 {
		return nil
	}
	return res
}

func GetLongestStraight(counts map[int]int, minCount int) []int {
	arr := make([]int, 0)
	for num, count := range counts {
		if count >= minCount {
			//if num == 0 { // A 应该接到 K 后面 见下文
			//	num = 13
			//}
			arr = append(arr, num)
		}
	}
	// 太短了
	if len(arr) < 2 {
		return nil
	}
	// 0 ~ 12   A ~ K   注意从2开始串  应该是  1 ～ 12 0  上面已经处理过了
	sort.Slice(arr, func(i, j int) bool {
		return arr[i] < arr[j]
	})
	res := make([]int, 0)
	current := make([]int, 0)
	lastNum := arr[0] - 1
	for i := 0; i < len(arr); i++ {
		if arr[i] != lastNum+1 {
			// 同样长 取最数字小的
			if len(current) > len(res) {
				res = current
			}
			current = make([]int, 0)
		}
		current = append(current, arr[i])
	}
	l := len(res)
	if l < 2 {
		return nil
	} // 恢复数字
	//if res[l-1] == 13 {
	//	res[l-1] = 0
	//}
	return res
}

// 牌src 是否大于 牌tar
func IsBigger(tar []int, src []int) bool {
	tar = SortCard(tar)
	src = SortCard(src)
	if !IsLegal(src) {
		return false
	}
	if len(tar) == 0 {
		return true
	}
	// 炸弹判断
	value1 := ParseBomb(tar)
	value2 := ParseBomb(src)
	if value2 != value1 { // 出的是炸弹 可以提前结束
		return value2 > value1
	} // 下面的是都不是炸弹都情况
	if len(tar) != len(src) {
		return false
	}
	// 单张判断
	value1 = ParseCard1(tar)
	value2 = ParseCard1(src)
	if value1 != -1 && value2 != -1 {
		return value2 > value1
	}
	// 对子判断
	value1 = ParseCard2(tar)
	value2 = ParseCard2(src)
	if value1 != -1 && value2 != -1 {
		return value2 > value1
	}
	// 三带1判断
	value1 = ParseCard3(tar)
	value2 = ParseCard3(src)
	if value1 != -1 && value2 != -1 {
		return value2 > value1
	}
	// 顺子判断
	value1 = ParseStraight1(tar)
	value2 = ParseStraight1(src)
	if value1 != -1 && value2 != -1 {
		return value2 > value1
	}
	// 连队判断
	value1 = ParseStraight2(tar)
	value2 = ParseStraight2(src)
	if value1 != -1 && value2 != -1 {
		return value2 > value1
	}
	// 飞机判断
	value1 = ParseStraight3(tar)
	value2 = ParseStraight3(src)
	if value1 != -1 && value2 != -1 {
		return value2 > value1
	}
	return false // 出的不是一类牌
}

func ParseStraight3(nums []int) int {
	value, l := ParseStraight(nums, 2)
	if l*4 != len(nums) || l < 2 {
		return -1
	}
	return value
}

func ParseStraight2(nums []int) int {
	value, l := ParseStraight(nums, 2)
	if l*2 != len(nums) || l < 3 {
		return -1
	}
	return value
}

func ParseStraight1(nums []int) int {
	value, l := ParseStraight(nums, 1)
	if l != len(nums) || l < 5 {
		return -1
	}
	return value
}

func ParseStraight(nums []int, minCount int) (int, int) {
	counts := GetCount(nums)
	arr := make([]int, 0)
	for num, count := range counts {
		if count >= minCount {
			arr = append(arr, num)
		}
	}
	if len(arr) < 2 {
		return -1, -1
	}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i] < arr[j]
	})
	for i := 1; i < len(arr); i++ {
		if arr[i-1] != arr[i]-1 {
			return -1, -1
		}
	}
	return arr[0], len(arr)
}

func ParseCard3(nums []int) int {
	if len(nums) != 4 {
		return -1
	}
	return ParseCard(nums, 3)
}

func ParseCard2(nums []int) int {
	if len(nums) != 2 {
		return -1
	}
	return ParseCard(nums, 2)
}

func ParseCard1(nums []int) int {
	if len(nums) != 1 {
		return -1
	}
	return ParseCard(nums, 1)
}

func ParseCard(nums []int, minCount int) int {
	res := GetCount(nums)
	for num, count := range res {
		if count >= minCount {
			return num
		}
	}
	return -1
}

func SortCard(nums []int) []int {
	sort.Slice(nums, func(i, j int) bool {
		return nums[i] < nums[j]
	})
	return nums
}

// 出的牌是否合法
func IsLegal(src []int) bool {
	if ParseBomb(src) != -1 {
		return true
	}
	if ParseCard1(src) != -1 {
		return true
	}
	if ParseCard2(src) != -1 {
		return true
	}
	if ParseCard3(src) != -1 {
		return true
	}
	if ParseStraight1(src) != -1 {
		return true
	}
	if ParseStraight2(src) != -1 {
		return true
	}
	if ParseStraight3(src) != -1 {
		return true
	}
	return false
}
