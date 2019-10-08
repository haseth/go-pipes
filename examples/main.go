package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"

	"github.com/pipes"
)

func main() {
	states := []pipes.Commander{
		&GetURL{url: "https://curl.haxx.se"},
		&OsCommand{cmd: []string{"grep", "curl"}},
	}
	//Create a pipeline for executing commands
	pipe := pipes.NewPipeline(states)
	out, err := pipe.Run()
	if err.Error() != "" {
		fmt.Println(err.Error())
	}
	fmt.Println(out)
}

//OsCommand ...
type OsCommand struct {
	cmd []string
}

//Execute ...
func (o *OsCommand) Execute(stdin, stdout *bytes.Buffer) error {
	cmd := exec.Command(o.cmd[0], o.cmd[1:]...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	err := cmd.Run()
	if err != nil {
		return errors.New("Error in running the command:  " + err.Error())
	}
	return nil
}

//GetURL ...
type GetURL struct {
	url string
}

//Execute ...
func (g *GetURL) Execute(stdin, stdout *bytes.Buffer) error {
	cmd := exec.Command("curl", g.url)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	err := cmd.Run()
	if err != nil {
		return errors.New("Error in running the command:  " + err.Error())
	}
	return nil
}
