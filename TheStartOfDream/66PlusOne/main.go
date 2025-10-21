package main

import "fmt"

// You are given a large integer represented as an integer array digits, where each digits[i] is the ith digit of the integer. The digits are ordered from most significant to least significant in left-to-right order. The large integer does not contain any leading 0's.

// Increment the large integer by one and return the resulting array of digits.
// func plusOne(digits []int) []int {

// }

func main() {
	fmt.Println(plusOne([]int{1, 2, 3}))    // [1 2 4]
	fmt.Println(plusOne([]int{1, 2, 9}))    // [1 3 0]
	fmt.Println(plusOne([]int{9, 9, 9}))    // [1 0 0 0]
	fmt.Println(plusOne([]int{9}))          // [1 0]
	fmt.Println(plusOne([]int{2, 9, 9}))    // [3 0 0]
	fmt.Println(plusOne([]int{4, 3, 2, 1})) // [4 3 2 2]
	fmt.Println(plusOne([]int{}))           // [1]
}

func plusOne(digits []int) []int {
	if len(digits) == 0 {
		return []int{1}
	}
	for i := len(digits) - 1; i >= 0; i-- {
		if digits[i] < 9 {
			digits[i]++
			return digits
		}
		digits[i] = 0
	}
	result := make([]int, len(digits)+1)
	result[0] = 1
	return result
}
