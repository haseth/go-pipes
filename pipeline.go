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

//Common stderr channel for pipeline
var stdErrChan = make(chan *bytes.Buffer)

//Commander interfaces the execute of each node
type Commander interface {
	Execute(stdin, stdout *bytes.Buffer) error
}

//Pipeline contains [][]commands and buffer for entire pipeline to store errors.
type Pipeline struct {
	commands []Commander
	stderr   *bytes.Buffer
}

//NewPipeline intializes a Pipeline struct
func NewPipeline(commands []Commander) *Pipeline {
	return &Pipeline{
		commands: commands,
		stderr:   new(bytes.Buffer),
	}
}

//Run starts the execution of the pipeline by creating
//n nodes and n+1 channels which connects the nodes.
func (p *Pipeline) Run() (string, error) {
	//defines number of nodes of pipe
	numOfCommands := len(p.commands)
	if numOfCommands == 0 {
		return "", errors.New(commandNil)
	}

	//make channel to link states
	channels, err := makeChannels(numOfCommands)
	if err != nil {
		return "", errors.New(err.Error())
	}

	//start the pipe by blank stdin
	go firstStdin(&channels[0])
	temp := 0
	for index, cmd := range p.commands {
		//new struct for each command with input channels[i] as input and channels[i+1] as outputs
		//pipeline stderr is comman buffer for storing errors.
		go cmdExecute(cmd, &channels[index], &channels[index+1], p.stderr)
		temp = index + 1
	}

	//There can be two outputs
	//1. from stdout
	//2. from stderr
	o := make([]byte, buffSize)
	e := make([]byte, buffSize)

	//retrieve stdout
	f := <-channels[temp]
	f.Read(o)

	//retrieve stderr
	p.stderr.Read(e)

	return string(o), errors.New(string(e))
}

//makeChannels initializes n+1 number of channels for n commands.
func makeChannels(n int) ([]chan *bytes.Buffer, error) {
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

//cmdExecute takes care for execution of each command.
func cmdExecute(command Commander, ip, op *chan *bytes.Buffer, stderr *bytes.Buffer) {
	//create node
	node, err := NewNode(command, stderr)
	if err != nil {
		commitError(err, stderr)
	}

	//input
	if err = node.Input(ip); err != nil {
		commitError(err, stderr)
	}

	//process
	if err = node.Process(); err != nil {
		commitError(err, stderr)
	}

	//output
	if err = node.Output(op); err != nil {
		commitError(err, stderr)
	}
}

//commit error of node
func commitError(err error, stderr *bytes.Buffer) {
	_, e := stderr.Write([]byte(err.Error()))
	if e != nil {
		return
	}
}

//firstStdin initializes first cmd's IP buffer.
func firstStdin(channel *chan *bytes.Buffer) {
	b := new(bytes.Buffer)
	*channel <- b
}
