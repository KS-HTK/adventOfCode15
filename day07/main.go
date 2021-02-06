package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

//wire
//	source:
//		constant
//		wire
//		gate
//	destinations:
//		list of wires and gates
// Wire: value or source (one of both has to be known)
type wire struct {
	value  uint16
	known  bool //true if value = signal
	source gate
}

//gates
//	inputs
//	type (oppcode)
//		0 := & AND
//		1 := | OR
//		2 := ^ XOR (NOT := 1111 ^ x)
//		3 := >> LSHIFT
//		4 := << RSHIFT
//	destinations
// Gate: 'inA', 'inB' and a oppcode 'opp'
type gate struct {
	inA string
	inB string
	opp string
}

//I know this is getting a little clutterd and overly complex (but learning...)
//stack for keeping a todo stack.
type stack []string

// IsEmpty: check if stack is empty
func (s *stack) IsEmpty() bool {
	return len(*s) == 0
}

// Pop: Remove and return top element of stack. Return false if empty
func (s *stack) Pop() (string, bool) {
	if s.IsEmpty() {
		return "", false
	} else {
		index := len(*s) - 1
		elem := (*s)[index]
		*s = (*s)[:index]
		return elem, true
	}
}

// Push: add a new value onto the stack
func (s *stack) Push(str string) {
	*s = append(*s, str)
}

// Peek: look at top value (not removing it)
func (s *stack) Peek() (string, bool) {
	if s.IsEmpty() {
		return "", false
	} else {
		return (*s)[len(*s)-1], true
	}
}

//Stacks did'nt work... Lets try Queues.
//Well, the Stack works fine, my Solution useing Stacks is the Problem.
type queue []string

// IsEmpty: returns true if queue is empty.
func (q *queue) IsEmpty() bool {
	return len(*q) == 0
}

//Push: add a new value to the end of the queue
func (q *queue) Push(str string) {
	*q = append(*q, str)
}

//Pop: remove and return first element of queue
func (q *queue) Pop() (string, bool) {
	if q.IsEmpty() {
		return "", false
	}
	elem := (*q)[0]
	*q = (*q)[1:]
	return elem, true
}

func errchk(e error) {
	if e != nil {
		panic(e)
	}
}

//check if a string is a constant
func isNumeric(s string) bool {
	_, err := strconv.ParseInt(s, 10, 16)
	return err == nil
}

//store all wires nicely on the map, referenced by their name
var diagram map[string]wire

func main() {
	diagram := make(map[string]wire)
	diagram["true"] = wire{value: 65535, known: true}
	dat, err := ioutil.ReadFile("input")
	errchk(err)
	//seperate lines
	commands := strings.Split(string(dat), "\n")
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
				val, _ := strconv.ParseInt(tmp[0], 10, 16)
				diagram[out] = wire{value: uint16(val), known: true}
			} else {
				//wire (maxINT AND x = x)
				diagram[out] = wire{known: false, source: gate{inA: "true", inB: tmp[0], opp: "AND"}}
			}
		} else if len(tmp) == 2 {
			//if input is NOT
			diagram[out] = wire{known: false, source: gate{inA: "true", inB: tmp[1], opp: "XOR"}}
		} else if len(tmp) == 3 {
			//if input is any other opperation
			diagram[out] = wire{known: false, source: gate{inA: tmp[0], inB: tmp[2], opp: tmp[1]}}

		}

	}
	fmt.Printf("Part 1: %d\n", pt1())
	fmt.Printf("Part 2: %d\n", pt2())
}

func pt1() uint16 {
	if diagram["a"].known {
		return diagram["a"].value
	}
	var todo stack
	todo.Push("a")
	for !todo.IsEmpty() {
		top, _ := todo.Peek()
		src := diagram[top].source
		if !isNumeric(src.inA) || !diagram[src.inA].known {
			todo.Push(src.inA)
			continue
		}
		if !isNumeric(src.inB) || !diagram[src.inB].known {
			todo.Push(src.inB)
			continue
		}
		top, _ = todo.Pop()
		assign(top, src)
	}
	fmt.Println(diagram)
	return diagram["a"].value
}

func assign(name string, source gate) {
	obj := diagram[name]
	var aVal, bVal uint16
	if isNumeric(source.inA) {
		tmp, err := strconv.ParseInt(source.inA, 10, 16)
		errchk(err)
		aVal = uint16(tmp)
	} else {
		aVal = diagram[source.inA].value
	}
	if isNumeric(source.inB) {
		tmp, err := strconv.ParseInt(source.inB, 10, 16)
		errchk(err)
		bVal = uint16(tmp)

	} else {
		bVal = diagram[source.inB].value
	}
	switch source.opp {
	case "AND":
		obj.value = aVal & bVal
	case "OR":
		obj.value = aVal | bVal
	case "XOR":
		obj.value = aVal ^ bVal
	case "LSHIFT":
		obj.value = aVal >> bVal
	case "RSHIFT":
		obj.value = aVal << bVal
	}
	obj.known = true
}

func pt2() int {
	return -1
}
