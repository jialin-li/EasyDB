# EasyDB

### Building
`$ make`

### Running
`$ cd bin/`  
`$ ./master`  

### Testing
Will run every test case available in `tests/`  
`$ make & make test`  

___

### Creating a test case
Tests go in the `tests/` directory. They can be in `tests/` or a subdirectory 1
level deep in `tests/`  
  
Inputs files are named as `{name}.txt` (empty lines are ignored)  
Ouput files are named as `{name}_out.txt`  
Both should be in the same directory.  

___

### Running the system manually
`$ cd bin/`  
`$ ./master`  
`$ ./server <server id>`  
`$ ./client <client id>`  

___

### TODO
* Look into switching to bufio.Scanner or fmt.Fscanf
* Make sure this can be compiled standalone
