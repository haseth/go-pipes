package pipes

import (
	"bytes"
	"errors"
	"log"
	"os/exec"
	"strings"
	"testing"
)

func EqualSlice(a, b []string) bool {
	for in := range a {
		if a[in] != b[in] {
			return false
		}
	}
	return true
}
func TestCreatingNewNode(t *testing.T) {
	var testExec Executer

	t.Run("Testing Factory method for NodeState", func(t *testing.T) {
		_, err := NewNodeState(testExec)
		if err != nil {
			log.Fatalf("Error creating new node" + err.Error())
		}
	})
}
func TestNodeInput(t *testing.T) {
	// Node Setup
	node := &NodeState{}
	buff := new(bytes.Buffer)
	inpChan := make(chan *bytes.Buffer)

	// TEST1:
	t.Run("Checking by sending buffer address", func(t *testing.T) {
		//sending the buffer address to the IP channel
		go func(InpChan chan *bytes.Buffer, InpBuf *bytes.Buffer) {
			InpChan <- InpBuf
		}(inpChan, buff)

		err := node.Input(inpChan)
		if err != nil {
			log.Fatalf("Some error in receiving the input buffer.")
		}
		if node.stdin != buff {
			log.Fatalf("Err: difference in InpBuf and stdin got %s want %s", node.stdin, buff)
		}
	})

	// TEST2:
	t.Run("Checking by nil address", func(t *testing.T) {
		//sending the nil address to the IP channel
		go func(InpChan chan *bytes.Buffer, InpBuf *bytes.Buffer) {
			InpChan <- InpBuf
		}(inpChan, nil)

		err := node.Input(inpChan)
		if err == nil {
			log.Fatalf("Error should be received as nil address is passed in inpChan")
		}
		if err.Error() != stdInBufNil {
			log.Fatalf("string got %s want %s", err.Error(), buff)
		}
	})
}
func TestNodeOutput(t *testing.T) {
	// SETUP
	node := &NodeState{}
	node.stdout = new(bytes.Buffer)
	opChan := make(chan *bytes.Buffer)

	// TEST1:
	t.Run("Checking by sending buffer address", func(t *testing.T) {
		//sending the buffer address to the IP channel
		var OpBuff *bytes.Buffer
		go func(OpChan chan *bytes.Buffer, OpBuf *bytes.Buffer) {
			OpBuff = <-OpChan
		}(opChan, OpBuff)
		err := node.Output(opChan)
		if err != nil {
			log.Fatalf("Some error in sending the output buffer.")
		}
	})

	// SETUP
	node.stdout = nil
	// TEST2
	t.Run("Checking by sending nil address", func(t *testing.T) {
		//sending the nil address to the IP channel
		var OpBuff *bytes.Buffer
		go func(OpChan chan *bytes.Buffer, OpBuf *bytes.Buffer) {
			OpBuff = <-OpChan
		}(opChan, OpBuff)
		err := node.Output(opChan)
		if err == nil {
			log.Fatalf("Should have received error passing nill in stdout.")
		}
		if err.Error() != stdOutBufNil {
			log.Fatalf("string got %s want %s", err.Error(), stdOutBufNil)
		}
	})
}

func TestProcessNode(t *testing.T) {
	n := NodeState{}
	command := []string{"echo", "hello"}
	wrongCommand := []string{"lws"}

	n.stdin = new(bytes.Buffer)
	n.stdout = new(bytes.Buffer)
	n.executer = &OsExec{command}

	t.Run("Testing passing correct command", func(t *testing.T) {
		err := n.executer.Execute(n.stdin, n.stdout)
		if err != nil {
			log.Fatalf("We shouldn't face any error running the correct command: %s", command)
		}
	})

	n.executer = &OsExec{wrongCommand}
	t.Run("Testing passing incorrect command", func(t *testing.T) {
		err := n.executer.Execute(n.stdin, n.stdout)
		if err == nil {
			log.Fatalf("We should face any error running the incorrect command: %s", command)
		}
		if !strings.Contains(err.Error(), "not found") {
			log.Fatalf("string got %s want %s", err.Error(), "Command not found")
		}
	})
}

//TODO: Create a mock of it.
type OsExec struct {
	Cmds []string
}

func (o *OsExec) Execute(stdin, stdout *bytes.Buffer) error {
	cmd := exec.Command(o.Cmds[0], o.Cmds[1:]...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	err := cmd.Run()
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}
