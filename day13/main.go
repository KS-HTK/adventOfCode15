package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
	"sync"

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

func calcHappieness(res chan<- int, names []string, data map[string]map[string]int) {
	score := 0
	for i, name := range names {
		score += data[name][names[mod(i-1, len(names))]]
		score += data[name][names[mod(i+1, len(names))]]
	}
	res <- score
}

func findMaxHappieness(result chan<- int, data map[string]map[string]int) {
	keySet := sort.StringSlice(keys(data))
	mathutil.PermutationFirst(keySet)
	scores := make(chan int, 100)
	go func() {
		max := 0
		for score := range scores {
			if score > max {
				max = score
			}
		}
		result <- max
		close(result)
	}()
	var w sync.WaitGroup
	for {
		if !mathutil.PermutationNext(keySet) {
			break
		}
		tmp := make([]string, keySet.Len())
		for i, v := range keySet {
			tmp[i] = v
		}
		w.Add(1)
		go func() {
			calcHappieness(scores, tmp, data)
			w.Done()
		}()
	}
	w.Wait()
	close(scores)
}

func main() {
	dat, err := ioutil.ReadFile("input")
	errchk(err)
	lines := strings.Split(string(dat), ".\n")
	interpret(lines)
	res1, res2 := make(chan int), make(chan int)
	go findMaxHappieness(res1, dataMap)
	dataMap2 := make(map[string]map[string]int)
	dataMap2["me"] = make(map[string]int)
	for k, v := range dataMap {
		dataMap2[k] = make(map[string]int)
		for key, val := range v {
			dataMap2[k][key] = val
		}
		dataMap2[k]["me"] = 0
		dataMap2["me"][k] = 0
	}
	go findMaxHappieness(res2, dataMap2)

	fmt.Printf("Part 1: %d\n", <-res1)
	fmt.Printf("Part 2: %d\n", <-res2)
}
