package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"sync"
)

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
				if arr, ok := trans[in[0]]; ok {
					trans[in[0]] = append(arr, in[1])
				} else {
					trans[in[0]] = append(make([]string, 0), in[1])
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
func countTransformations(in string, t map[string][]string) int {
	return len(getAllTransformations(in, t))
}

//get all single step transformations possible on in string useing t(ransition) map
//result contains no duplicats as it is the keyset of the returned map
//value set of returned map is only 0 (no Information is stored here)
func getAllTransformations(in string, t map[string][]string) map[string]int {
	resultMap := make(map[string]int)
	previousInd := 0
	previousChar := ""
	for i, c := range in {
		if slice, ok := t[string(c)]; ok {
			for _, rStr := range slice {
				resStr := in[:i] + rStr + in[i+1:]
				resultMap[resStr] = 0
			}
		}
		c2 := previousChar + string(c)
		if slice, ok := t[c2]; ok {
			for _, rStr := range slice {
				resStr := in[:previousInd] + rStr + in[i+1:]
				resultMap[resStr] = 0
			}
		}
		previousChar = string(c)
		previousInd = i
	}
	return resultMap
}

//Helper Function: get index of all occurences of substring in string
func indexAll(str string, substr string) []int {
	ind := make(map[int]int, 0)
	center := len(str) / 2
	if center < len(substr) {
		for i := 0; i < len(str); i++ {
			index := strings.Index(str[i:], substr)
			if index != -1 {
				ind[index+i] = 0
			}
		}
	} else {
		for i := 0; i <= center; i++ {
			index := strings.Index(str[i:center], substr)
			lIndex := strings.LastIndex(str[center:len(str)-i], substr)
			if index != -1 {
				ind[index+i] = 0
			}
			if lIndex != -1 {
				ind[lIndex+center] = 0
			}
		}
	}
	out := make([]int, 0)
	for k := range ind {
		out = append(out, k)
	}
	sort.Ints(out)
	return out
}

//Helper Function: reverse the Mappings and sort them by Length
func reverseTransitions(in map[string][]string) map[string]string {
	out := make(map[string]string)
	for key, val := range in {
		for _, v := range val {
			if exists, ok := out[v]; ok {
				fmt.Println("Value Overwritten: ", v, exists, key)
			}
			out[v] = key
		}
	}
	return out
}

//job strct stores a molecule and the amount of steps used to synthesize it
type job struct {
	mol  string
	prev string
	step int
}

//store all molecules and minimal steps to get to them in the molindex
//if a mol has already been indexed do not add it to the todo channel
//if the med molecule has been reached close all further jobs
func indexer(jobs <-chan job, todo chan<- job, med string, res chan<- int) {
	//store all found molecules and the min steps to synthesize it
	molindex := make(map[string]job)
	open := true
	for j := range jobs {
		if j.mol == "e" {
			if open {
				open = false
				close(todo)
			}
		} else if strings.Contains(j.mol, "e") {
			//if string contains any "e" discard any further job
			continue
		}
		if mol, ok := molindex[j.mol]; ok {
			if mol.step > j.step {
				molindex[j.mol] = j
			}
		} else {
			molindex[j.mol] = j
			if open {
				todo <- j
			}
		}
	}
	res <- getStepCount(molindex)
}

//worker for generating the future jobs based on a incomming job
func worker(jobs <-chan job, addJob chan<- job, rt map[string]string) {
	for j := range jobs {
		for key, val := range rt {
			ind := indexAll(j.mol, key)
			for _, i := range ind {
				molecule := j.mol[:i] + val + j.mol[i+len(key):]
				addJob <- job{mol: molecule, prev: j.mol, step: j.step + 1}
			}
		}
	}
}

//function for tracing the steps trough the map
func getStepCount(m map[string]job) int {
	count := 0
	j := m["e"]
	for j.prev != "" {
		j = m[j.prev]
		count++
	}
	return count
}

//function for initalizing indexer and worker and overseeing their work
func makeMed(med string, rt map[string]string, res chan<- int) {
	jobsForIndexer := make(chan job, 10000000)
	jobsForWorkers := make(chan job, 10000000)
	go indexer(jobsForIndexer, jobsForWorkers, med, res)
	jobsForIndexer <- job{mol: med, step: 0, prev: ""}
	var wg sync.WaitGroup
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func() {
			worker(jobsForWorkers, jobsForIndexer, rt)
			wg.Done()
		}()
	}
	wg.Wait()
	close(jobsForIndexer)
}

//main function.
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
		makeMed(med, reverseTransitions(transitions), res2)
	}()
	fmt.Printf("Part 1: %d\n", <-res1)
	fmt.Printf("Part 2: %d\n", <-res2)
}
