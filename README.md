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

