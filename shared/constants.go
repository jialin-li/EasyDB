package shared

const (
	// enums used to identify type of rpc client notifying the master
	ClientType = iota
	ServerType = iota

	// port numbers hopefully these are not used?
	MasterPort = 1234
	// client ports are determined by 1240 + clientId
	ClientPort = 1240
	// server ports are determined by 1250 + serverId
	ServerPort = 1250

	ERR_DEP = "ERR_DEP"
	ERR_KEY = "ERR_KEY"
)
