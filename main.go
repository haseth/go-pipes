package main

import (
	"bytes"
	"errors"
	"fmt"
	"help/go-pipes/pipes"
	"os/exec"
)

func main() {
	states := []pipes.Commander{
		&OsCommand{cmd: []string{"sent", "google"}}, //wrong OS command
		&OsCommand{cmd: []string{"curl", "google.com"}},
		&OsCommand{cmd: []string{"grep", "google"}},
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
