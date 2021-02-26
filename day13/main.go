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

func interpret(lines []string) map[string]map[string]int {
	data := make(map[string]map[string]int)
	for _, line := range lines {
		if line == "" {
			continue
		}
		words := strings.Split(line, " ")
		if data[words[0]] == nil {
			data[words[0]] = make(map[string]int)
		}
		val, err := strconv.Atoi(words[3])
		errchk(err)
		if words[2] == "lose" {
			val = -val
		}
		data[words[0]][words[len(words)-1]] = val
	}
	return data
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
	var dataMap, dataMap2 map[string]map[string]int
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		dataMap = interpret(lines)
		wg.Done()
	}()
	go func() {
		dataMap2 = interpret(lines)
		wg.Done()
	}()
	wg.Wait()
	dataMap2["me"] = make(map[string]int)
	for k, _ := range dataMap {
		dataMap2[k]["me"] = 0
		dataMap2["me"][k] = 0
	}

	res1, res2 := make(chan int), make(chan int)
	go findMaxHappieness(res1, dataMap)
	go findMaxHappieness(res2, dataMap2)

	fmt.Printf("Part 1: %d\n", <-res1)
	fmt.Printf("Part 2: %d\n", <-res2)
}
