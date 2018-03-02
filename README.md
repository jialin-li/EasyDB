# EasyDB

### Student Information
Will Lin - wl7524 - wlsaidhi  
Jialin Li - jl48629 - jialinli 

### Set up  


### Building (with debug prints for `make test`)
`$ make`  

### Building for release (no debug prints)
`$ make release`  

___

### Running the system through master
`$ make release` in the EasyDB directory  
`$ cd bin`   
`$ ./master` and you can type in your commands :D  

### Testing
Will run every test case available in `tests/`  

Tests with debug output:  
`$ make test`  

Tests without debug output:  
`$ make rtest`  
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

### Protocol 

In order to keep track of the order of writes for each key, we generate an independent timestamp for each key.
We use vector clock to represent the timestamp. The vector clock has 10 entries, top 5 are for the five clients 
and bottom 5 are for the servers. Each client/server is assigned an unique index in the vector clock on creation.
When a client first issues a request to write a key value pair to the server, it first generates a vector clock 
for the given key with all entries initialized to 0. It then increments its the value at its index, and sends 
this timestamp along with the key value pair to a server it's connected to. Upon receiving this request, the server adds/updates the key value pair to its database and update the timestamp for this key with the timestamp provided by the client. We update timestamps (A with B) by comparing the value at each index of the timestamp,
if B's value at index i is greater than A's value at index i, then we update A's value at index i to B's value.
Once the server updates its timestamp with the client provided one, it increments the timestamp by 1 at its corresponding position, saves it in its database with this kv pair and responds to the client this timestamp.
When client receives the timestamp from the server, it updates(synchronize) its timestamp to the one returned by the server. Every client has a map of keys that it has put into the kvstore and the corresponding timestamps.

When a client issues a read, it sends its timestamp in the read request. The responding server first checks to
see if the key exists in its database, if not, it responds with ERR_KEY as the value. If so, it checks to see if 
the client timestamp is later than its current timestamp for that key. Timestamp A is later than timestamp B if 
A's value at an index is greater than B's value at the same index. We start this check from entry 0. 
This implicitly places priority on the clients that are created early. If client has a timestamp that's later 
than server's timestamp, we return ERR_DEP. Else, we respond with the value and the server's timestamp to the 
client. And upon receiving a successfully completed request, client updates its timestamp to sync with server.

In stabilize, each server communicates with the other servers in its partition and send its entire database.
Each database includes kv pairs and the corresponding timestamps. When one server receives the another server's 
database, it goes through all the entries, insert any new entries and update existing entries if the other 
server's timestamp for that key is later than its timestamp. To ensure we can propogate all updates from 
different kinds of partitions, for each stabilize call, we actually do multiple rounds of stabilize to the 
servers to ensure all the updates propagate through (in case of a chain connection).

### TODO
* Make sure this can be compiled standalone
