package pipes

import (
	"bytes"
	"fmt"
	"os/exec"
)

//Buffer ...
type Buffer *bytes.Buffer

//Node make input file
type Node struct {
	stdin  Buffer
	cmd    []string
	stdout Buffer
}

//NewNode initialises the values.
func NewNode(cmd []string, stdin, stdout Buffer) *Node {
	n := new(Node)
	n.stdin = stdin
	n.cmd = cmd
	n.stdout = stdout
	return n
}

//SetCommand ...
func (n *Node) SetCommand(command []string) {
	//TODO: check if commands are all good.
	n.cmd = command
	fmt.Println(command)
}

//Input ...
func (n *Node) Input(ip chan Buffer) {
	//take stdibytes.Buffer
	//TODO: check if buffer nil
	n.stdin = <-ip
	fmt.Println(" Taken the buffer")
}

//Process ...
func (n *Node) Process() {
	//process the command  by setting the input/output streaming file
	//TODO: cheking for error in stderr
	outputBuffer := new(Buffer)
	execCmd := exec.Command(n.cmd[0], n.cmd[1:]...)

	execCmd.Stdin = n.stdin

	execCmd.Stdout = outputBuffer

	err := execCmd.Run()
	if err != nil {
		fmt.Println("Error in running file " + err.Error())
	}
	n.stdout = outputBuffer
}

//Output ...
func (n *Node) Output(op chan Buffer) {
	//TODO: Check if nil
	op <- n.stdout
}
