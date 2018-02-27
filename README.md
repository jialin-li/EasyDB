# EasyDB

### Building (with debug prints for `make test`)
`$ make`  

### Building for release (no debug prints)
`$ make release`  

### Running
`$ cd bin/`  
`$ ./master`  

### Testing
Will run every test case available in `tests/`  
`$ make & make test`  

___

### Directory structure
`bin/` executables  
`client/` src files for client  
`master/` src files for master  
`output/` output produced by tests in tests/  
`scripts/` scripts used by Makefile  
`server/` src files for server  
`shared/` src files for shared interfaces and util functions  
`tests/` tests cases and their expected output  
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
* Create a benchmark test
