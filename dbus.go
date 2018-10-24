package jackdbus

import "github.com/godbus/dbus"

type args []interface{}

func call(object dbus.BusObject, iface string, method string, flags dbus.Flags, args args, ret args) error {
	call := object.Call(iface+"."+method, flags, args...)
	if call.Err != nil {
		return call.Err
	}

	if ret != nil {
		return call.Store(ret...)
	}
	return nil
}
