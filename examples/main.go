package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"

	pipes ".."
)

func main() {

	executers := []pipes.Executer{
		&GetURLExecuter{url: "https://curl.haxx.se"},
		&OSCmdExecuter{cmd: []string{"grep", "img"}},
	}

	pipe := pipes.NewPipeline(executers)

	out, err := pipe.Run()
	if err.Error() != "" {
		fmt.Println(err.Error())
	}

	fmt.Println(out)
}

//OSCmdExecuter implements pipes.Executer
type OSCmdExecuter struct {
	cmd []string
}

//Execute executes os command
func (o *OSCmdExecuter) Execute(stdin, stdout *bytes.Buffer) error {
	cmd := exec.Command(o.cmd[0], o.cmd[1:]...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	err := cmd.Run()
	if err != nil {
		return errors.New("Error in running the command:  " + err.Error())
	}
	return nil
}

//GetURLExecuter implements pipes.Executer
type GetURLExecuter struct {
	url string
}

//Execute gets url
func (g *GetURLExecuter) Execute(stdin, stdout *bytes.Buffer) error {
	cmd := exec.Command("curl", g.url)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	err := cmd.Run()
	if err != nil {
		return errors.New("Error in running the command:  " + err.Error())
	}
	return nil
}
