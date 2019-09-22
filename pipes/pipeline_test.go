package pipes

import (
	"bytes"
	"testing"
)

func TestMakeChannel(t *testing.T) {
	numOfCmd := 10
	gotChan, err := makeChannels(numOfCmd)
	got := len(gotChan)
	want := numOfCmd + 1
	if err != nil {
		t.Fatalf("Error in creating channels")
	}

	//should be numOfCmd +1  channels
	if got != want {
		t.Fatalf("Num of channels should be num of cmd +1 got: %d want : %d ", got, want)
	}
}

func TestCommandExecute(t *testing.T) {
	var ipChan chan *bytes.Buffer
	var opChan chan *bytes.Buffer

	firstStdin(&ipChan)
	go cmdExecute([]string{"echo", "hello"}, ipChan, opChan, 1)
	gotBuffer := <-opChan
	got := gotBuffer.String()
	want := "hello"
	if got != want {
		t.Fatalf("pipeline execution failure got: %s want: %s", got, want)
	}

}

func TestRun(t *testing.T) {

}
func TestfirstStdin(t *testing.T) {

}
