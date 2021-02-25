package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"

	"modernc.org/mathutil"
)

func errchk(e error) {
	if e != nil {
		panic(e)
	}
}

var dataMap map[string](map[string]int)

func interpret(lines []string) {
	dataMap = make(map[string]map[string]int)
	for _, line := range lines {
		if line == "" {
			continue
		}
		words := strings.Split(line, " ")
		if dataMap[words[0]] == nil {
			dataMap[words[0]] = make(map[string]int)
		}
		val, err := strconv.Atoi(words[3])
		errchk(err)
		if words[2] == "lose" {
			val = -val
		}
		dataMap[words[0]][words[len(words)-1]] = val
	}
}

func keys(in map[string]map[string]int) []string {
	out := make([]string, 0, len(in))
	for k := range in {
		out = append(out, k)
	}
	return out
}

//quick and dirty fix for % being a shitty remainder and not a mod function
//!only for positiv b!
func mod(a int, b int) int {
	out := a % b
	if out < 0 {
		return out + b
	}
	return out
}

func calcHappieness(names []string) int {
	score := 0
	for i, name := range names {
		score += dataMap[name][names[mod(i-1, len(names))]]
		score += dataMap[name][names[mod(i+1, len(names))]]
	}
	return score
}

func findMaxHappieness() int {
	keySet := sort.StringSlice(keys(dataMap))
	mathutil.PermutationFirst(keySet)
	max, score := 0, 0
	for mathutil.PermutationNext(keySet) {
		score = calcHappieness(keySet)
		if score > max {
			max = score
		}
	}
	return max
}

func main() {
	dat, err := ioutil.ReadFile("input")
	errchk(err)
	lines := strings.Split(string(dat), ".\n")
	interpret(lines)

	fmt.Printf("Part 1: %d\n", findMaxHappieness())
	keySet := keys(dataMap)
	dataMap["me"] = make(map[string]int)
	for _, name := range keySet {
		dataMap[name]["me"] = 0
		dataMap["me"][name] = 0
	}
	fmt.Printf("Part 2: %d\n", findMaxHappieness())
}
