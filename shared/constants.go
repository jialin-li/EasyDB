package shared

const (
	ClientType = iota
	ServerType = iota

	// port numbers hopefully these are not used?
	// client ports are determined by 1240 + clientId
	ClientPort = 1240
	// server ports are determined by 1250 + serverId
	ServerPort = 1250
)
