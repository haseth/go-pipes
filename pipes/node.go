package pipes

import (
	"bytes"
	"errors"
)

//Buffer ...
//Need help
//type Buffer bytes.Buffer

//Node make input file
const (
	errorInCommand = "Error running the command"
	stdErrBufNil   = "Std Error buf nil"
	stdInBufNil    = "StdIn buf nil"
	stdOutBufNil   = "StdOut buf nil"
)

//Processor ...
type Processor interface {
	Process()
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
	n.cmd = cmd
	n.stderr = stderr
	return n, nil
}

//Input takes ip buffer address
func (n *Node) Input(ip *chan *bytes.Buffer) error {
	n.stdin = <-*ip
	if n.stdin == nil {
		return errors.New(stdInBufNil)
	}
	return nil
}

//Process runs the commad attaching stdin, stdout and stderr.
func (n *Node) Process() error {
	outputBuffer := new(bytes.Buffer)
	n.stdout = outputBuffer
	if n.stdout == nil && n.stdin == nil && n.stderr == nil && n.cmd == nil {
		return errors.New("Buffers/State cannot be empty")
	}
	err := n.cmd.Execute(n.stdin, n.stdout)
	if err != nil {
		n.stderr.Write([]byte(err.Error()))
		//stdErrChan <- n.stderr
	}
	return nil
}

//Output produces the address of output buffer.
func (n *Node) Output(op *chan *bytes.Buffer) error {
	if n.stdout == nil {
		return errors.New(stdOutBufNil)
	}
	*op <- n.stdout
	return nil
}
