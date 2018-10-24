package jackdbus

type ClientInfo struct {
	ClientID   uint64
	ClientName string
	Ports      []PortInfo
}

type PortInfo struct {
	PortID   uint64
	PortName string
	Flags    uint32
	Type     uint32
}

type ConnectionHandle struct {
	ClientID   uint64
	ClientName string
	PortID     uint64
	PortName   string
}

type Connection struct {
	Client1Info  ConnectionHandle
	Client2Info  ConnectionHandle
	ConnectionID uint64
}

type connection struct {
	Client1ID       uint64
	Client1Name     string
	Client1PortID   uint64
	Client1PortName string
	Client2ID       uint64
	Client2Name     string
	Client2PortID   uint64
	Client2PortName string
	ConnectionID    uint64
}

func parseConnections(src []connection) []Connection {
	result := []Connection{}
	for _, raw := range src {
		result = append(result, Connection{
			Client1Info: ConnectionHandle{
				ClientID:   raw.Client1ID,
				ClientName: raw.Client1Name,
				PortID:     raw.Client1PortID,
				PortName:   raw.Client1PortName,
			},
			Client2Info: ConnectionHandle{
				ClientID:   raw.Client2ID,
				ClientName: raw.Client2Name,
				PortID:     raw.Client2PortID,
				PortName:   raw.Client2PortName,
			},
			ConnectionID: raw.ConnectionID,
		})
	}
	return result
}

type Graph struct {
	ClientsAndPorts []ClientInfo
	Connections     []Connection
}
