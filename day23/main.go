package main

import (
	"errors"
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

var a, b uint
var ins_count int

func hlf(r *uint) {
	*r /= 2
}
func tpl(r *uint) {
	*r *= 3
}
func inc(r *uint) {
	*r += 1
}
func jmp(offset int) {
	ins_count += offset
}
func jie(r *uint, offset int) {
	if *r%2 == 0 {
		jmp(offset)
	}
}
func jio(r *uint, offset int) {
	if *r == 1 {
		jmp(offset)
	}
}

func interpret(ins []string) {
	ins_count = 0
	for ins_count < len(ins) {
		if ins[ins_count] == "" {
			return
		}
		inst := ins[ins_count][:3]
		args := strings.Split(ins[ins_count][4:], ", ")
		var offset int
		var reg *uint
		var err error
		if inst == "jmp" {
			offset, err = strconv.Atoi(args[0])
			errchk(err)
		} else {
			if args[0] == "a" {
				reg = &a
			} else {
				reg = &b
			}
		}
		if inst == "jie" || inst == "jio" {
			offset, err = strconv.Atoi(args[1])
			errchk(err)
		}
		offset--
		switch inst {
		case "hlf":
			hlf(reg)
		case "tpl":
			tpl(reg)
		case "inc":
			inc(reg)
		case "jmp":
			jmp(offset)
		case "jie":
			jie(reg, offset)
		case "jio":
			jio(reg, offset)
		default:
			errchk(errors.New("Unknown Instruction " + inst))
		}
		ins_count++
	}
}

func main() {
	dat, err := ioutil.ReadFile("input")
	errchk(err)
	lines := strings.Split(string(dat), "\n")

	//channel for returning results
	res1, res2 := make(chan uint), make(chan uint)

	go func() {
		a = 0
		b = 0
		interpret(lines)
		res1 <- b
		a = 1
		b = 0
		interpret(lines)
		res2 <- b
	}()
	fmt.Printf("Part 1: %d\n", <-res1)
	fmt.Printf("Part 2: %d\n", <-res2)
}
