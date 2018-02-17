package shared

type Args struct {
	// required args
	Msg   string
	Key   string
	Value string
	Time  [ClockLen]int
}

type NotifyArgs struct {
	Type int
	ID   int
}

// we might want to have multiple args struct for different calls

type Response struct {
	Result string
	Time   [ClockLen]int
}

type Server interface {
	// Terminate a server
	Terminate(args *Args, reply *Response) error
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

// exposed so master can talk to client through this
type Client interface {
	// Ask the client to connect to another client or server
	Connect(args *Args, reply *Response) error
	// Disconnect ourselves with the specified client or server
	Disconnect(args *Args, reply *Response) error
	// Put a KV pair
	Put(args *Args, reply *Response) error
	// Get a Value based on a key
	Get(args *Args, reply *Response) error

	// No more new write requests will be send until the next write
	// Stabilize(args *Args, reply *Response) error
}

type Master interface {
	// Notify the master that it's up
	Notify(args *NotifyArgs, reply *Response) error
}
