package main

import (
	"fmt"
	"sort"
)

// Given an array of intervals where intervals[i] = [starti, endi],
// merge all overlapping intervals, and return an array of the non-overlapping
// intervals that cover all the intervals in the input.

// func merge(intervals [][]int) [][]int {

// }

func main() {
	fmt.Println(merge([][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}})) // [[1 6] [8 10] [15 18]]
	fmt.Println(merge([][]int{{1, 4}, {4, 5}}))                    // [[1 5]]
}

func merge(intervals [][]int) [][]int {
	if len(intervals) == 0 {
		return [][]int{}
	}

	sort.Slice(intervals, func(i, j int) bool {
		if intervals[i][0] == intervals[j][0] {
			return intervals[i][1] < intervals[j][1]
		}
		return intervals[i][0] < intervals[j][0]
	})

	cur := []int{intervals[0][0], intervals[0][1]}
	res := make([][]int, 0, len(intervals))
	for _, v := range intervals[1:] {
		if cur[1] >= v[0] { // overlap
			if cur[1] < v[1] {
				cur[1] = v[1]
			}
		} else {
			res = append(res, cur)
			cur = []int{v[0], v[1]}
		}
	}
	res = append(res, cur)
	return res
}
