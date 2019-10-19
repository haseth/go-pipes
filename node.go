package pipes

import (
	"bytes"
	"errors"
)

const (
	errorInCommand = "Error running the command"
	stdErrBufNil   = "Std Error buf nil"
	stdInBufNil    = "StdIn buf nil"
	stdOutBufNil   = "StdOut buf nil"
)

//NodeState extracts node functionality
type NodeState interface {
	Input(*chan *bytes.Buffer) error
	Process()
	Output(op *chan *bytes.Buffer) error
}

//Node represents single command with IP/OP buffers.
type Node struct {
	stdin  *bytes.Buffer
	cmd    Commander
	stdout *bytes.Buffer
	stderr *bytes.Buffer //common for entire pipeline
}

//NewNode initialises the values.
func NewNode(cmd Commander, stderr *bytes.Buffer) (*Node, error) {
	//what if error is therer.
	n := new(Node)
	if n == nil {
		return nil, errors.New("Error in creating Node")
	}

	//define node values
	n.cmd = cmd
	n.stderr = stderr
	n.stdout = new(bytes.Buffer)

	return n, nil
}

//Input takes ip buffer address
func (n *Node) Input(ip *chan *bytes.Buffer) error {
	n.stdin = <-*ip
	if n.stdin == nil {
		//as next state will use this as stdin.
		ipBuffer := new(bytes.Buffer)
		n.stdin = ipBuffer
		return errors.New(stdInBufNil)
	}

	return nil
}

//Process runs the commad attaching stdin, stdout and stderr.
func (n *Node) Process() error {
	//error checking
	if n.stdout == nil && n.stdin == nil && n.stderr == nil && n.cmd == nil {
		return errors.New("Buffers/State cannot be empty")
	}

	//execute
	err := n.cmd.Execute(n.stdin, n.stdout)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

//Output produces the address of output buffer.
func (n *Node) Output(op *chan *bytes.Buffer) error {
	if n.stdout == nil {
		outputBuffer := new(bytes.Buffer)
		n.stdout = outputBuffer
		*op <- n.stdout
		return errors.New(stdOutBufNil)
	}
	*op <- n.stdout

	return nil
}
