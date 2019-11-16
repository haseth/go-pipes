package pipes

import (
	"bytes"
	"errors"
)

const (
	commandNotFound  string = "not found"
	commandNil              = "Command cannot be empty"
	negativeChannels        = "Negative number of commands cannot be passed."
	buffSize                = 10000
)

var pipeErrChan = make(chan *bytes.Buffer)

/*
Pipeline pipes nodes together
*/
type Pipeline struct {
	nodes   []Node
	links   []chan *bytes.Buffer
	pipeErr *bytes.Buffer
}

// NewPipeline creates nodes, its links and std err
func NewPipeline(executers []Executer) *Pipeline {
	nodeLinks, _ := makeLinks(len(executers))

	nodes := make([]Node, 0)
	for _, exec := range executers {
		node, err := NewNodeState(exec)
		if err != nil {
			return nil
		}
		nodes = append(nodes, node)
	}

	return &Pipeline{
		nodes:   nodes,
		links:   nodeLinks,
		pipeErr: new(bytes.Buffer),
	}
}

/*
Run executes the pipe nodes with IP and OP links
*/
func (p *Pipeline) Run() (string, error) {

	// start the pipe by blank stdin
	go firstStdin(p.links[0])

	lastCh := 0
	for index, node := range p.nodes {
		go nodeExecute(node, p.links[index], p.links[index+1], p.pipeErr)
		lastCh = index + 1
	}

	//OUTPUT

	//There can be two outputs
	//1. from stdout
	//2. from stderr
	o := make([]byte, buffSize)
	e := make([]byte, buffSize)

	//retrieve stdout
	f := <-p.links[lastCh]
	f.Read(o)

	//retrieve stderr
	p.pipeErr.Read(e)

	return string(o), errors.New(string(e))
}

//makeLinks initializes n+1 number of channels for n executers.
func makeLinks(n int) ([]chan *bytes.Buffer, error) {
	//error checking
	if n == 0 {
		return nil, errors.New(commandNil)
	}
	if n < 0 {
		return nil, errors.New(negativeChannels)
	}

	channels := make([]chan *bytes.Buffer, n+1)
	for index := range channels {
		channels[index] = make(chan *bytes.Buffer)
	}
	return channels, nil
}

//nodeExecute takes care for execution of each node.
func nodeExecute(node Node, ip, op chan *bytes.Buffer, stderr *bytes.Buffer) {
	//input
	if err := node.Input(ip); err != nil {
		commitError(err, stderr)
	}

	//process
	if err := node.Process(); err != nil {
		commitError(err, stderr)
	}

	//output
	if err := node.Output(op); err != nil {
		commitError(err, stderr)
	}
}

//commit error logs to pipe error buffer
func commitError(err error, stderr *bytes.Buffer) {
	_, e := stderr.Write([]byte(err.Error()))
	if e != nil {
		return
	}
}

//firstStdin initializes first cmd's IP buffer.
func firstStdin(channel chan *bytes.Buffer) {
	b := new(bytes.Buffer)
	channel <- b
}
