package pipes

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

//Buffer ...
//Need help
//type Buffer bytes.Buffer

//Node make input file

const errorInCommand = "Error running the command"
const stdErrBufNil = "Std Error buf nil"

//Node ...
type Node struct {
	stdin  *bytes.Buffer
	cmd    []string
	stdout *bytes.Buffer
	//common for entire pipeline
	stderr *bytes.Buffer
}

//NewNode initialises the values.
func NewNode(cmd []string, stderr *bytes.Buffer) (*Node, error) {

	n := new(Node)
	if len(cmd) == 0 {
		return nil, errors.New(commandNil)
	}
	if stderr == nil {
		return nil, errors.New(stdErrBufNil)
	}
	n.cmd = cmd
	n.stderr = stderr

	return n, nil
}

//SetCommand ...
func (n *Node) SetCommand(command []string) error {
	//TODO: check if commands are all good.
	if len(command) == 0 {
		return errors.New(commandNil)
	}
	n.cmd = command
	return nil
}

//Input ...
func (n *Node) Input(ip *chan *bytes.Buffer) error {
	//take stdibytes.Buffer
	//TODO: check if buffer nil
	n.stdin = <-*ip
	if n.stdin == nil {
		return errors.New(stdErrBufNil)
	}
	return nil
}

//Process ...
func (n *Node) Process() {
	//process the command  by setting the input/output streaming file
	//TODO: cheking for error in stderr

	outputBuffer := new(bytes.Buffer)
	execCmd := exec.Command(n.cmd[0], n.cmd[1:]...)

	execCmd.Stdin = n.stdin
	execCmd.Stdout = outputBuffer

	fmt.Println(outputBuffer)
	err := execCmd.Run()
	if err != nil {
		fmt.Println("Error in running file " + err.Error())
		e := "Error in running file " + strings.Join(n.cmd, " ") + err.Error()
		n.stderr.Write([]byte(e))
		stdErrChan <- n.stderr
	}
	n.stdout = outputBuffer
}

//Output ...
func (n *Node) Output(op *chan *bytes.Buffer) error {
	//TODO: Check if nil
	*op <- n.stdout
	if *op == nil {
		return errors.New(stdErrBufNil)
	}
	return nil
}
