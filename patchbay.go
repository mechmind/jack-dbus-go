package jackdbus

type Patchbay struct {
	object
	iface string

	eventRouter *eventRouter
}

// methods

func (pb *Patchbay) GetAllPorts() ([]string, error) {
	ports := []string{}
	err := call(pb.Object(), pb.iface, "GetAllPorts", 0, nil, args{&ports})
	if err != nil {
		return nil, err
	}
	return ports, nil
}

func (pb *Patchbay) GetGraph(lastVersion uint64) (uint64, *Graph, error) {
	var graph Graph
	var version uint64
	var connections []connection
	err := call(pb.Object(), pb.iface, "GetGraph", 0, args{&lastVersion},
		args{&version, &graph.ClientsAndPorts, &connections})
	if err != nil {
		return 0, nil, err
	}

	graph.Connections = parseConnections(connections)
	return version, &graph, nil
}

func (pb *Patchbay) ConnectPortsByID(port1, port2 uint64) error {
	return call(pb.Object(), pb.iface, "ConnectPortsByID", 0, args{port1, port2}, nil)
}

func (pb *Patchbay) ConnectPortsByName(port1, port2 string) error {
	return call(pb.Object(), pb.iface, "ConnectPortsByName", 0, args{port1, port2}, nil)
}

func (pb *Patchbay) DisconnectPortsByID(port1, port2 uint64) error {
	return call(pb.Object(), pb.iface, "DisconnectPortsByID", 0, args{port1, port2}, nil)
}

func (pb *Patchbay) DisconnectPortsByName(port1, port2 string) error {
	return call(pb.Object(), pb.iface, "DisconnectPortsByName", 0, args{port1, port2}, nil)
}

func (pb *Patchbay) DisconnectPortsByConnectioID(connID uint64) error {
	return call(pb.Object(), pb.iface, "DisconnectPortsByConnectionID", 0, args{connID}, nil)
}

func (pb *Patchbay) GetClientPID(clientID uint64) (pid int64, err error) {
	err = call(pb.Object(), pb.iface, "GetClientPID", 0, args{clientID}, args{&pid})
	return pid, err
}

// signal handlers

func (pb *Patchbay) OnGraphChanged(handler func(graphID uint64)) (detach func() error, err error) {
	return pb.eventRouter.addHandler(pb.Object(), pb.iface, "GraphChanged", handler)
}

func (pb *Patchbay) OnClientAppeared(handler func(uint64, uint64, string)) (detach func() error, err error) {
	return pb.eventRouter.addHandler(pb.Object(), pb.iface, "ClientAppeared", handler)
}

func (pb *Patchbay) OnClientDisappeared(handler func(uint64, uint64, string)) (detach func() error, err error) {
	return pb.eventRouter.addHandler(pb.Object(), pb.iface, "ClientDisappeared", handler)
}

func (pb *Patchbay) OnPortAppeared(handler func(uint64, uint64, string, uint64, string, uint32, uint32)) (detach func() error, err error) {
	return pb.eventRouter.addHandler(pb.Object(), pb.iface, "PortAppeared", handler)
}

func (pb *Patchbay) OnPortDisappeared(handler func(uint64, uint64, string, uint64, string)) (detach func() error, err error) {
	return pb.eventRouter.addHandler(pb.Object(), pb.iface, "PortDisappeared", handler)
}

func (pb *Patchbay) OnPortsConnected(handler func(uint64, uint64, string, uint64, string, uint64, string, uint64)) (detach func() error, err error) {
	return pb.eventRouter.addHandler(pb.Object(), pb.iface, "PortsConnected", handler)
}

func (pb *Patchbay) OnPortsDisconnected(handler func(uint64, uint64, string, uint64, string, uint64, string, uint64)) (detach func() error, err error) {
	return pb.eventRouter.addHandler(pb.Object(), pb.iface, "PortsDisconnected", handler)
}
