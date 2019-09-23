package pipes

import (
	"bytes"
	"log"
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
	var stdErr *bytes.Buffer = new(bytes.Buffer)

	ipChan := make(chan *bytes.Buffer)
	opChan := make(chan *bytes.Buffer)

	t.Run("Checking with correct command", func(t *testing.T) {
		execCmd(t, &ipChan, &opChan, stdErr, []string{"echo", "hello"}, "hello")
	})

	t.Run("Checking with incorrect bash command", func(t *testing.T) {
		execCmd(t, &ipChan, &opChan, stdErr, []string{"lw", "test.txt"}, commandNotFound)
	})

	// t.Run("Checking with empty command", func(t *testing.T) {
	// 	execCmd(t, &ipChan, &opChan, stdErr, []string{}, commandNil)
	// })
}
func execCmd(t *testing.T, ipChan, opChan *chan *bytes.Buffer, stdErr *bytes.Buffer, cmd []string, want string) {
	t.Helper()
	go firstStdin(ipChan)
	go cmdExecute(cmd, ipChan, opChan, stdErr)

	for {
		select {
		case gotBuffer := <-*opChan:
			got := gotBuffer.String()
			if got != want {
				t.Fatalf("pipeline execution failure got: %s want: %s", got, want)
			}
		case errBuffer := <-stdErrChan:
			got := errBuffer.String()
			if got != want {
				t.Fatalf("Should have received error,  got %s want %s", got, want)
			}
		}
	}
}

func TestRun(t *testing.T) {
	commands := [][]string{
		[]string{"echo", "hello seth"},
		[]string{"grep hello"},
	}
	pipe := NewPipeline(commands)

	t.Run("Sending n correct commands", func(t *testing.T) {
		runCmd(pipe, commands, "hello")
	})
	t.Run("Sending n-1 correct commands  and one incorrect", func(t *testing.T) {
		runCmd(pipe, commands, commandNotFound)
	})
	t.Run("Sending 0 commands", func(t *testing.T) {
		runCmd(pipe, commands, commandNil)
	})

}
func runCmd(pipe *Pipeline, commands [][]string, want string) {
	out, err := pipe.Run()
	if err != nil {
		if err.Error() != want {
			log.Fatalf("Should have received error, got %s want %s", err.Error(), want)
		}
	} else {
		if out != want {
			log.Fatalf("Output mismatched, got %s want %s", out, want)
		}
	}

}
