package jackdbus

import "github.com/godbus/dbus"

type object struct {
	obj dbus.BusObject
}

func (obj object) Object() dbus.BusObject {
	return obj.obj
}
