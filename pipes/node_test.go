package pipes

// func TestNewNode(t *testing.T) {
// 	node := &Node{}
// 	testIPBuff := new(bytes.Buffer)
// 	testOPBuff := new(bytes.Buffer)
// 	testStdErrBuff := new(bytes.Buffer)

// 	testCmd := []string{"cat", "test.txt"}
// 	node = NewNode(testCmd, testStdErrBuff)

// 	CheckTest(t, node.cmd, testCmd, "Error in intializing the node.cmd.")
// 	CheckTest(t, node.stdin, testIPBuff, "Error in intializing the node.stdin buffer.")
// 	CheckTest(t, node.stdout, testOPBuff, "Error in intializing the node.stdout buffer.")

// 	//test 1: what if no cmd
// }
// func TestSetCommands(t *testing.T) {
// 	cmd := []string{"testing", "command"}
// 	node := &Node{}
// 	node.SetCommand(cmd)
// 	CheckTest(t, node.cmd, cmd, "Setting the command not working in SetCommands")
// }
// func TestNodeInput(t *testing.T) {
// 	node := &Node{}
// 	InpBuf := new(bytes.Buffer)
// 	var InpChan chan *bytes.Buffer
// 	go func(InpBuf *bytes.Buffer, InpChan chan *bytes.Buffer) { InpChan <- InpBuf }(InpBuf, InpChan)
// 	node.Input(InpChan)
// 	CheckTest(t, node.stdin, InpBuf, "Error in setting the address of buffer.")

// }
// func TestNodeOutput(t *testing.T) {
// 	node := &Node{}
// 	var OutChan chan *bytes.Buffer

// 	go func(node *Node, OutChan chan *bytes.Buffer) {
// 		OutBuff := new(bytes.Buffer)
// 		node.stdout = OutBuff
// 		node.Output(OutChan)
// 	}(node, OutChan)
// 	b := <-OutChan

// 	CheckTest(t, node.stdout, b, "Error in setting the address of buffer.")
// }
// func TestNode(t *testing.T) {
// 	cmd := []string{"cat", "test.txt"}
// 	node := &Node{cmd: cmd}

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
// 	node.SetCommand(cmd)
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
