package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

func errchk(e error) {
	if e != nil {
		panic(e)
	}
}

type group struct {
	rest int
	list []int
}

func sum(s []int) int {
	sum := 0
	for _, i := range s {
		sum += i
	}
	return sum
}
func prod(s []int) int {
	prod := 1
	for _, i := range s {
		prod *= i
	}
	return prod
}

func unique_combination(l, sum_local int, target int, group []int, slc []int, grp_chan chan<- []int) {
	// If a unique combination is found
	if sum_local == target {
		rtn := make([]int, len(group))
		copy(rtn, group)
		grp_chan <- rtn
		return
	}
	// For all other combinations
	for i := l; i < len(slc); i++ {
		// Check if the sum exceeds K
		if sum_local+slc[i] > target {
			continue
		}
		// Check if it is repeated or not
		if i > l && slc[i] == slc[i-1] {
			continue
		}
		// Recursive call
		unique_combination(i+1, sum_local+slc[i], target, append(group, slc[i]), slc, grp_chan)
	}
}

// Function to find all combination
// of the given elements

func combination(slc []int, target int, res chan<- []int) {
	// Sort the given elements
	less := func(i int, j int) bool {
		return i < j
	}
	sort.Slice(slc, less)

	group := make([]int, 0)

	unique_combination(0, 0, target, group, slc, res)
	close(res)
}

func main() {
	dat, err := ioutil.ReadFile("input")
	errchk(err)
	lines := strings.Split(string(dat), "\n")
	weights := make([]int, 0)
	for _, line := range lines {
		if line == "" {
			continue
		}
		w, err := strconv.Atoi(line)
		errchk(err)
		weights = append(weights, w)
	}
	target := [2]int{sum(weights) / 3, sum(weights) / 4}

	//channel for returning results
	res := []chan int{make(chan int), make(chan int)}
	for i := 0; i < 2; i++ {
		grp_chan := make(chan []int, 5)
		go combination(weights, target[i], grp_chan)
		group := <-grp_chan
		go func(grp_chan <-chan []int, res chan<- int) {
			for grp := range grp_chan {
				if len(grp) < len(group) {
					group = grp
				} else if len(grp) == len(group) && prod(grp) < prod(group) {
					group = grp
				}
			}
			res <- prod(group)
		}(grp_chan, res[i])
	}
	fmt.Printf("Part 1: %d\n", <-res[0])
	fmt.Printf("Part 2: %d\n", <-res[1])
}
