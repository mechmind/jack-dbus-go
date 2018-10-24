package jackdbus

type SessionManager struct {
	object
	iface string

	eventRouter *eventRouter
}

// methods

func (sm *SessionManager) GetClientNameByUUID(uuid string) (name string, err error) {
	err = call(sm.Object(), sm.iface, "GetClientNameByUuid", 0, args{uuid}, args{&name})
	return
}

func (sm *SessionManager) GetUUIDForClientName(name string) (uuid string, err error) {
	err = call(sm.Object(), sm.iface, "GetUuidForClientName", 0, args{name}, args{&uuid})
	return
}

func (sm *SessionManager) GetState() (t uint32, target string, err error) {
	err = call(sm.Object(), sm.iface, "GetState", 0, args{}, args{&t, &target})
	return
}

func (sm *SessionManager) HasSessionCallback(name string) (hasCallback bool, err error) {
	err = call(sm.Object(), sm.iface, "HasSessionCallback", 0, args{name}, args{&hasCallback})
	return
}

type NotifyResult struct {
	UUID       string
	ClientName string
	Command    string
	Flags      uint32
}

type NotifyType uint32

const (
	NotifySessionSave         NotifyType = 1
	NotifySessionSaveAndQuit             = 2
	NotifySessionSaveTemplate            = 3
)

func (sm *SessionManager) Notify(queue bool, target string, typ NotifyType, path string) (result []NotifyResult, err error) {
	err = call(sm.Object(), sm.iface, "Notify", 0, args{queue, target, typ, path}, args{&result})
	return
}

func (sm *SessionManager) ReserveClientName(name string, uuid string) (err error) {
	return call(sm.Object(), sm.iface, "ReserveClientName", 0, args{name, uuid}, nil)
}

// signals

func (sm *SessionManager) OnStateChanged(handler func(uint32, string)) (detach func() error, err error) {
	return sm.eventRouter.addHandler(sm.Object(), sm.iface, "StateChanged", handler)
}
