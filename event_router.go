package jackdbus

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"sync"

	"github.com/godbus/dbus"
)

type eventKey struct {
	iface string
	name  string
	path  dbus.ObjectPath
}

type handler struct {
	callback reflect.Value
}

type eventRouter struct {
	sync.RWMutex

	handlers map[eventKey][]*handler
	conn     *dbus.Conn
}

func newEventRouter() *eventRouter {
	return &eventRouter{
		handlers: make(map[eventKey][]*handler),
	}
}

func (router *eventRouter) DeliverSignal(iface, name string, signal *dbus.Signal) {
	key := eventKey{iface, name, signal.Path}

	router.RLock()
	defer router.RUnlock()

	handlers := router.handlers[key]
	for _, handler := range handlers {
		go router.handleOne(handler, signal)
	}
}

func (router *eventRouter) handleOne(handler *handler, signal *dbus.Signal) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("catched panic handling '%s' on '%s': %v", signal.Name, signal.Path, err)
		}
	}()

	args := prepareArgs(signal.Body)
	handler.callback.Call(args)
}

func (router *eventRouter) addHandler(object dbus.BusObject, iface, name string, callback interface{}) (detach func() error, err error) {
	val := reflect.ValueOf(callback)
	if val.Type().Kind() != reflect.Func {
		return nil, errors.New("handler is not a function")
	}

	handler := &handler{callback: val}
	key := eventKey{iface, name, object.Path()}
	router.Lock()
	handlers, any := router.handlers[key]
	router.handlers[key] = append(handlers, handler)
	router.Unlock()

	detach = func() error {
		return router.removeHandler(object, iface, name, callback)
	}

	if !any {
		return detach, router.listen(object.Path(), iface, name)
	} else {
		return detach, nil
	}
}

func (router *eventRouter) removeHandler(object dbus.BusObject, iface, name string, callback interface{}) error {
	val := reflect.ValueOf(callback)
	if val.Type().Kind() != reflect.Func {
		return errors.New("handler is not a function")
	}

	noHandlerError := fmt.Errorf("no handler for '%s' on '%s'", iface, name)

	router.Lock()
	defer router.Unlock()

	key := eventKey{iface, name, object.Path()}
	handlers, ok := router.handlers[key]
	if !ok {
		return noHandlerError
	}
	for id, h := range handlers {
		if h.callback == val {
			copy(handlers[id:], handlers[id+1:])
			handlers[len(handlers)-1] = nil
			handlers = handlers[:len(handlers)-1]
			router.handlers[key] = handlers
		}
	}

	if len(handlers) == 0 {
		delete(router.handlers, key)
		return router.unlisten(object.Path(), iface, name)
	}

	return nil
}

func prepareArgs(raws []interface{}) []reflect.Value {
	var args []reflect.Value
	for _, raw := range raws {
		args = append(args, reflect.ValueOf(raw))
	}

	return args
}

func (router *eventRouter) listen(path dbus.ObjectPath, iface, name string) error {
	if router.conn == nil {
		return errors.New("router has no conn")
	}

	root := router.conn.Object("org.freedesktop.DBus", "/org/freedesktop/DBus")
	matchSelector := fmt.Sprintf("type='signal',path='%s',interface='%s',member='%s'", path, iface, name)
	return call(root, "org.freedesktop.DBus", "AddMatch", 0, args{matchSelector}, nil)
}

func (router *eventRouter) unlisten(path dbus.ObjectPath, iface, name string) error {
	if router.conn == nil {
		return errors.New("router has no conn")
	}

	root := router.conn.Object("org.freedesktop.DBus", "/org/freedesktop/DBus")
	matchSelector := fmt.Sprintf("type='signal',path='%s',interface='%s',member='%s'", path, iface, name)
	return call(root, "org.freedesktop.DBus", "RemoveMatch", 0, args{matchSelector}, nil)
}
