package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

//commands
//	inputs
//	opperation
//		0 := & AND (constants are stored as true AND value)
//		1 := | OR
//		2 := ^ XOR (NOT is stored as true XOR value)
//		3 := >> LSHIFT
//		4 := << RSHIFT
// Command: 'input A', 'input B' and a 'opperation'
type cmd struct {
	inA  string
	inB  string
	opp  string
	name string
}

//Stacks did'nt work... Lets try Queues.
//Well, the Stack works fine, my Solution useing Stacks is the Problem.
type queue []cmd

// IsEmpty: returns true if queue is empty.
func (q *queue) IsEmpty() bool {
	return len(*q) == 0
}

//Push: add a new value to the end of the queue
func (q *queue) Push(com cmd) {
	*q = append(*q, com)
}

//Pop: remove and return first element of queue
func (q *queue) Pop() (cmd, bool) {
	if q.IsEmpty() {
		return cmd{}, false
	}
	elem := (*q)[0]
	*q = (*q)[1:]
	return elem, true
}

//Peek: return first element of queue
func (q *queue) Peek() (cmd, bool) {
	if q.IsEmpty() {
		return cmd{}, false
	}
	return (*q)[0], true
}

func errchk(e error) {
	if e != nil {
		panic(e)
	}
}

//check if a string is a number
func isNumeric(s string) bool {
	_, err := strconv.ParseInt(s, 0, 32)
	return err == nil
}

//store all wire-values nicely on the map, referenced by their name
var diagram map[string]uint16

//queue of all unsolved commands
var todo queue

func main() {
	//read input file
	dat, err := ioutil.ReadFile("input")
	errchk(err)
	initial(string(dat))
	solve()
	res := getA()
	fmt.Printf("Part 1: %d\n", res)
	initial(string(dat))
	diagram["b"] = res
	solve()
	fmt.Printf("Part 2: %d\n", getA())
}

func initial(dat string) {
	diagram = make(map[string]uint16)
	diagram["true"] = 0b1111111111111111

	//seperate lines
	commands := strings.Split(dat, "\n")
	for _, s := range commands {
		if s == "" {
			continue
		}
		//split line into input and output
		tmp := strings.Split(s, " -> ")
		out := tmp[1]
		//split input into parts
		tmp = strings.Split(tmp[0], " ")
		//if input has only one part (wire or constant)
		if len(tmp) == 1 {
			if isNumeric(tmp[0]) {
				//input is constant
				val, _ := strconv.ParseInt(tmp[0], 10, 32)
				diagram[out] = uint16(val)
			} else {
				//input is wire
				//wire (true AND val = val)
				todo.Push(cmd{inA: "true", inB: tmp[0], opp: "AND", name: out})
			}
		} else if len(tmp) == 2 {
			//input is NOT
			//true XOR val = NOT val
			todo.Push(cmd{inA: "true", inB: tmp[1], opp: "XOR", name: out})
		} else if len(tmp) == 3 {
			//if input is any other opperation
			todo.Push(cmd{inA: tmp[0], inB: tmp[2], opp: tmp[1], name: out})

		}
	}
}

func solve() {
	for !todo.IsEmpty() {
		//queue should not be empty.
		c, _ := todo.Pop()
		if !assign(c) {
			todo.Push(c)
		}
	}
}

func assign(src cmd) bool {
	var aVal, bVal uint16
	if isNumeric(src.inA) {
		tmp, err := strconv.ParseInt(src.inA, 10, 32)
		errchk(err)
		aVal = uint16(tmp)
	} else if val, ok := diagram[src.inA]; ok {
		aVal = val
	} else {
		return false
	}
	if isNumeric(src.inB) {
		tmp, err := strconv.ParseInt(src.inB, 10, 32)
		errchk(err)
		bVal = uint16(tmp)
	} else if val, ok := diagram[src.inB]; ok {
		bVal = val
	} else {
		return false
	}
	switch src.opp {
	case "AND":
		diagram[src.name] = aVal & bVal
	case "OR":
		diagram[src.name] = aVal | bVal
	case "XOR":
		diagram[src.name] = aVal ^ bVal
	case "LSHIFT":
		diagram[src.name] = aVal << bVal
	case "RSHIFT":
		diagram[src.name] = aVal >> bVal
	}
	return true
}

func getA() uint16 {
	if val, ok := diagram["a"]; ok {
		return val
	}
	return 0
}
