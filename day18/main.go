package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
)

func errchk(e error) {
	if e != nil {
		panic(e)
	}
}

func interpret(lines []string) [100][100]bool {
	var grid [100][100]bool
	for y, line := range lines {
		if line == "" {
			continue
		}
		x := 0
		for _, c := range line {
			if c == '#' {
				grid[y][x] = true
			}
			x += 1
		}
	}
	return grid
}

func countOnNeighbors(in [100][100]bool, x int, y int) int {
	var n [][]bool
	var xMin, xMax int
	if x > 0 {
		xMin = x - 1
	} else {
		xMin = x
	}
	if x < 99 {
		xMax = x + 2
	} else {
		xMax = x + 1
	}
	if y > 0 {
		n = append(n, in[y-1][xMin:xMax])
	}
	n = append(n, in[y][xMin:xMax])
	if y < 99 {
		n = append(n, in[y+1][xMin:xMax])
	}
	count := 0
	for _, l := range n {
		for _, b := range l {
			if b {
				count++
			}
		}
	}
	if in[y][x] {
		return count - 1
	}
	return count
}

func countOn(in [100][100]bool) int {
	count := 0
	for _, line := range in {
		for _, light := range line {
			if light {
				count++
			}
		}
	}
	return count
}

func next(in [100][100]bool, two bool) [100][100]bool {
	var out [100][100]bool
	for y, line := range in {
		for x, state := range line {
			count := countOnNeighbors(in, x, y)
			if state {
				if count == 2 || count == 3 {
					out[y][x] = true
				}
			} else {
				if count == 3 {
					out[y][x] = true
				}
			}
		}
	}
	if two {
		out[0][0] = true
		out[0][99] = true
		out[99][0] = true
		out[99][99] = true
	}
	return out
}

func run(steps int, in [100][100]bool, two bool) [100][100]bool {
	if steps == 1 {
		return next(in, two)
	}
	return run(steps-1, next(in, two), two)
}

func main() {
	dat, err := ioutil.ReadFile("input")
	errchk(err)
	lines := strings.Split(string(dat), "\n")
	var grid, grid2 [100][100]bool
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		grid = interpret(lines)
		wg.Done()
	}()
	go func() {
		grid2 = interpret(lines)
		grid2[0][0] = true
		grid2[0][99] = true
		grid2[99][0] = true
		grid2[99][99] = true
		wg.Done()
	}()
	wg.Wait()

	grid = run(100, grid, false)
	count := countOn(grid)
	grid2 = run(100, grid2, true)
	count2 := countOn(grid2)

	fmt.Printf("Part 1: %d\n", count)
	fmt.Printf("Part 2: %d\n", count2)

}
