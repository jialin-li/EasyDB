package shared

type Args struct {
	// required args
	Msg   string
	Key   string
	Value string
}

// we might want to have multiple args struct for different calls

type Response struct {
	Result string
}

type Server interface {
	// Terminate a server
	Terminate() error
	// Ask the server to connect to another server
	Connect(args *Args, reply *Response) error
	// Send back the current key value store
	DumpStore(args *Args, reply *Response) error
	// Disconnect if we are clients of any other servers
	Disconnect(args *Args, reply *Response) error
	// No more new write requests will be send until the next write
	Stabilize(args *Args, reply *Response) error
	// Put a KV pair
	Put(args *Args, reply *Response) error
	// Get a Value based on a key
	Get(args *Args, reply *Response) error
}
