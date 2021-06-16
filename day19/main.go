package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

// so brute force was way to complicated.
// with a bit of reddit I got hint to the pattern:
// https://www.reddit.com/r/adventofcode/comments/3xflz8/day_19_solutions/

//error checking function
//causes panic if error occurs
func errchk(e error) {
	if e != nil {
		panic(e)
	}
}

//interpret the input from the file into usable dataStructure (and send startString to Channel)
func interpret(lines []string, startStr chan<- string) map[string][]string {
	trans := make(map[string][]string)
	for _, line := range lines {
		if line != "" {
			if strings.Contains(line, " => ") {
				in := strings.Split(line, " => ")
				if _, ok := trans[in[0]]; ok {
					trans[in[0]] = append(trans[in[0]], in[1])
				} else {
					trans[in[0]] = []string{in[1]}
				}
			} else {
				startStr <- line
				close(startStr)
			}
		}
	}
	return trans
}

//count the amount of transformations possibel on the in string useing the t(ransition) map
//for each key in the map there are t[key] possible results
//duplicat filtering is problematic
func countTransformations(in string, t map[string][]string) int {
	transmuted := make(map[string]int, 0)
	for key, val := range t {
		re := regexp.MustCompile(key)
		ind := re.FindAllStringIndex(in, -1)
		for _, i := range ind {
			for _, rep := range val {
				tmp := in[:i[0]] + rep + in[i[0]+len(key):]
				transmuted[tmp] = 0
			}
		}
	}
	return len(transmuted)
}

func makeMed(med string, t map[string][]string, res chan<- int) {
	//treating Rn, Y, Ar as special
	//Y needs to be counted twice
	// I used special chars as replacements as they can be easily distingushed from keys
	med = strings.ReplaceAll(med, "Rn", "")
	med = strings.ReplaceAll(med, "Ar", "")
	med = strings.ReplaceAll(med, "Y", "−")
	for key := range t {
		med = strings.ReplaceAll(med, key, "_")
	}
	pos := regexp.MustCompile("_")
	neg := regexp.MustCompile("−")
	pMatches := pos.FindAllStringIndex(med, -1)
	nMatches := neg.FindAllStringIndex(med, -1)
	pNum := len(pMatches)
	nNum := len(nMatches)
	res <- pNum - nNum - 1
}

func main() {
	dat, err := ioutil.ReadFile("input")
	errchk(err)
	lines := strings.Split(string(dat), "\n")
	var med string
	strChannel := make(chan string)
	transitions := make(map[string][]string)
	go func() {
		transitions = interpret(lines, strChannel)
	}()
	med = <-strChannel

	res1, res2 := make(chan int), make(chan int)
	go func() {
		res1 <- countTransformations(med, transitions)
		makeMed(med, transitions, res2)
	}()
	fmt.Printf("Part 1: %d\n", <-res1)
	fmt.Printf("Part 2: %d\n", <-res2)
}
