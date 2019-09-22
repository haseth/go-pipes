package pipes

import (
	"bytes"
	"fmt"
)

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
}

//Run ...
func (p *Pipeline) Run() (string, error) {
	numOfCommands := len(p.commands)
	//exec.Command("ls").String()
	//Array of channels of bytes.Buffer one more than num of commands.
	channels, err := makeChannels(numOfCommands)
	if err != nil {
		fmt.Println("Somme issue in creating channels")
		panic(err)
	}
	go firstStdin(&channels[0])
	//first stdin file
	// go func(channels []chan *bytes.Buffer) {
	// 	b := new(bytes.Buffer)
	// 	channels[0] <- b
	// }(channels)

	temp := 0
	for index, cmd := range p.commands {
		//new struct for each command with input channels[i] as input and channels[i+1] as outputs
		go cmdExecute(cmd, channels[index], channels[index+1], index)
		temp = index + 1
	}
	//last channel should retrieve the output.
	f := <-channels[temp]
	return f.String(), nil
}

//NewPipeline ...
func NewPipeline(commands [][]string) *Pipeline {
	return &Pipeline{commands: commands}
}

//makeChannels initializes n+1 number of channels
func makeChannels(n int) ([]chan *bytes.Buffer, error) {
	channels := make([]chan *bytes.Buffer, n+1)
	//var f *os.File
	for index := range channels {
		channels[index] = make(chan *bytes.Buffer)
	}
	//TODO: if any error could occur.
	return channels, nil
}
func cmdExecute(command []string, ip, op chan *bytes.Buffer, index int) {
	cmdStruct := &Node{}
	cmdStruct.SetCommand(command)
	cmdStruct.Input(ip)
	cmdStruct.Process(index)
	cmdStruct.Output(op)
}
func firstStdin(channel *chan *bytes.Buffer) {
	//some pointer checking
	// if channel != nil {
	// 	b := new(bytes.Buffer)
	// 	*channel <- b
	// 	return nil
	// } else {
	// 	return "Did not get address for first stdin channel"
	// }
	b := new(bytes.Buffer)
	*channel <- b
}
