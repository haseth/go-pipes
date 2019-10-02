package pipes

import (
	"bytes"
	"errors"
)

const commandNotFound string = "not found"
const commandNil = "Command cannot be empty"
const negativeChannels = "Negative number of commands cannot be passed."

var stdErrChan = make(chan *bytes.Buffer)

//Commander ...
type Commander interface {
	Execute(stdin, stdout *bytes.Buffer) error
}

//Pipeline contains [][]commands and buffer for entire pipeline to store errors.
type Pipeline struct {
	commands []Commander
	stderr   *bytes.Buffer
}

//Run starts the execution of the pipeline by creating
//n nodes and n+1 channels which connects the nodes.
func (p *Pipeline) Run() (string, error) {
	numOfCommands := len(p.commands)
	if numOfCommands == 0 {
		return "", errors.New(commandNil)
	}

	channels, err := makeChannels(numOfCommands)
	if err != nil {
		return "", errors.New(err.Error())
	}

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

	o := make([]byte, 100)
	e := make([]byte, 1000)

	//last channel should retrieve the output.
	f := <-channels[temp]
	f.Read(o)
	p.stderr.Read(e)
	return string(o), errors.New(string(e))
}

//NewPipeline intializes a Pipeline struct.
func NewPipeline(commands []Commander) *Pipeline {
	return &Pipeline{
		commands: commands,
		stderr:   new(bytes.Buffer),
	}
}

//makeChannels initializes n+1 number of channels for n commands.
func makeChannels(n int) ([]chan *bytes.Buffer, error) {
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
	// if command == nil {
	// 	e := errors.New(commandNil)
	// 	commitError(e, stderr)
	// 	return
	// }
	node, err := NewNode(command, stderr)
	if err != nil {
		commitError(err, stderr)
	}
	if err = node.Input(ip); err != nil {

		commitError(err, stderr)
	}
	if err = node.Process(); err != nil {
		commitError(err, stderr)
	}
	if err = node.Output(op); err != nil {
		commitError(err, stderr)
	}
}
func commitError(err error, stderr *bytes.Buffer) {
	//fmt.Println(err.Error())
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
