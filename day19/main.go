package main

import (
	"fmt"
	"io/ioutil"
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
func getResultCount(in string, t map[string][]string) int {
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

//job strct stores a molecule and the amount of steps used to synthesize it
type job struct {
	mol  string
	step int
}

//store all molecules and minimal steps to get to them in the molindex
//if a mol has already been indexed do not add it to the todo channel
//if the med molecule has been reached close all further jobs
func indexer(jobs <-chan job, todo chan<- job, med string) int {
	//store all found molecules and the min steps to synthesize it
	molindex := make(map[string]int)
	open := true
	for j := range jobs {
		if j.mol == med {
			if open {
				open = false
				close(todo)
			}
		}
		if len(j.mol) > len(med) {
			continue
		}
		if min, ok := molindex[j.mol]; ok {
			if min > j.step {
				molindex[j.mol] = j.step
			}
		} else {
			molindex[j.mol] = j.step
			if open {
				todo <- j
			}
		}
	}
	return molindex[med]
}

//worker for generating the future jobs based on a incomming job
func worker(jobs <-chan job, addJob chan<- job, t map[string][]string) {
	for j := range jobs {
		lst := getAllTransformations(j.mol, t)
		for key := range lst {
			addJob <- job{mol: key, step: j.step + 1}
		}
	}
}

//function for initalizing indexer and worker and overseeing their work
func makeMed(med string, t map[string][]string, res chan int) {
	jobsForIndexer := make(chan job, 1000000)
	jobsForWorkers := make(chan job, 1000000)
	go func() {
		res <- indexer(jobsForIndexer, jobsForWorkers, med)
	}()
	jobsForIndexer <- job{mol: "e", step: 0}
	var wg sync.WaitGroup
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func() {
			worker(jobsForWorkers, jobsForIndexer, t)
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
	strChannel := make(chan string, 1)
	transitions := interpret(lines, strChannel)
	med = <-strChannel
	res1, res2 := make(chan int), make(chan int)
	go func() {
		res1 <- getResultCount(med, transitions)
		makeMed(med, transitions, res2)
	}()
	fmt.Printf("Part 1: %d\n", <-res1)
	fmt.Printf("Part 2: %d\n", <-res2)
}
