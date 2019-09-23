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

//Pipeline ...
type Pipeline struct {
	commands [][]string
	stderr   *bytes.Buffer
}

//Run ...
func (p *Pipeline) Run() (string, error) {
	numOfCommands := len(p.commands)
	//Array of channels of bytes.Buffer one more than num of commands.

	channels, err := makeChannels(numOfCommands)
	if err != nil {
		fmt.Println("Somme issue in creating channels")
		return "", errors.New("Error in creating channel")
	}
	go firstStdin(&channels[0])

	temp := 0
	for index, cmd := range p.commands {
		//new struct for each command with input channels[i] as input and channels[i+1] as outputs
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

//NewPipeline ...
func NewPipeline(commands [][]string) *Pipeline {
	return &Pipeline{
		commands: commands,
		stderr:   new(bytes.Buffer),
	}
}

//makeChannels initializes n+1 number of channels
func makeChannels(n int) ([]chan *bytes.Buffer, error) {
	if n == 0 {
		return nil, errors.New(commandNil)
	}
	if n < 0 {
		return nil, errors.New(negativeChannels)
	}
	channels := make([]chan *bytes.Buffer, n+1)
	//var f *os.File
	for index := range channels {
		channels[index] = make(chan *bytes.Buffer)
	}
	//TODO: if any error could occur.
	return channels, nil
}
func cmdExecute(command []string, ip, op *chan *bytes.Buffer, stderr *bytes.Buffer) {
	//will give address of stderr to each nodes.
	if len(command) == 0 {
		e := commandNil
		stderr.Write([]byte(e))
		stdErrChan <- stderr
	}
	node := NewNode(command, stderr)
	//cmdStruct.SetCommand(command)
	node.Input(ip)
	node.Process()
	//TODO: need to properly take care of the output/err
	node.Output(op)
}
func firstStdin(channel *chan *bytes.Buffer) {

	b := new(bytes.Buffer)
	*channel <- b
}
