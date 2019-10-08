<h2> GO-PIPES </h2>
<HR>

It is a library to implement Linux pipes using GOlang. 

There are two branched: 

1. Master: 
- Here we use *os.File as the channel type. 
- But there are some issue comming as IP channel holds the *File and doesn't allow to write in O/P Channel. 


2. Master2: 
- Here we use *bytes.Buffer as the channel type. 
- We are getting the resulting output, but still code needs to be reafactored. 


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
		&OsCommand{cmd: []string{"curl", "google.com"}},
}

pipe := pipes.NewPipeline(states)
out:= pipe.Run()
```



