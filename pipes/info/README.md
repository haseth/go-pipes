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



# Synchronization 
Node1 -> Node2 --> Node3 

Synchronization is requried as :
- Node3 should use node2's output 
- Node2 should use node1's output 
- Node1 should use Input given to it. 

## Ways to do 

# Without channels: 
file > node1 > file then 
file > node2 > file then 
file > node3 > file.

Display the files. 

# With channels:

Single Channel (Wrong Approach)

->ch 
ch-> node1 -> ch 
ch-> node2 -> ch 
ch-> node3 -> ch 
(Anyone could take, breaking synchronization)

Number of Nodes +1 Channel
->ch1
ch1 -> Node1 ->ch2
ch2 -> Node2 ->ch3
ch3 -> Node3 ->ch4

Output <- ch4 

