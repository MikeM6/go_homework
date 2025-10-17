package main

import "fmt"

//Q
// Given a non-empty array of integers nums, every element appears twice except for one. Find that single one.

// You must implement a solution with a linear runtime complexity and use only constant extra space.

func main() {
	// Quick verification examples
	// fmt.Println(singleNumber([]int{2, 2, 1}))       // expect 1
	// fmt.Println(singleNumber([]int{4, 1, 2, 1, 2})) // expect 4
	// fmt.Println(singleNumber([]int{1}))             // expect 1

	fmt.Println(singleNumber2([]int{1, 1, 2, 3}))
}

func singleNumber1(nums []int) int {
	res := 0
	for _, n := range nums {
		res ^= n
	}
	return res
}

// func singleNumber2(nums []int) []int {
// 	freqMap := make(map[int]int, len(nums))
// 	for _, n := range nums {
// 		freqMap[n]++
// 	}
// 	res := []int{}
// 	for k, v := range freqMap {
// 		if v == 1 {
// 			res = append(res, k)
// 		}
// 	}
// 	return res
// }

func singleNumber2(nums []int) []int {
	freqMap := make(map[int]int, len(nums))
	res := []int{}

	for _, n := range nums {
		freqMap[n]++
	}
	for k, v := range freqMap {
		if v == 1 {
			res = append(res, k)
		}
	}
	return res
}
