package pipes

import (
	"bytes"
	"errors"
	"fmt"
)

const commandNotFound string = "Command not found"
const commandNil = "Command cannot be empty"
const negativeChannels = "Negative number of commands cannot be passed."

var stdErrChan = make(chan *bytes.Buffer)

// type Commander interface{
// 	Execute() ([]byte,error)
// }

// type CatCommand struct{}

// func (c *CatCommand) Execute() ([]byte,error){
// 	res,err:=os.exec("cat","")
// }

// type GrepCommand struct{}

// func (c *GrepCommand) Execute() ([]byte,error){
// 	res,err:=os.exec("grep","")
// }

//Pipeline contains [][]commands and buffer for entire pipeline to store errors.
type Pipeline struct {
	commands [][]string
	stderr   *bytes.Buffer
}

//Run starts the execution of the pipeline by creating
//n nodes and n+1 channels which connects the nodes.
func (p *Pipeline) Run() (string, error) {
	numOfCommands := len(p.commands)

	channels, err := makeChannels(numOfCommands)
	if err != nil {
		fmt.Println("Somme issue in creating channels")
		return "", errors.New("Error in creating channel")
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
func NewPipeline(commands [][]string) *Pipeline {
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
func cmdExecute(command []string, ip, op *chan *bytes.Buffer, stderr *bytes.Buffer) {
	if len(command) == 0 {
		e := commandNil
		stderr.Write([]byte(e))
		stdErrChan <- stderr
		return
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
