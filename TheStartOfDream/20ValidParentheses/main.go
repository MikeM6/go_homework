package main

import (
	"fmt"
)

// Given a string s containing just the characters '(', ')', '{', '}', '[' and ']', determine if the input string is valid.

// An input string is valid if:

// 1.Open brackets must be closed by the same type of brackets.
// 2.Open brackets must be closed in the correct order.
// 3.Every close bracket has a corresponding open bracket of the same type.

// func isValid(s string) bool {

// }

func main() {
	samples := []string{"", "()", "([)]", "([]{})", "{[(])}", "(((())))"}
	for _, s := range samples {
		fmt.Printf("%q -> %v\n", s, isValid(s))
	}
}

// use array simulate stack
func isValid(s string) bool {
	stack := []rune{}
	pairsMap := map[rune]rune{
		')': '(',
		'[': ']',
		'{': '}',
	}
	for _, c := range s {
		switch c {
		case '(', '[', '{':
			stack = append(stack, c)
		case ')', '}', ']':
			if len(stack) == 0 || stack[len(stack)-1] != pairsMap[c] {
				return false
			}
			stack = stack[:len(stack)-1]
		default:
			return false

		}
	}
	return true
}
