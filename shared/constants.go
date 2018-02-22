package shared

const (
	// enums used to identify type of rpc client notifying the master
	ClientType = iota
	ServerType = iota

	// port numbers hopefully these are not used?
	MasterPort = 1234
	// common port based
	BasePort = 1240
	// client starts allocating port number from base port
	ClientStart = 0
	ServerStart = 1000

	MaxClient = 5

	ERR_DEP = "ERR_DEP"
	ERR_KEY = "ERR_KEY"
)
