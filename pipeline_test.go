package pipes

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestMakeLinks(t *testing.T) {
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

func TestCommandExecute(t *testing.T) {
	stdErr := new(bytes.Buffer)

	executer := &OsExec{Cmds: []string{"echo", "hello"}}
	node, err := NewNodeState(executer)
	if err != nil {
		t.Fatal(err)
	}

	ipChan := make(chan *bytes.Buffer)
	opChan := make(chan *bytes.Buffer)

	t.Run("Checking with correct command", func(t *testing.T) {
		execCmd(t, &ipChan, &opChan, stdErr, node, "hello")
	})

	executer = &OsExec{Cmds: []string{"lsw"}}
	node, err = NewNodeState(executer)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Checking with incorrect command", func(t *testing.T) {
		execCmd(t, &ipChan, &opChan, stdErr, node, commandNotFound)
	})
}
func execCmd(t *testing.T, ipChan, opChan *chan *bytes.Buffer, stdErr *bytes.Buffer, node Node, want string) {
	t.Helper()
	go firstStdin(*ipChan)
	go nodeExecute(node, *ipChan, *opChan, stdErr)
	gotBuffer := <-*opChan

	got := gotBuffer.String()
	if !strings.ContainsAny(got, want) {
		//check in stderr file
		e := make([]byte, 100)
		stdErr.Read(e)
		if !strings.Contains(string(e), want) {
			t.Fatalf("Should have received error executing command,  got %s want %s", string(e), want)
		}
		//t.Fatalf("pipeline execution failure got: %s want: %s", got, want)
	}

}

func TestRunWithCorrectCommand(t *testing.T) {
	commands := []Executer{
		&OsExec{[]string{"echo", "hello"}},
	}
	want := "hello"
	pipe := NewPipeline(commands)
	out, _ := pipe.Run()
	// if err.Error() != "" {
	// 	log.Fatalf("Error received, got %s, want %s", err.Error(), want)
	// }
	if !strings.Contains(out, want) {
		log.Fatalf("Error in running command, got %s, want %s", out, want)
	}
}

func TestRunWithInCorrectCommand(t *testing.T) {
	commands := []Executer{
		&OsExec{[]string{"lsw", "hello"}},
	}
	want := commandNotFound
	pipe := NewPipeline(commands)
	_, err := pipe.Run()
	if err.Error() == "" {
		log.Fatalf("Error received, got %s, want %s", err.Error(), want)
	}
	if !strings.Contains(err.Error(), want) {
		log.Fatalf("Error in running command, got %s, want %s", err.Error(), want)
	}
}

// commands = []Commander{}
// pipe = NewPipeline(commands)
// t.Run("Sending 0 commands", func(t *testing.T) {
// 	runCmd(pipe, commands, commandNil)
// })

// commands = []Commander{
// 	&OsExec{[]string{"lsw"}},
// }
// pipe = NewPipeline(commands)
// t.Run("Sending wrong commands", func(t *testing.T) {
// 	runCmd(pipe, commands, commandNotFound)
// })
//}

// func runCmd(pipe *Pipeline, commands []Commander, want string) {
// 	out, err := pipe.Run()
// 	if err.Error() != "" {
// 		if !strings.Contains(err.Error(), want) {
// 			log.Fatalf("Should have received error, got %s, want %s", err.Error(), want)
// 		} else {
// 			var b []byte
// 			pipe.stderr.Read(b)
// 			if !strings.Contains(string(b), want) {
// 				log.Fatalf("")
// 			}
// 		}
// 	} else {
// 		if !strings.Contains(out, want) {
// 			log.Fatalf("Contains Output mismatched, got %s want %s", out, want)
// 		}
// 	}

// }

// helper
func numChannels(t *testing.T, msg string, value int) {
	t.Helper()
	numofCmd := value
	gotChan, err := makeLinks(numofCmd)
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
