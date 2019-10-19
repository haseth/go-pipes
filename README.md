<h2> GO-PIPES </h2>
<HR>

It is a library to implement Linux pipes using GOlang. 

# Example of Nodes
``` 
wget "something.com" <pipe> filterLinks <pipe> whatElseYouWantToDo 
```
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
			&GetURL{url: "https://curl.haxx.se"},
			&OsCommand{cmd: []string{"grep", "curl"}},
}

pipe := pipes.NewPipeline(states)
out:= pipe.Run()
```



