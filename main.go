package main

import (
	"fmt"
	"help/go-pipes/pipes"
)

func main() {
	//commands := "cat test.txt | grep -i test"
	//TODO: change to something which can take from any sort of I/P webpage, from file, or mock...
	commandsArray := [][]string{
		//[]string{},
		[]string{"cat", "/Users/haseth/testing/go-pipes/pipes/test.txt"},
		[]string{"grep", "-i", "harsh"},
		[]string{"wc", "-clw"},
		[]string{"ls", "-ltrh"},
		// []string{"sort"},
		// []string{"less"},
	}

	//Create a pipeline for executing commands
	pipe := pipes.NewPipeline(commandsArray)

	//Output of the commands.
	out, err := pipe.Run()
	if err != nil {
		fmt.Println("Some issue" + err.Error())
	}
	fmt.Println(out)
}
