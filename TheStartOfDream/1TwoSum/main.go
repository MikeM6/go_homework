package main

import "fmt"

// Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.

// You may assume that each input would have exactly one solution, and you may not use the same element twice.

// You can return the answer in any order.

// func twoSum(nums []int, target int) []int {

// }

func main() {
	fmt.Println(twoSum([]int{2, 7, 11, 15}, 9))        // [0 1]
	fmt.Println(twoSum([]int{3, 2, 4}, 6))             // [1 2]
	fmt.Println(twoSum([]int{3, 3}, 6))                // [0 1]
	fmt.Println(twoSum([]int{-1, -2, -3, -4, -5}, -8)) // [2 4]
	fmt.Println(twoSum([]int{0, 4, 3, 0}, 0))          // [0 3]
}

func twoSum(nums []int, target int) []int {
	recordMap := make(map[int]int, len(nums))
	for i, v := range nums {
		if mapValue, ok := recordMap[target-v]; ok {
			return []int{mapValue, i}
		}
		recordMap[v] = i
	}
	return nil
}
