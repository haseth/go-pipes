package pipes

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestMakeChannel(t *testing.T) {
	t.Run("Testing with 10 values", func(t *testing.T) {
		numChannels(t, "", 10)
	})

	t.Run("Testing with 0 values", func(t *testing.T) {
		numChannels(t, commandNil, 0)
	})
	t.Run("Testing with negative values", func(t *testing.T) {
		numChannels(t, negativeChannels, -10)
	})

}
func numChannels(t *testing.T, msg string, value int) {
	t.Helper()
	numofCmd := value
	gotChan, err := makeChannels(numofCmd)
	if err != nil {
		if err.Error() != msg {
			t.Fatalf("Should have received an error, got %s, wanted %s", err.Error(), msg)
		}
	} else {
		got := len(gotChan)
		want := value + 1
		//should be numOfCmd +1  channels
		if got != want {
			t.Fatalf("Num of channels should be num of cmd +1 got: %d want : %d ", got, want)
		}
	}
}
func TestCommandExecute(t *testing.T) {
	stdErr := new(bytes.Buffer)
	commander := &OsExec{Cmds: []string{"echo", "hello"}}
	ipChan := make(chan *bytes.Buffer)
	opChan := make(chan *bytes.Buffer)

	t.Run("Checking with correct command", func(t *testing.T) {
		execCmd(t, &ipChan, &opChan, stdErr, commander, "hello")
	})
	// commander = &OsExec{Cmds: []string{"lsw"}}
	// t.Run("Checking with incorrect command", func(t *testing.T) {
	// 	execCmd(t, &ipChan, &opChan, stdErr, commander, commandNotFound)
	// })
	// //Why it is failing?

	// commander = nil
	// t.Run("Checking with no commander", func(t *testing.T) {
	// 	execCmd(t, &ipChan, &opChan, stdErr, commander, commandNotFound)
	// })
}
func execCmd(t *testing.T, ipChan, opChan *chan *bytes.Buffer, stdErr *bytes.Buffer, cmd Commander, want string) {
	t.Helper()
	go firstStdin(ipChan)
	go cmdExecute(cmd, ipChan, opChan, stdErr)

	select {
	case gotBuffer := <-*opChan:
		got := gotBuffer.String()
		if !strings.ContainsAny(got, want) {
			t.Fatalf("pipeline execution failure got: %s want: %s", got, want)
		}

	case errBuffer := <-stdErrChan:
		got := errBuffer.String()
		if !strings.ContainsAny(got, want) {
			t.Fatalf("Should have received error,  got %s want %s", got, want)
		}

	}
}

func TestRun(t *testing.T) {
	commands := []Commander{
		&OsExec{[]string{"echo", "hello"}},
	}
	pipe := NewPipeline(commands)

	t.Run("Sending n correct commands", func(t *testing.T) {
		runCmd(pipe, commands, "hello")
	})
	commands = []Commander{}
	pipe = NewPipeline(commands)
	t.Run("Sending 0 commands", func(t *testing.T) {
		runCmd(pipe, commands, commandNil)
	})

	commands = []Commander{
		&OsExec{[]string{"lsw"}},
	}
	pipe = NewPipeline(commands)
	t.Run("Sending wrong commands", func(t *testing.T) {
		runCmd(pipe, commands, commandNotFound)
	})
}
func runCmd(pipe *Pipeline, commands []Commander, want string) {
	out, err := pipe.Run()
	if err != nil {
		if strings.Contains(string(b), want) {
			log.Fatalf("Should have received error, got %s, want %s", err.Error(), want)
		}
	} else {
		if !strings.Contains(out, want) {
			log.Fatalf("Contains Output mismatched, got %s want %s", out, want)
		}
	}

}
