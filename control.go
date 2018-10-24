package jackdbus

type Control struct {
	object
	iface string

	eventRouter *eventRouter
}

// methods

func (control *Control) IsStarted() (bool, error) {
	var started bool
	err := call(control.Object(), control.iface, "IsStarted", 0, nil, args{&started})
	return started, err
}

func (control *Control) StartServer() error {
	return call(control.Object(), control.iface, "StartServer", 0, nil, nil)
}

func (control *Control) StopServer() error {
	return call(control.Object(), control.iface, "StopServer", 0, nil, nil)
}

func (control *Control) AddSlaveDriver(driver string) error {
	return call(control.Object(), control.iface, "AddSlaveDriver", 0, args{driver}, nil)
}

func (control *Control) RemoveSlaveDriver(driver string) error {
	return call(control.Object(), control.iface, "RemoveSlaveDriver", 0, args{driver}, nil)
}

func (control *Control) GetBufferSize() (uint32, error) {
	var bufferSize uint32
	err := call(control.Object(), control.iface, "GetBufferSize", 0, nil, args{&bufferSize})
	return bufferSize, err
}

func (control *Control) GetLatency() (float64, error) {
	var latency float64
	err := call(control.Object(), control.iface, "GetLatency", 0, nil, args{&latency})
	return latency, err
}

func (control *Control) GetLoad() (float64, error) {
	var load float64
	err := call(control.Object(), control.iface, "GetLoad", 0, nil, args{&load})
	return load, err
}

func (control *Control) GetSampleRate() (uint32, error) {
	var sampleRate uint32
	err := call(control.Object(), control.iface, "GetSampleRate", 0, nil, args{&sampleRate})
	return sampleRate, err
}

func (control *Control) GetXruns() (uint32, error) {
	var xruns uint32
	err := call(control.Object(), control.iface, "GetXruns", 0, nil, args{&xruns})
	return xruns, err
}

func (control *Control) IsRealtime() (bool, error) {
	var isRealtime bool
	err := call(control.Object(), control.iface, "IsRealtime", 0, nil, args{&isRealtime})
	return isRealtime, err
}

func (control *Control) LoadInternal(internal string) error {
	return call(control.Object(), control.iface, "LoadInternal", 0, args{internal}, nil)
}

func (control *Control) UnloadInternal(internal string) error {
	return call(control.Object(), control.iface, "UnloadInternal", 0, args{internal}, nil)
}

func (control *Control) ResetXruns() error {
	return call(control.Object(), control.iface, "ResetXruns", 0, nil, nil)
}

func (control *Control) SwitchMaster() error {
	return call(control.Object(), control.iface, "SwitchMaster", 0, nil, nil)
}

// signal handlers

func (control *Control) OnServerStart(handler func()) (detach func() error, err error) {
	return control.eventRouter.addHandler(control.Object(), control.iface, "ServerStart", handler)
}

func (control *Control) OnServerStop(handler func()) (detach func() error, err error) {
	return control.eventRouter.addHandler(control.Object(), control.iface, "ServerStop", handler)
}
