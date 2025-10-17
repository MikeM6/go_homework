package main

// Write a function to find the longest common prefix string amongst an array of strings.

// If there is no common prefix, return an empty string "".

// func longestCommonPrefix(strs []string) string {

// }

func main() {
	println(longestCommonPrefix([]string{"flower", "flow", "flight"})) // fl
	println(longestCommonPrefix([]string{"dog", "racecar", "car"}))    // ""
	println(longestCommonPrefix([]string{}))                           // ""
	println(longestCommonPrefix([]string{"cir", "car"}))               // "c"
	println(longestCommonPrefix([]string{"inter", "internet"}))        // "inter"
	println(longestCommonPrefix([]string{"", ""}))
}

func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	for i := 0; i < len(strs[0]); i++ {
		c := strs[0][i]
		for j := 1; j < len(strs); j++ {
			// the other str is end || the other substr not match
			if i == len(strs[j]) || c != strs[j][i] {
				return strs[0][:i]
			}
		}
	}
	return strs[0]
}
