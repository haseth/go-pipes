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

	//last channel should retrieve the output.
	for {
		select {
		case f := <-channels[temp]:
			return f.String(), nil

		case e := <-stdErrChan:
			return "", errors.New(e.String())
		}
	}
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
	if command == nil {
		e := commandNil
		stderr.Write([]byte(e))
		stdErrChan <- stderr
	}
	node, err := NewNode(command, stderr)
	if err != nil {
		stderr.Write([]byte(err.Error()))
		stdErrChan <- stderr
	}
	node.Input(ip)
	node.Process()
	node.Output(op)
}

//firstStdin initializes first cmd's IP buffer.
func firstStdin(channel *chan *bytes.Buffer) {
	b := new(bytes.Buffer)
	*channel <- b
}
