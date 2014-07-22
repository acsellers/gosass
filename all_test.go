package sassy

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"testing"
)

func runParallel(testFunc func(chan bool), concurrency int) {
	runtime.GOMAXPROCS(4)
	done := make(chan bool, concurrency)
	for i := 0; i < concurrency; i++ {
		go testFunc(done)
	}
	for i := 0; i < concurrency; i++ {
		<-done
		<-done
	}
	runtime.GOMAXPROCS(1)
}

const numConcurrentRuns = 200

// const testFileName1     = "test1.scss"
// const testFileName2     = "test2.scss"
// const desiredOutput     = "div {\n  color: black; }\n  div span {\n    color: blue; }\n"

func compileTest(t *testing.T, fileName string) (result string) {

	ctx := FileSet{
		Style:      NestedStyle,
		IncludeDir: []string{},
	}

	f, err := ctx.ParseFile(fileName)

	if err != nil {
		t.Error("ERROR: ", err)
	} else {
		result = f.Output
	}

	return result
}

const numTests = 4 // TO DO: read the test dir and set this dynamically

func TestConcurrent(t *testing.T) {
	testFunc := func(done chan bool) {
		done <- false
		for i := 1; i <= numTests; i++ {
			inputFile := fmt.Sprintf("test/test%d.scss", i)
			result := compileTest(t, inputFile)
			desiredOutput, err := ioutil.ReadFile(fmt.Sprintf("test/test%d.css", i))
			if err != nil {
				t.Error(fmt.Sprintf("ERROR: couldn't read test/test%d.css", i))
			}
			if result != string(desiredOutput) {
				t.Error(result, string(desiredOutput))
				t.Error("ERROR: incorrect output")
			}
		}
		done <- true
	}
	runParallel(testFunc, numConcurrentRuns)
}
