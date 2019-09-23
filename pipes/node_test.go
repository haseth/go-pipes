package pipes

import (
	"bytes"
	"log"
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
func TestNewNode(t *testing.T) {
	testStdErrBuff := new(bytes.Buffer)
	testCmd := []string{"cat", "test.txt"}
	t.Run("Running passsing the command", func(t *testing.T) {
		node, err := NewNode(testCmd, testStdErrBuff)
		if err != nil {
			log.Fatalf("Error creating new node" + err.Error())
		}
		//TODO: is the right way to test the slices?
		if !EqualSlice(node.cmd, testCmd) {
			log.Fatalf("Error in intializing the node.cmd got %s want %s", node.cmd, testCmd)
		}

		if node.stderr != testStdErrBuff {
			log.Fatalf("Error in intializing the buffer got %v want %v", node.stderr, testStdErrBuff)
		}
	})
	testCmd = []string{}
	t.Run("Running passsing no command", func(t *testing.T) {
		_, err := NewNode(testCmd, testStdErrBuff)
		if err == nil {
			log.Fatalf("We should have received the error passing empty command")
		}
		if err.Error() != commandNil {
			log.Fatalf("Incorrect error received got %s want %s", err.Error(), commandNil)
		}
	})
	t.Run("Running passsing no buff", func(t *testing.T) {
		_, err := NewNode(testCmd, testStdErrBuff)
		if err == nil {
			log.Fatalf("We should have received the error passing empty command")
		}
		if err.Error() != commandNil {
			log.Fatalf("Incorrect error received got %s want %s", err.Error(), commandNil)
		}
	})

	//test 1: what if no cmd
	//test2: what if no buf
}
func TestNodeInput(t *testing.T) {
	node := &Node{}
	InpBuf := new(bytes.Buffer)
	var InpChan chan *bytes.Buffer
	InpChan = make(chan *bytes.Buffer)

	//TEST1: checking by sending the buffer address
	//TEST2: checking by sending nill address.
	//Test3: Not sending.. deadlock state handling

	t.Run("Checking by sending buffer address", func(t *testing.T) {
		//sending the buffer address to the IP channel
		go func(InpChan chan *bytes.Buffer, InpBuf *bytes.Buffer) {
			InpChan <- InpBuf
		}(InpChan, InpBuf)

		err := node.Input(&InpChan)
		if err != nil {
			log.Fatalf("Some error in receiving the input buffer.")
		}
		if node.stdin != InpBuf {
			log.Fatalf("Err: difference in InpBuf and stdin got %s want %s", node.stdin, InpBuf)
		}
	})
	t.Run("Checking by nil address", func(t *testing.T) {
		//sending the nil address to the IP channel
		go func(InpChan chan *bytes.Buffer, InpBuf *bytes.Buffer) {
			InpChan <- InpBuf
		}(InpChan, nil)

		err := node.Input(&InpChan)
		if err == nil {
			log.Fatalf("Error should be received ")
		}
		if err.Error() != stdErrBufNil {
			log.Fatalf("string got %s want %s", err.Error(), InpBuf)
		}
	})
}
