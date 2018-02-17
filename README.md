# EasyDB

### Building
`$ make`

### Running
`$ cd bin/`  
`$ ./master`  

___

### Running the system manually
`$ cd bin/`  
`$ ./master`  
`$ ./server <server id>`  
`$ ./client <client id>`  

___

### TODO
* Look into switching to bufio.Scanner or fmt.Fscanf
* Extract internal functions from client and server's shared interface
implementaion
* Use TcpConn instead?
* Make sure this can be compiled standalone
