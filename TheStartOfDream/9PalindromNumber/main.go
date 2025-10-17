package main

import "fmt"

// Given an integer x, return true if x is a palindrome, and false otherwise.

// func isPalindrome(x int) bool {

// }

func main() {
	// quick examples
	fmt.Println(isPalindromeNumber1([]int{3, 1, 2, 1, 3})) // true
	fmt.Println(isPalindromeNumber2(-121))                 // false
	fmt.Println(isPalindromeNumber2(10))                   // false
	fmt.Println(isPalindromeNumber2(31213))
}

func isPalindromeNumber1(nums []int) bool {
	numsLen := len(nums)
	if numsLen == 0 {
		return false
	}
	i, j := 0, numsLen-1
	for i < j {
		if nums[i] != nums[j] {
			return false
		}
		i++
		j--
	}
	return true

}

func isPalindromeNumber2(num int) bool {
	if num < 0 {
		return false
	}
	// xxx0 never panlindrom
	if num > 0 && num%10 == 0 {
		return false
	}
	revered := 0
	for num > revered {
		revered = revered*10 + num%10
		num /= 10
	}
	return num == revered || num == revered/10
}
