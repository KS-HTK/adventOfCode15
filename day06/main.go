package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

func errchk(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	dat, err := ioutil.ReadFile("input")
	errchk(err)

	fmt.Printf("Part 1: %d\n", pt1(string(dat)))
	fmt.Printf("Part 2: %d\n", pt2(string(dat)))
}

func pt1(input string) int {
	var grid [1000][1000]bool
	command := strings.Split(input, "\n")
	re := regexp.MustCompile(`[0-9]*,[0-9]*`)
	for _, cmd := range command {
		if cmd == "" {
			continue
		}
		res := re.FindAll([]byte(cmd), -1)
		from := strings.Split(string(res[0]), ",")
		to := strings.Split(string(res[1]), ",")

		fromx, err := strconv.Atoi(from[0])
		errchk(err)
		fromy, err := strconv.Atoi(from[1])
		errchk(err)
		tox, err := strconv.Atoi(to[0])
		errchk(err)
		toy, err := strconv.Atoi(to[1])
		errchk(err)

		if strings.Contains(cmd, "toggle") {
			for i := fromx; i <= tox; i++ {
				for j := fromy; j <= toy; j++ {
					grid[i][j] = !grid[i][j]
				}
			}
		} else if strings.Contains(cmd, "off") {
			for i := fromx; i <= tox; i++ {
				for j := fromy; j <= toy; j++ {
					grid[i][j] = false
				}
			}
		} else {
			for i := fromx; i <= tox; i++ {
				for j := fromy; j <= toy; j++ {
					grid[i][j] = true
				}
			}
		}
	}
	acc := 0
	for i := 0; i < 1000; i++ {
		for j := 0; j < 1000; j++ {
			if grid[i][j] {
				acc++
			}
		}
	}
	return acc
}

func pt2(input string) int {
	var grid [1000][1000]int
	command := strings.Split(input, "\n")
	re := regexp.MustCompile(`[0-9]*,[0-9]*`)
	for _, cmd := range command {
		if cmd == "" {
			continue
		}
		res := re.FindAll([]byte(cmd), -1)
		from := strings.Split(string(res[0]), ",")
		to := strings.Split(string(res[1]), ",")

		fromx, err := strconv.Atoi(from[0])
		errchk(err)
		fromy, err := strconv.Atoi(from[1])
		errchk(err)
		tox, err := strconv.Atoi(to[0])
		errchk(err)
		toy, err := strconv.Atoi(to[1])
		errchk(err)

		if strings.Contains(cmd, "toggle") {
			for i := fromx; i <= tox; i++ {
				for j := fromy; j <= toy; j++ {
					grid[i][j] += 2
				}
			}
		} else if strings.Contains(cmd, "off") {
			for i := fromx; i <= tox; i++ {
				for j := fromy; j <= toy; j++ {
					grid[i][j] -= 1
					if grid[i][j] < 0 {
						grid[i][j] = 0
					}
				}
			}
		} else {
			for i := fromx; i <= tox; i++ {
				for j := fromy; j <= toy; j++ {
					grid[i][j] += 1
				}
			}
		}
	}
	acc := 0
	for i := 0; i < 1000; i++ {
		for j := 0; j < 1000; j++ {
			acc += grid[i][j]
		}
	}
	return acc
}
