package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

//Node make input file
type Node struct {
	stdin  *os.File
	cmd    []string
	stdout *os.File
}

//TakeInputOutputFile ...
func (n *Node) TakeInputOutputFile(ip chan *os.File) {
	//take stdin file
	n.stdin = <-ip
	fmt.Println(" Taken the file")
}

//Process ...
func (n *Node) Process(op chan *os.File, index int) {
	//process the command  by setting the input/output streaming file
	execCmd := exec.Command(n.cmd[0], n.cmd[1:]...)
	execCmd.Stdin = n.stdin
	_, err := os.Open(n.stdin.Name())
	if err != nil {
		fmt.Println("Error in opening file" + err.Error())
	}
	//make a stdout file
	s := "temp" + strconv.Itoa(index+1)
	fmt.Println("Testing name " + s)
	fileOutput, err := os.Create(s)

	if err != nil {
		fmt.Println("Error in creating file" + err.Error())

	}

	execCmd.Stdout = fileOutput

	err = execCmd.Run()
	if err != nil {
		fmt.Println("Error in reading file" + err.Error())
	}
	n.stdout = fileOutput
	op <- n.stdout
	// make a output file where we will put the output
	// n.stdout= testFile
	// cmd.Stdout("testFile")
	// run the command
}

//OutputFile ...
func (n *Node) OutputFile(op chan *os.File) {
	op <- n.stdout
}

func main() {

	//commands := "cat test.txt | grep -i test"
	commandsArray := [][]string{
		//[]string{"echo", "hello"},
		[]string{"cat", "/Users/haseth/testing/go-pipes/pipes/test.txt"},
		[]string{"grep", "-i", "harsh"},
		[]string{"ls"},
		[]string{"sort"},
		[]string{"less"},
	}
	lenOfCmdArray := len(commandsArray)
	//exec.Command("ls").String()

	//Array of channels of os.File
	var channels []chan *os.File
	channels = make([]chan *os.File, lenOfCmdArray+1)
	//var f *os.File

	for index := range channels {
		channels[index] = make(chan *os.File)
	}

	//first stdin file
	go func(channels []chan *os.File) {
		f, err := os.Create("temp123")

		if err != nil {
			fmt.Println("Error in creating file" + err.Error())
		}
		channels[0] <- f
	}(channels)

	//testing
	// index := 0
	// go commandStruct(commandsArray[0], channels[index], channels[index+1])
	// fmt.Println("Waiting for the output file to come")

	temp := 0
	for index, cmd := range commandsArray {
		//for each command
		//new struct for each command with
		//input channels[i] as ip and channels[i+1] as outputs
		go commandStruct(cmd, channels[index], channels[index+1], index)
		temp = index + 1
	}

	f := <-channels[temp]
	fmt.Println(f.Name())

	// var b []byte
	// b = make([]byte, 10)
	// //fmt.Println(f.Rea())
	// n, err := f.Read(b)
	// if err != nil {
	// 	fmt.Println("You got error in reading file " + err.Error())
	// }
	// fmt.Println("No. of chars", n)
	// fmt.Println(string(b))

}

func commandStruct(command []string, ip, op chan *os.File, index int) {
	cmdStruct := &Node{cmd: command}
	cmdStruct.TakeInputOutputFile(ip)
	cmdStruct.Process(op, index)

}
