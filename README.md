# golook
golook on any go application code directory to get stats a high level look of the code.

This is an application that will show you stats about a golang project - like line count -- number of functions, structs etc.

Though golang has godoc that shows you details of a project that are Public but if you want to start contributing to some project 
and want to learn what it really does internally this is a good place to start.

This project shows you a High Level picture of a project and gives you clues about what the application you are trying to 
work with is about. Think of this like a long index listing or a skeleton of an application. Once you know what the application is about
you can deep dive right into the right portions of the code and work with it.

#### Features
This is what you can see on a Package level:
- Files used in this package and details of
    - Number of lines
    - Counts of Functions, Structs & Methods defined
    - Breakdown of Public / Private   
- Struct definitions
    - Public and private methods on them with their method signatures
- Functions list (Public / Private) with number of lines count    
        
#### Installing
you can clone this code and just run: <br/>
`go install golook.go`

Now you will be able to run golook from anywhere.<br/>
`golook` 

The application takes the current working directory as the starting point and then recursively parses all folders for go files.

It leverages the native golang ast to parse the files.
    