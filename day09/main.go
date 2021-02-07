package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func errchk(e error) {
	if e != nil {
		panic(e)
	}
}

type Distance struct {
	A, B string
	Dist int
}

//list of all locations and Distances to other Locations
var locations map[string][]Distance

//list of keys of locations
var keys []string

//There is some unknown behavior in this Code making it not always come to the correct solution.
//Running it about 5 Times gave me the correct awnsers. Not nice, but hey...
func main() {
	dat, err := ioutil.ReadFile("input")
	errchk(err)
	tmp := strings.Split(string(dat), "\n")
	locations = make(map[string][]Distance)
	for _, s := range tmp {
		if s == "" {
			continue
		}
		tmp3 := strings.Split(s, " = ")
		dist := tmp3[1]
		tmp4 := strings.Split(tmp3[0], " to ")
		orig := tmp4[0]
		dest := tmp4[1]

		distint, err := strconv.Atoi(dist)
		errchk(err)
		locations[orig] = append(locations[orig], Distance{A: orig, B: dest, Dist: distint})
		locations[dest] = append(locations[dest], Distance{A: dest, B: orig, Dist: distint})
	}
	keys = make([]string, 0, len(locations))
	for k := range locations {
		keys = append(keys, k)
	}
	lDest, mindist := getRoute(keys, lesser)
	fmt.Printf("%s\n", lDest)

	lDest, maxdest := getRoute(keys, greater)
	fmt.Printf("%s\n", lDest)
	fmt.Printf("Part 1: %d\n", mindist)
	fmt.Printf("Part 2: %d\n", maxdest)
}

//compare type compare to values
type comp func(int, int) bool

//Compare if a greater b
func greater(a int, b int) bool {
	return a > b
}

//Compare if a lesser than b
func lesser(a int, b int) bool {
	return a < b
}

//get Route between all points of 'todo' based on the distance as compared by 'fn'
func getRoute(todo []string, fn comp) (string, int) {
	rtn := -1
	lDest := ""
	if len(todo) == 1 {
		return todo[0], 0
	}
	for i, l := range todo {
		lastDest, dist := getRoute(remove(todo, i), fn)
		dist += getDist(lastDest, l)
		if fn(dist, rtn) || rtn == -1 {
			rtn = dist
			lDest = l
		}
	}
	return lDest, rtn
}

//Get Distance between inA and inB
func getDist(inA string, inB string) int {
	for _, dest := range locations[inA] {
		if dest.B == inB {
			return dest.Dist
		}
	}
	return -1
}

//Clone input slice and remove element at index
func remove(slice []string, index int) []string {
	new := append([]string{}, slice[:index]...)
	return append(new, slice[index+1:]...)
}
