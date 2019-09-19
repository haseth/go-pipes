package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

//Node make input file
type Node struct {
	stdin  *bytes.Buffer
	cmd    []string
	stdout *bytes.Buffer
}

//TakeInputOutputFile ...
func (n *Node) TakeInputOutputFile(ip chan *bytes.Buffer) {
	//take stdibytes.Buffer
	n.stdin = <-ip
	fmt.Println(" Taken the buffer")
}

//Process ...
func (n *Node) Process(op chan *bytes.Buffer, index int) {
	//process the command  by setting the input/output streaming file
	fmt.Println("Got command")
	execCmd := exec.Command(n.cmd[0], n.cmd[1:]...)
	//execCmd.Run()
	if index != 0 {
		execCmd.Stdin = n.stdin
	}
	fmt.Println("Stdin:" + n.stdin.String())
	//So we have the inpubytes.Buffer for input stream.
	// _, err := os.Open(n.stdin.Name())
	// if err != nil {
	// 	fmt.Println("Error in opening file" + err.Error())
	// }

	//make a stdoubytes.Buffer and attach as stdout
	// s := "temp" + strconv.Itoa(index+1)
	// fmt.Println("Testing name " + s)
	// fileOutput, err := os.Create(s)
	outputBuffer := new(bytes.Buffer)
	n.stdout = outputBuffer
	fmt.Println("Stdout:" + n.stdout.String())
	// execCmd.Stdout = os.Stdout
	// execCmd.Stderr = os.Stderr
	execCmd.Stdout = outputBuffer
	//execCmd.Stderr = fileOutput

	err := execCmd.Run()
	if err != nil {
		fmt.Println("Error in running file " + err.Error())
	} else {
		//fmt.Println(string(b))
	}
	fmt.Println("Runed the command.")

	op <- n.stdout
	// make a output file where we will put the output
	// n.stdout= testFile
	// cmd.Stdout("testFile")
	// run the command
}

//OutputFile ...
func (n *Node) OutputFile(op chan *bytes.Buffer) {
	op <- n.stdout
}

func main() {

	//commands := "cat test.txt | grep -i test"
	commandsArray := [][]string{
		[]string{"cat", "/Users/haseth/testing/go-pipes/pipes/test.txt"},
		[]string{"grep", "-i", "lopsem"},
		// []string{"ls", "-l"},
		// []string{"sort"},
		// []string{"less"},
	}
	lenOfCmdArray := len(commandsArray)
	//exec.Command("ls").String()

	//Array of channels of bytes.Buffer
	var channels []chan *bytes.Buffer
	channels = make([]chan *bytes.Buffer, lenOfCmdArray+1)
	//var f *os.File

	for index := range channels {
		channels[index] = make(chan *bytes.Buffer)
	}

	//first stdin file
	go func(channels []chan *bytes.Buffer) {
		b := new(bytes.Buffer)
		channels[0] <- b
	}(channels)

	//all good till here ..............
	//testing
	// index := 0
	// go commandStruct(commandsArray[0], channels[index], channels[index+1])
	// fmt.Println("Waiting for the output file to come")

	temp := 0
	for index, cmd := range commandsArray {
		//for each command new struct for each command with
		//input channels[i] as ip and channels[i+1] as outputs
		go commandStruct(cmd, channels[index], channels[index+1], index)
		temp = index + 1
	}
	//last channel should retrieve the output.
	f := <-channels[temp]
	fmt.Println(f.String())

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

func commandStruct(command []string, ip, op chan *bytes.Buffer, index int) {
	cmdStruct := &Node{cmd: command}
	cmdStruct.TakeInputOutputFile(ip)
	cmdStruct.Process(op, index)
}
