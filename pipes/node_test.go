package pipes

import (
	"bytes"
	"log"
	"reflect"
	"strings"
	"testing"
)

func TestNewNode(t *testing.T) {
	testStdErrBuff := new(bytes.Buffer)
	testCmd := []string{"cat", "test.txt"}
	t.Run("Running passsing the command", func(t *testing.T) {
		node, err := NewNode(testCmd, testStdErrBuff)
		if err != nil {
			log.Fatalf("Error creating new node" + err.Error())
		}
		p := string(strings.Join(cmd, " "))
		if == node.stdin {
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
	// t.Run("Running passsing no buff", func(t *testing.T) {
	// 	node, err := NewNode(testCmd, testStdErrBuff)
	// 	if err == nil {
	// 		log.Fatalf("We should have received the error passing empty command")
	// 	}
	// 	if err.Error() != commandNil {
	// 		log.Fatalf("Incorrect error received got %s want %s", err.Error(), commandNil)
	// 	}
	// })

	//test 1: what if no cmd
	//test2: what if no buf
}
func TestSetCommands(t *testing.T) {
	cmd := []string{"testing", "command"}
	node := &Node{}
	//test1: passing command
	//test2: passsing no command
	t.Run("Testing with random command", func(t *testing.T) {
		err := node.SetCommand(cmd)
		if err != nil {
			log.Fatalf("Error setting the command")
		}
		if reflect.DeepEqual(node.cmd, cmd) {
			log.Fatalf("Setting the command not working in SetCommands got %s want %s", node.cmd, cmd)
		}
	})
	cmd = []string{}
	t.Run("Testing with no command", func(t *testing.T) {
		err := node.SetCommand(cmd)
		if err == nil {
			log.Fatalf("We should have got error for setting empty command.")
		} else {

		}
	})

}
func TestNodeInput(t *testing.T) {
	node := &Node{}
	InpBuf := new(bytes.Buffer)
	var InpChan chan *bytes.Buffer
	InpChan = make(chan *bytes.Buffer)

	//TEST1: checking by sending the buffer address
	//TEST2: checking by sending nill address.
	// Test3: Not sending.. deadlock state handling

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

// func TestNodeOutput(t *testing.T) {
// 	node := &Node{}
// 	var OutChan chan *bytes.Buffer

// 	go func(node *Node, OutChan chan *bytes.Buffer) {
// 		OutBuff := new(bytes.Buffer)
// 		node.stdout = OutBuff
// 		node.Output(&OutChan)
// 	}(node, OutChan)
// 	b := <-OutChan

// 	CheckTest(t, node.stdout, b, "Error in setting the address of buffer.")
// }
// func TestNode(t *testing.T) {
// 	cmd := []string{"cat", "test.txt"}

// 	stdErr := new(bytes.Buffer)
// 	node, err := NewNode(cmd, stdErr)
// 	if err != nil {
// 		log.Fatalf("Error in allocating node")
// 	}

// 	InpBuff := new(bytes.Buffer)
// 	//InpBuff.Write([]byte{""})
// 	OutBuff := new(bytes.Buffer)
// 	//var InpChannel chan *bytes.Buffer

// 	// //we will put the address of the InpBuff on the channel..
// 	// go func(InpBuff bytes.Buffer, InpChannel chan *bytes.Buffer){
// 	// 	InpChannel <- & InpBuff
// 	// }(InpBuff,InpChannel)
// 	node.stdin = InpBuff
// 	node.stdout = OutBuff
// 	node.Process()

// 	CheckTest(t, *OutBuff, "harsh seth", "Output string not matching after piping.")
// 	//test 1: what if wrong command given
// 	//test 2: what if no command given
// }

// func CheckTest(t *testing.T, got, want interface{}, err string) {
// 	t.Helper()
// 	if got != want {
// 		t.Fatalf(err+"got %s wanted %s", got, want)
// 	}
// 	// switch v := got.(type) {
// 	// case string:
// 	// 	//check strings
// 	// 	i
// 	// 	fmt.Printf("Twice %v is %v\n", v, v*2)
// 	// case string:
// 	// 	fmt.Printf("%q is %v bytes long\n", v, len(v))
// 	// default:
// 	// 	fmt.Printf("I don't know about type %T!\n", v)
// 	// }
// }
