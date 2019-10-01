## GO-Pipes
Implementing pipes in Go. 


# Example of Nodes
wget "something.com" <pipe> filterLinks <pipe> whatElseYouWantToDo 
	Node 1         |        Node2      |      Node3  

Display stdout.

# Usage
Define your own state: 

```
type NodeState struct{}

func (s *NodeState) Execute(stdin, stdout *bytes.Buffer){
    //Implementation goes here using stdin and put in stdout.
}   
```
Pipe different states 

```
states := []pipes.Commander{
		&OsCommand{cmd: []string{"curl", "google.com"}},
}

pipe := pipes.NewPipeline(states)
out:= pipe.Run()
```



