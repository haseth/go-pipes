# GO-Pipes 

GO-Pipes is a library to implement ```Linux pipes``` using Go. 

# Usage  

Each node should have a ```Execute``` method which should define how data will be processed by node taking from ```stdin``` and putting back to ```stdout```. 

```
type Executer interface {
	Execute(stdin, stdout *bytes.Buffer) error
}
```

```
type myNodeExecuter struct{}

func (s *myNodeExecuter) Execute(stdin, stdout *bytes.Buffer){
    //Implement taking data from stdin and put in stdout.
}   
```

After filtering from all the nodes in pipeline, it will display ```error (stderr)``` and ```output (stdout)``` of the pipe. 

# Example

We will perform following operations:

1. ``` Get (wget)``` data from ```google.com```
2. ```Filter``` links from data


```Node 1 ```

	
	type GetURLExecuter struct {
		url string
	}

	//Execute will take input/output from stdin/stdout and get data from URL. 
	func (g *GetURLExecuter) Execute(stdin, stdout *bytes.Buffer) error {
		cmd := exec.Command("curl", g.url)
		cmd.Stdin = stdin
		cmd.Stdout = stdout
		err := cmd.Run()
		if err != nil {
			return errors.New("Error in running the command:  " + err.Error())
		}
		return nil
	}
	

```Node 2 ```

	
	//OsCommandExecuter takes array of commands
	type OsCommandExecuter struct {
		cmd []string
	}

	//Execute takes input/output from stdin/stdout, runs the os command. 
	func (o *OsCommandExecuter) Execute(stdin, stdout *bytes.Buffer) error {
		cmd := exec.Command(o.cmd[0], o.cmd[1:]...)
		cmd.Stdin = stdin
		cmd.Stdout = stdout
		err := cmd.Run()
		if err != nil {
			return errors.New("Error in running the command:  " + err.Error())
		}
		return nil
	}

	
 
Link the defined nodes and run the pipe. 

```
nodes := []pipes.Executer{
			&GetURL{url: "https://curl.haxx.se"},
			&OsCommand{cmd: []string{"grep", "curl"}},
}

pipe := pipes.NewPipeline(nodes)
out:= pipe.Run()
```

``` Note: Order of nodes matter ```


# Fetaure 

1. Pipe of nodes contains ```common error buffer```, failure of single node wouldn't ```break``` the pipe. 
2. Each node ```waits``` for previous node to finish ```order is preserved```.



# Reference 

- [Pipes](http://www.linfo.org/pipes.html): A Brief Introduction


