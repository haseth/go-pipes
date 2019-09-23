package main

import (
	"fmt"
	"help/go-pipes/pipes"
)

// type Commander interface{
// 	Execute() ([]byte,error)
// }

// type CatCommand struct{}

// func (c *CatCommand) Execute() ([]byte,error){
// 	res,err:=os.exec("cat","")
// }

// type GrepCommand struct{}

// func (c *GrepCommand) Execute() ([]byte,error){
// 	res,err:=os.exec("grep","")
// }

// type MockCatCommand struct{}

// func (c *MockCatCommand) Execute() ([]byte,error){
// 	return []byte{"hello"},nil
// }

func main() {
	//commands := "cat test.txt | grep -i test"
	//TODO: change to something which can take from any sort of I/P webpage, from file, or mock...
	commandsArray := [][]string{
		[]string{"cat", "/Users/haseth/testing/go-pipes/pipes/test.txt"},
		[]string{"grep", "-i", "harsh"},
		[]string{"wc", "-clw"},
		[]string{"ls", "-ltrh"},
		// []string{"sort"},
		// []string{"less"},
	}

	//Create a pipeline
	pipe := pipes.NewPipeline(commandsArray)

	//Output
	out, err := pipe.Run()
	if err != nil {
		fmt.Println("Some issue" + err.Error())
	}
	fmt.Println(out)
}
