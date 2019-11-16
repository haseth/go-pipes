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

//Executer interfaces the execute of each node
type Executer interface {
	Execute(stdin, stdout *bytes.Buffer) error
}

//Node extracts Node functionality
type Node interface {
	Input(chan *bytes.Buffer) error
	Process() error
	Output(op chan *bytes.Buffer) error
}

//NodeState represents single command with IP/OP buffers.
type NodeState struct {
	stdin    *bytes.Buffer
	executer Executer
	stdout   *bytes.Buffer
}

//NewNodeState initialises the values.
func NewNodeState(executer Executer) (Node, error) {
	n := new(NodeState)
	if n == nil {
		return nil, errors.New("Error in creating NodeState")
	}

	n.stdin = new(bytes.Buffer)
	n.executer = executer
	n.stdout = new(bytes.Buffer)

	return n, nil
}

//Input takes ip buffer address
func (n *NodeState) Input(ip chan *bytes.Buffer) error {
	n.stdin = <-ip
	// if received nil buffer it might raise error
	if n.stdin == nil {
		//as next state will use this as stdin.
		ipBuffer := new(bytes.Buffer)
		n.stdin = ipBuffer
		return errors.New(stdInBufNil)
	}

	return nil
}

//Process runs the commad attaching stdin, stdout and stderr.
func (n *NodeState) Process() error {
	//error checking
	if n.stdout == nil && n.stdin == nil && n.executer == nil {
		return errors.New("Buffers/State cannot be empty")
	}

	//execute
	err := n.executer.Execute(n.stdin, n.stdout)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

//Output produces the address of output buffer.
func (n *NodeState) Output(op chan *bytes.Buffer) error {
	if n.stdout == nil {
		outputBuffer := new(bytes.Buffer)
		n.stdout = outputBuffer
		op <- n.stdout
		return errors.New(stdOutBufNil)
	}
	op <- n.stdout

	return nil
}
